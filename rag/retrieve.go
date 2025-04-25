package rag

import (
	"context"
	"log"

	"github.com/firebase/genkit/go/ai"
	"github.com/firebase/genkit/go/genkit"
)

func Retrieve(g *genkit.Genkit, indexer *ai.Indexer, retriever *ai.Retriever) {
	ctx := context.Background()

	// Below all done in Init

	// g, err := genkit.Init(ctx,
	// 	genkit.WithPlugins(&googlegenai.GoogleAI{}),
	// 	genkit.WithDefaultModel("googleai/gemini-1.5-flash"),
	// )
	// if err != nil {
	// 	log.Fatal(err)
	// }

	// if err := localvec.Init(); err != nil {
	// 	log.Fatal(err)
	// }

	// _, menuPdfRetriever, err := localvec.DefineIndexerAndRetriever(
	// 	g, "menuQA",
	// 	localvec.Config{Embedder: googlegenai.GoogleAIEmbedder(g, "text-embedding-004")},
	// )

	// if err != nil {
	// 	log.Fatal(err)
	// }

	var err error

	retrieveFlow := genkit.DefineFlow(
		g, "ragRetrieve",
		func(ctx context.Context, question string) (string, error) {
			// Retrieve text relevant to the user's question.
			resp, err := ai.Retrieve(ctx, *retriever, ai.WithTextDocs(question))

			if err != nil {
				return "", err
			}

			// Call Generate, including the menu information in your prompt.
			return genkit.GenerateText(ctx, g,
				ai.WithModelName("googleai/gemini-2.0-flash"),
				ai.WithDocs(resp.Documents...),
				// Prompt for the menu sample:
				ai.WithSystem(`
You are acting as a helpful AI assistant that can answer questions about the
food available on the menu at Frattelino's Italian Restaurant.
Use only the context provided to answer the question. If you don't know, do not
make up an answer. Do not add or change items on the menu.`),
				ai.WithPrompt(question),
			)
		})

	res, err := retrieveFlow.Run(ctx, "What are the specials on Monday?")
	if err != nil {
		log.Printf("Error running flow: %v", err)
		log.Fatal(err)
	}

	log.Println(res)
}
