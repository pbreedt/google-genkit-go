package main

import (
	"context"
	"fmt"
	"log"

	"github.com/firebase/genkit/go/ai"
	"github.com/firebase/genkit/go/genkit"
	"github.com/firebase/genkit/go/plugins/googlegenai"
)

type ToolInput struct {
	Location string `json:"location"`
}

func tools() {
	ctx := context.Background()

	g, err := genkit.Init(ctx,
		genkit.WithPlugins(&googlegenai.GoogleAI{}),
		genkit.WithDefaultModel("googleai/gemini-2.0-flash"),
	)
	if err != nil {
		log.Fatal(err)
	}

	// var not required when using option 3 - prompt file
	// genkit.DefineTool(
	getWeatherTool := genkit.DefineTool(
		g, "getWeather", "Gets the current weather in a given location",
		func(ctx *ai.ToolContext, input ToolInput) (string, error) {
			// Here, we would typically make an API call or database query.
			// For this example, we just return a fixed value.
			return fmt.Sprintf("The current weather in %s is 63°F and sunny.", input.Location), nil
		})

	// Invoke option 1: genkit.Generate
	// resp, err := genkit.Generate(ctx, g,
	// 	ai.WithPrompt("What is the weather in San Francisco?"),
	// 	ai.WithTools(getWeatherTool),
	// )
	// if err != nil {
	// 	log.Fatal(err)
	// }

	// Invoke option 2: ai.Generate
	weatherPrompt, err := genkit.DefinePrompt(g, "weatherPrompt",
		ai.WithPrompt("What is the weather in {{location}}?"),
		ai.WithTools(getWeatherTool),
	)
	if err != nil {
		log.Fatal(err)
	}

	resp, err := weatherPrompt.Execute(ctx,
		ai.WithInput(map[string]any{"location": "San Francisco"}),
	)

	// Invoke option 3: genkit.LookupPrompt with prompt file
	// in prompt file, tool name [getWeather] corresponds to name provided in genkit.DefineTool
	// weatherPrompt := genkit.LookupPrompt(g, "weather")
	// if weatherPrompt == nil {
	// 	log.Fatal("no prompt named 'weatherPrompt' found")
	// }

	// resp, err := weatherPrompt.Execute(ctx,
	// 	ai.WithInput(map[string]any{"location": "Chicago"}),
	// )
	// if err != nil {
	// 	log.Fatal(err)
	// }

	log.Println(resp.Text())

}
