package main

// import "github.com/pbreedt/google-genkit-go/rag"
import "github.com/pbreedt/google-genkit-go/rag"

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

func main() {
	rag.DogFacts()
	// Evaluate()
	// g, i, r := rag.Init()
	// rag.Index(g, i, r)
	// rag.Retrieve(g, i, r)
	// rag.Serve(g)
	// tools()
	// flowPrompt()
	// filePrompt()
	// simplePrompt()
	select {} // keep the program running, required by genkit dev ui
}
