package main

import (
	"context"
	"fmt"
	"log"

	"github.com/firebase/genkit/go/ai"
	"github.com/firebase/genkit/go/genkit"
	"github.com/firebase/genkit/go/plugins/evaluators"
	"github.com/firebase/genkit/go/plugins/googlegenai"
)

/*
Start Genkit: "genkit start -- go run ."
Create dataset for "qaFlow", add exaples:

#	Input								Reference output
1	Who is man's best friend?			(?i)dog
2	Can I give milk to my cats?			(?i)(don't know|cannot answer)
3	From which animals did dogs evolve?	(?i)wolf|wolves
*/
func Evaluate() {
	ctx := context.Background()

	// Define metrics for evaluation
	metrics := []evaluators.MetricConfig{
		{
			MetricType: evaluators.EvaluatorRegex,
		},
	}

	// Initialize Genkit
	g, err := genkit.Init(ctx,
		genkit.WithPlugins(
			&googlegenai.GoogleAI{},
			&evaluators.GenkitEval{Metrics: metrics}, // Added metrics evaluator here
		),
		genkit.WithDefaultModel("googleai/gemini-2.0-flash"),
	)
	if err != nil {
		log.Fatalf("Genkit initialization error: %v", err)
	}

	// Dummy retriever that always returns the same facts
	dummyRetrieverFunc := func(ctx context.Context, req *ai.RetrieverRequest) (*ai.RetrieverResponse, error) {
		facts := []string{
			"Dog is man's best friend",
			"Dogs have evolved and were domesticated from wolves",
		}
		// Just return facts as documents.
		var docs []*ai.Document
		for _, fact := range facts {
			docs = append(docs, ai.DocumentFromText(fact, nil))
		}
		return &ai.RetrieverResponse{Documents: docs}, nil
	}
	factsRetriever := genkit.DefineRetriever(g, "local", "dogFacts", dummyRetrieverFunc)

	m := googlegenai.GoogleAIModel(g, "gemini-2.0-flash")
	if m == nil {
		log.Fatal("failed to find model")
	}

	// A simple question-answering flow
	genkit.DefineFlow(g, "qaFlow", func(ctx context.Context, query string) (string, error) {
		factDocs, err := ai.Retrieve(ctx, factsRetriever, ai.WithTextDocs(query))
		if err != nil {
			return "", fmt.Errorf("retrieval failed: %w", err)
		}
		llmResponse, err := genkit.Generate(ctx, g,
			ai.WithModelName("googleai/gemini-2.0-flash"),
			ai.WithPrompt("Answer this question with the given context: %s", query),
			ai.WithDocs(factDocs.Documents...),
		)
		if err != nil {
			return "", fmt.Errorf("generation failed: %w", err)
		}
		return llmResponse.Text(), nil
	})
}
