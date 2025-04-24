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

func main() {
	// flowPrompt()
	filePrompt()
	// simplePrompt()
	// select {} // keep the program running, required by genkit dev ui
}
