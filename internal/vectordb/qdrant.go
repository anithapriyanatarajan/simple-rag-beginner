package vectordb

import (
	"context"
	"log"

	"github.com/google/uuid"
	"github.com/qdrant/go-client/qdrant"

	"simple-rag-beginner/internal/embedding"
)

var (
	client         *qdrant.Client
	collectionName = "rag_collection"
)

// InitQdrant initializes the Qdrant client and creates collection if needed
func InitQdrant(url string) error {
	// Create client with localhost configuration
	var err error
	client, err = qdrant.NewClient(&qdrant.Config{
		Host: "localhost",
		Port: 6334,
	})
	if err != nil {
		log.Printf("Failed to create Qdrant client: %v", err)
		return err
	}

	// Determine vector dimension by testing embedding
	testVec := embedding.TextToVector("test")
	vectorSize := len(testVec)
	log.Printf("Using vector dimension: %d", vectorSize)

	// Create collection with dynamic vector configuration
	err = client.CreateCollection(context.Background(), &qdrant.CreateCollection{
		CollectionName: collectionName,
		VectorsConfig: qdrant.NewVectorsConfig(&qdrant.VectorParams{
			Size:     uint64(vectorSize),
			Distance: qdrant.Distance_Cosine,
		}),
	})
	if err != nil {
		log.Printf("Collection might already exist: %v", err)
		// Don't return error if collection already exists
	}

	// Insert sample data with more diverse content
	samples := map[string]string{
		"greeting": "Hi there! This is context for 'hello'. Welcome to our system!",
		"weather":  "Today's weather is sunny and bright. Perfect day for outdoor activities.",
		"agent":    "Agents are autonomous entities that can perform tasks independently.",
		"identity": "You are a helpful AI assistant designed to answer questions and provide information.",
		"time":     "Time is a continuous progression of events from past to future.",
		"learning": "Machine learning involves training algorithms on data to make predictions.",
	}
	return InsertSampleData(samples)
}

// CloseQdrant closes the Qdrant client connection
func CloseQdrant() {
	if client != nil {
		client.Close()
	}
}

// InsertSampleData inserts a map of {id: text} into the Qdrant collection
func InsertSampleData(samples map[string]string) error {
	if client == nil {
		return nil // Skip if client not initialized
	}

	points := []*qdrant.PointStruct{}
	for _, text := range samples {
		vec := embedding.TextToVector(text)
		uid := uuid.New().String()
		points = append(points, &qdrant.PointStruct{
			Id: &qdrant.PointId{
				PointIdOptions: &qdrant.PointId_Uuid{Uuid: uid},
			},
			Vectors: qdrant.NewVectors(vec...),
			Payload: qdrant.NewValueMap(map[string]any{"text": text}),
		})
	}

	_, err := client.Upsert(context.Background(), &qdrant.UpsertPoints{
		CollectionName: collectionName,
		Points:         points,
	})
	return err
}

// RetrieveTopK searches for the top K most similar texts to the query
func RetrieveTopK(query string, k int) ([]string, error) {
	if client == nil {
		// Fallback to in-memory if client not available
		return []string{"No relevant context found."}, nil
	}

	vec := embedding.TextToVector(query)

	searchResult, err := client.Query(context.Background(), &qdrant.QueryPoints{
		CollectionName: collectionName,
		Query:          qdrant.NewQuery(vec...),
		Limit:          qdrant.PtrOf(uint64(k)),
		WithPayload:    qdrant.NewWithPayload(true),
	})
	if err != nil {
		return nil, err
	}

	results := []string{}
	for _, point := range searchResult {
		if payload := point.GetPayload(); payload != nil {
			if textValue := payload["text"]; textValue != nil {
				if text := textValue.GetStringValue(); text != "" {
					results = append(results, text)
				}
			}
		}
	}

	// Fallback if no results
	if len(results) == 0 {
		results = append(results, "No relevant context found.")
	}

	return results, nil
}
