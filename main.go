package main

/*
Explore Google GenKit for Go
See: https://firebase.google.com/docs/genkit-go/get-started-go

Requires:
	export GEMINI_API_KEY=<your-api-key>
https://console.cloud.google.com/apis/credentials?project=gen-lang-client-0457768792&pli=1

Dev tools:
	npm i -g genkit-cli
		<running genkit dev UI to quickly try out different models, temps & prompts>
		genkit start -- go run .
		below require already running runtime in a separate terminal with the GENKIT_ENV=dev environment variable set.
		genkit flow:run <flowName>
		genkit eval:flow <flowName>
*/

import (
	"context"
	"log"
	"net/http"

	"github.com/firebase/genkit/go/plugins/server"

	"github.com/firebase/genkit/go/ai"
	"github.com/firebase/genkit/go/core"
	"github.com/firebase/genkit/go/genkit"
	"github.com/firebase/genkit/go/plugins/googlegenai"
)

func main() {
	flowPrompt()
	// filePrompt()
	// simplePrompt()
	select {} // keep the program running, required by genkit dev ui
}

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

func filePrompt() {
	ctx := context.Background()

	// Initialize Genkit with the Google AI plugin and Gemini 2.0 Flash.
	g, err := genkit.Init(ctx,
		genkit.WithPlugins(&googlegenai.GoogleAI{}),
		genkit.WithDefaultModel("googleai/gemini-2.0-flash"),
	)
	if err != nil {
		log.Fatalf("could not initialize Genkit: %w", err)
	}

	menuPrompt := genkit.LookupPrompt(g, "menu")
	if menuPrompt == nil {
		log.Fatal("no prompt named 'menu' found")
	}

	resp, err := menuPrompt.Execute(context.Background(),
		ai.WithInput(map[string]any{"theme": "medieval"}),
	)
	if err != nil {
		log.Fatal(err)
	}

	var output map[string]any
	if err := resp.Output(&output); err != nil {
		log.Fatal(err)
	}

	log.Println(output["dishname"])
	log.Println(output["description"])
	/*
		Output:
			2025/04/24 11:41:51 The King's Beef Stew
			2025/04/24 11:41:51 A hearty stew of slow-cooked beef, root vegetables (carrots, parsnips, turnips), and barley, simmered in a rich ale broth. Served with a crusty bread trencher.
	*/
}

func simplePrompt() {
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
