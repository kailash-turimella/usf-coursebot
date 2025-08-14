package main

import (
	"context"
	"encoding/csv"
	"fmt"
	"os"
	"strconv"
	"strings"

	chroma "github.com/amikos-tech/chroma-go"
	openai "github.com/amikos-tech/chroma-go/openai"
	"github.com/amikos-tech/chroma-go/types"
	"github.com/joho/godotenv"
	openaiapi "github.com/sashabaranov/go-openai"
)

type CollectionData struct {
	collection    *chroma.Collection
	context       context.Context
	embeddingFunc *openai.OpenAIEmbeddingFunction
	openaiClient  *openaiapi.Client
}

type Course struct {
	title    string // e.g. "CS110-01 Intro to Computer Science by Julia Nolfo"
	metadata map[string]any
}

func setUpCollections(fileName string) CollectionData {
	err := godotenv.Load()
	if err != nil {
		fmt.Printf("Error loading .env file: %v", err)
	}

	openaiKey := os.Getenv("OPENAI_API_KEY")
	if openaiKey == "" {
		fmt.Printf("OPENAI_API_KEY not set")
	}

	ctx := context.Background()

	client, err := chroma.NewClient(chromadbHost)
	if err != nil {
		fmt.Printf("Error creating client: %v", err)
		return CollectionData{}
	}

	embeddingFunction, err := openai.NewOpenAIEmbeddingFunction(openaiKey)
	if err != nil {
		fmt.Printf("Error creating embedding function: %v", err)
		return CollectionData{}
	}

	var loadFiles bool

	collections := CollectionData{
		collection:    createCollection(ctx, client, "Courses", embeddingFunction, &loadFiles),
		context:       ctx,
		embeddingFunc: embeddingFunction,
		openaiClient:  openaiapi.NewClient(openaiKey),
	}

	if loadFiles {
		if !useTerminal {
			fmt.Printf("Loading courses from %s...\n", fileName)
		}
		loadCourses(collections, fileName)
	}

	return collections
}

func createCollection(ctx context.Context, client *chroma.Client, collectionName string, embeddingFunction *openai.OpenAIEmbeddingFunction, loadFiles *bool) *chroma.Collection {
	collection, err := client.GetCollection(ctx, collectionName, embeddingFunction)
	if err != nil {
		metadata := map[string]interface{}{
			"description": fmt.Sprintf("%s collection for Fall 2024", collectionName),
		}

		collection, err = client.CreateCollection(ctx, collectionName, metadata, false, embeddingFunction, types.L2)
		if err != nil {
			fmt.Printf("Error creating collection %s: %v\n", collectionName, err)
			return nil
		}
		if !useTerminal {
			fmt.Printf("Created new collection: %s\n", collection.Name)
		}
		*loadFiles = true
	} else {
		fmt.Printf("Loaded existing collection: %s\n", collection.Name)
		*loadFiles = false
	}
	return collection
}

func loadCourses(collections CollectionData, fileName string) error {
	file, err := os.Open(fileName)
	if err != nil {
		fmt.Printf("failed to open CSV file: %v\n", err)
		return err
	}
	defer file.Close()

	reader := csv.NewReader(file)
	records, err := reader.ReadAll()
	if err != nil {
		fmt.Printf("failed to read CSV file: %v\n", err)
		return err
	}

	batchSize := 1000

	for i := 1; i < len(records); i += batchSize {
		end := min(i+batchSize, len(records))

		var ids []string
		var documents []string
		var metadatas []map[string]any
		var embeddings []*types.Embedding

		for j, row := range records[i:end] {
			meta := make(map[string]any)
			for x, colName := range records[0] {
				meta[colName] = row[x]
			}

			// document: e.g. "CS 110-01 Intro to Computer Science by Julia Nolfo"
			courseCode := row[getIndex(records[0], "SUBJ")] + " " + row[getIndex(records[0], "CRSE NUM")] + "-" + row[getIndex(records[0], "SEC")]
			title := row[getIndex(records[0], "Title Short Desc")]
			instructor := row[getIndex(records[0], "Primary Instructor First Name")] + " " + row[getIndex(records[0], "Primary Instructor Last Name")]
			document := courseCode + " " + title
			if strings.TrimSpace(instructor) == "" {
				document += " by " + instructor
			}

			ids = append(ids, strconv.Itoa(i+j+1))
			documents = append(documents, document)
			metadatas = append(metadatas, meta)
		}

		embeddings, err = collections.embeddingFunc.EmbedDocuments(collections.context, documents)
		if err != nil {
			fmt.Printf("failed to get embeddings: %v\n", err)
			return err
		}

		_, err = collections.collection.Add(collections.context, embeddings, metadatas, documents, ids)
		if err != nil {
			fmt.Printf("failed to add documents to the main collection: %v\n", err)
			return err
		}
		fmt.Printf("Added %d documents to the collection\n", len(documents))
	}
	return nil
}

func getIndex(header []string, colName string) int {
	for i, col := range header {
		if col == colName {
			return i
		}
	}
	return -1
}
