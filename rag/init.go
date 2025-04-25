package rag

// Import Genkit's file-based vector retriever, (Don't use in production.)
import (
	"context"
	"log"

	"github.com/firebase/genkit/go/ai"

	"github.com/firebase/genkit/go/genkit"
	"github.com/firebase/genkit/go/plugins/googlegenai"

	"github.com/firebase/genkit/go/plugins/localvec"
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
