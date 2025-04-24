package main

/*
Explore Google GenKit for Go
See: https://firebase.google.com/docs/genkit-go/get-started-go

Requires:
	export GEMINI_API_KEY=<your-api-key>
https://console.cloud.google.com/apis/credentials?project=gen-lang-client-0457768792&pli=1
*/

import (
	"context"
	"log"

	"github.com/firebase/genkit/go/ai"
	"github.com/firebase/genkit/go/genkit"
	"github.com/firebase/genkit/go/plugins/googlegenai"
)

func main() {
	ctx := context.Background()

	// Initialize Genkit with the Google AI plugin and Gemini 2.0 Flash.
	g, err := genkit.Init(ctx,
		genkit.WithPlugins(&googlegenai.GoogleAI{}),
		genkit.WithDefaultModel("googleai/gemini-2.0-flash"),
	)
	if err != nil {
		log.Fatalf("could not initialize Genkit: %w", err)
	}
	// Alternative:
	// model := googlegenai.GoogleAIModelRef("gemini-2.0-flash", &googlegenai.GeminiConfig{
	// 	MaxOutputTokens: 500,
	// 	StopSequences:   ["<end>", "<fin>"],
	// 	Temperature:     0.5,
	// 	TopP:            0.4,
	// 	TopK:            50,
	// })
	// resp, err := genkit.Generate(ctx, g,
	// 	ai.WithModel(model),
	// 	ai.WithPrompt("What is the meaning of life."),
	// )

	resp, err := genkit.Generate(ctx, g,
		ai.WithPrompt("What is the meaning of life?"),
		ai.WithConfig(&googlegenai.GeminiConfig{
			MaxOutputTokens: 500,
			// Temperature:     0.5,
			// TopP:            0.4,
			// TopK:            50,
		}),
		ai.WithOutputType(FormattedResponse{}),
	)
	if err != nil {
		log.Fatal("could not generate model response: %w", err)
	}

	// Alternative:
	// formatted, resp, err := genkit.GenerateData[FormattedResponse](ctx, g,
	// 	ai.WithPrompt("What is the meaning of life."),
	// )
	log.Println(resp.Text())
}

type FormattedResponse struct {
	AiResponse string `json:"ai-response"`
	UserPrompt string `json:"user-prompt"`
}
