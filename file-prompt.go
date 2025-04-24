package main

import (
	"context"
	"log"

	"github.com/firebase/genkit/go/ai"
	"github.com/firebase/genkit/go/genkit"
	"github.com/firebase/genkit/go/plugins/googlegenai"
)

/*
By default, uses prompt files under prompts/
*/

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
