package rag

// Import Genkit's file-based vector retriever, (Don't use in production.)
import (
	"context"
	"log"
	"net/http"

	"github.com/firebase/genkit/go/ai"

	"github.com/firebase/genkit/go/genkit"
	"github.com/firebase/genkit/go/plugins/googlegenai"
	"github.com/firebase/genkit/go/plugins/localvec"
	"github.com/firebase/genkit/go/plugins/server"
)

// Init returns a Genkit instancem, an Indexer, and a Retriever used in the Index and Retrieve functions
func Init() (*genkit.Genkit, *ai.Indexer, *ai.Retriever) {
	ctx := context.Background()

	// Vertex AI requires GOOGLE_CLOUD_PROJECT to be set and
	// export GOOGLE_CLOUD_LOCATION=us-central1
	// export GOOGLE_GENAI_USE_VERTEXAI=True
	// g, err := genkit.Init(ctx, genkit.WithPlugins(&googlegenai.VertexAI{}))

	// Used Google AI instead of Vertex AI
	g, err := genkit.Init(ctx,
		genkit.WithPlugins(&googlegenai.GoogleAI{}),
		genkit.WithDefaultModel("googleai/gemini-2.0-flash"),
	)
	if err != nil {
		log.Fatal(err)
	}

	if err = localvec.Init(); err != nil {
		log.Fatal(err)
	}

	// Also used Google AI here instead of Vertex AI
	indexer, retriever, err := localvec.DefineIndexerAndRetriever(g, "indexerRetreiever",
		localvec.Config{Embedder: googlegenai.GoogleAIEmbedder(g, "text-embedding-004")})
	// localvec.Config{Embedder: googlegenai.VertexAIEmbedder(g, "text-embedding-004")})
	if err != nil {
		log.Fatal(err)
	}

	return g, &indexer, &retriever
}

// Expose the flows via HTTP, allowing calls like:
// curl -X POST "http://localhost:3400/ragIndex" -H "Content-Type: application/json" -d '{"data": "./rag/menu.pdf"}'
// curl -X POST "http://localhost:3400/ragRetrieve" -H "Content-Type: application/json" -d '{"data": "Do you have meatballs?"}'

func Serve(g *genkit.Genkit) {
	ctx := context.Background()

	mux := http.NewServeMux()
	for _, flow := range genkit.ListFlows(g) {
		mux.HandleFunc("POST /"+flow.Name(), genkit.Handler(flow))
	}

	// server.Start does normal ListenAndServe but add logic for graceful shutdown
	log.Fatal(server.Start(ctx, "127.0.0.1:3400", mux))
}
