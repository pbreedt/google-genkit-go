package rag

// Import Genkit's file-based vector retriever, (Don't use in production.)
import (
	"context"
	"io"
	"log"

	"github.com/firebase/genkit/go/ai"
	"github.com/firebase/genkit/go/genkit"

	"github.com/ledongthuc/pdf"

	"github.com/tmc/langchaingo/textsplitter"
)

// Requires Init to be called first

func Index(g *genkit.Genkit, indexer *ai.Indexer, retriever *ai.Retriever) {
	// ctx := context.Background()

	// Below all done in Init

	// g, err := genkit.Init(ctx,
	// 	genkit.WithPlugins(&googlegenai.GoogleAI{}),
	// 	genkit.WithDefaultModel("googleai/gemini-2.0-flash"),
	// )
	// if err != nil {
	// 	log.Fatal(err)
	// }

	// if err = localvec.Init(); err != nil {
	// 	log.Fatal(err)
	// }

	// menuPDFIndexer, _, err := localvec.DefineIndexerAndRetriever(g, "menuQA",
	// localvec.Config{Embedder: googlegenai.GoogleAIEmbedder(g, "text-embedding-004")})
	// if err != nil {
	// 	log.Fatal(err)
	// }

	// var err error

	splitter := textsplitter.NewRecursiveCharacter(
		textsplitter.WithChunkSize(200),
		textsplitter.WithChunkOverlap(20),
	)

	// indexPDFFlow := genkit.DefineFlow(
	genkit.DefineFlow(
		g, "ragIndex",
		func(ctx context.Context, path string) (any, error) {
			log.Printf("Indexing %s", path)
			// Extract plain text from the PDF. Wrap the logic in Run so it
			// appears as a step in your traces.
			pdfText, err := genkit.Run(ctx, "extract", func() (string, error) {
				return readPDF(path)
			})
			if err != nil {
				return nil, err
			}

			// Split the text into chunks. Wrap the logic in Run so it appears as a
			// step in your traces.
			docs, err := genkit.Run(ctx, "chunk", func() ([]*ai.Document, error) {
				// log.Printf("Splitting text '%s'", pdfText)
				chunks, err := splitter.SplitText(pdfText)
				if err != nil {
					log.Printf("Error splitting text: %v", err)
					return nil, err
				}

				var docs []*ai.Document
				for _, chunk := range chunks {
					// log.Printf("Adding Chunk: %s", chunk)
					docs = append(docs, ai.DocumentFromText(chunk, nil))
				}
				return docs, nil
			})
			if err != nil {
				log.Printf("Error in chunking: %v", err)
				return nil, err
			}
			log.Println("Done chunking")

			// Add chunks to the index.
			err = ai.Index(ctx, *indexer, ai.WithDocs(docs...))
			if err != nil {
				log.Printf("Error saving index: %v", err)
				return nil, err
			}

			log.Println("Done indexing menu")
			return nil, err
		},
	)

	// _, err = indexPDFFlow.Run(ctx, "./rag/menu.pdf")
	// if err != nil {
	// 	log.Printf("Error running flow: %v", err)
	// 	log.Fatal(err)
	// }

}

// Helper function to extract plain text from a PDF. Excerpted from
// https://github.com/ledongthuc/pdf
func readPDF(path string) (string, error) {
	log.Printf("Reading PDF %s", path)
	f, r, err := pdf.Open(path)
	if f != nil {
		defer f.Close()
	}
	if err != nil {
		return "", err
	}

	reader, err := r.GetPlainText()
	if err != nil {
		return "", err
	}

	bytes, err := io.ReadAll(reader)
	if err != nil {
		return "", err
	}

	return string(bytes), nil
}
