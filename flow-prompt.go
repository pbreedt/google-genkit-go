package main

import (
	"context"
	"log"
	"net/http"

	"github.com/firebase/genkit/go/ai"
	"github.com/firebase/genkit/go/core"
	"github.com/firebase/genkit/go/genkit"
	"github.com/firebase/genkit/go/plugins/googlegenai"
	"github.com/firebase/genkit/go/plugins/server"
)

type Menu struct {
	Theme string     `json:"theme"`
	Items []MenuItem `json:"items"`
}

type MenuItem struct {
	Name        string `json:"name"`
	Description string `json:"description"`
}

func flowPrompt() {
	ctx := context.Background()

	// Initialize Genkit with the Google AI plugin and Gemini 2.0 Flash.
	g, err := genkit.Init(ctx,
		genkit.WithPlugins(&googlegenai.GoogleAI{}),
		genkit.WithDefaultModel("googleai/gemini-2.0-flash"),
	)
	if err != nil {
		log.Fatalf("could not initialize Genkit: %w", err)
	}

	// Used in non-streaming call:
	// menuSuggestionFlow.Run(ctx, "bistro")
	menuSuggestionFlow := genkit.DefineStreamingFlow(g, "menuSuggestionFlow",
		func(ctx context.Context, theme string, callback core.StreamCallback[string]) (Menu, error) {
			item, _, err := genkit.GenerateData[MenuItem](ctx, g,
				ai.WithPrompt("Invent a menu item for a %s themed restaurant.", theme),
				ai.WithStreaming(func(ctx context.Context, chunk *ai.ModelResponseChunk) error {
					// Here, you could process the chunk in some way before sending it to
					// the output stream using StreamCallback. In this example, we output
					// the text of the chunk, unmodified.
					log.Printf("Process chunk: %s", chunk)
					if callback == nil { // non-streaming, like when invoked by HTTP POST
						return nil
					}
					return callback(ctx, "XXX"+chunk.Text())
					// return nil
				}),
			)
			if err != nil {
				return Menu{}, err
			}

			return Menu{
				Theme: theme,
				Items: []MenuItem{*item},
			}, nil
		})

	streamCh := menuSuggestionFlow.Stream(ctx, "bistro")

	// Handle streaming results
	for result, err := range streamCh {
		if err != nil {
			log.Fatal("Stream error: %v", err)
		}
		if result.Done {
			log.Printf("Menu with %s theme:\n", result.Output.Theme)
			for _, item := range result.Output.Items {
				log.Printf(" - %s: %s", item.Name, item.Description)
			}
		} else {
			log.Println("Stream chunk:", result.Stream)
		}
	}

	mux := http.NewServeMux()
	// for below: use genkit.DefineFlow instead of genkit.DefineStreamingFlow
	mux.HandleFunc("POST /menuSuggestionFlow", genkit.Handler(menuSuggestionFlow))
	// alternative: loop thru all flows:
	// for _, flow := range genkit.ListFlows(g) {
	// 	mux.HandleFunc("POST /"+flow.Name(), genkit.Handler(flow))
	// }

	// server.Start does normal ListenAndServe but add logic for graceful shutdown
	log.Fatal(server.Start(ctx, "127.0.0.1:3400", mux))
}
