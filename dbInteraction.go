package main

import (
	"encoding/json"
	"fmt"
	"sort"

	chromago "github.com/amikos-tech/chroma-go"
	"github.com/sashabaranov/go-openai"
)

type LLMFilter struct {
	Column string `json:"column"`
	Value  string `json:"value"`
}

func getRelavantCourses(collectionsData CollectionData, question string) string {
	results, err := collectionsData.collection.Query(collectionsData.context, []string{question}, numRelevantCoursesToRetrieve, nil, nil, nil)
	if err != nil {
		fmt.Printf("Error retrieving relevant documents: %v\n", err)
		return "none\n"
	}
	retVal := "Relevant Courses (queried without where map):\n" + resultsToString(results)

	messages := []openai.ChatCompletionMessage{
		{
			Role:    openai.ChatMessageRoleSystem,
			Content: systemPrompt,
		},
		{
			Role:    openai.ChatMessageRoleUser,
			Content: getJSONPromptTemplate + question,
		},
	}

	whereMap := buildWhereFromLLMResponse(getLLMMessage(collectionsData, messages).Content)
	results, err = collectionsData.collection.Query(collectionsData.context, []string{question}, numRelevantCoursesToRetrieve, whereMap, nil, nil)
	if err != nil {
		fmt.Printf("Error retrieving relevant documents: %v\n", err)
		return retVal + "\n\nRelevant Courses (queried with where map):\nnone\n"
	}
	retVal += "\n\nRelevant Courses (queried with where map):\n" + resultsToString(results)

	return retVal
}

func buildWhereFromLLMResponse(JSONWhereMap string) map[string]any {
	var filter LLMFilter
	err := json.Unmarshal([]byte(JSONWhereMap), &filter)
	if err != nil {
		fmt.Println("Error unmarshalling JSON:", err)
		return nil
	}
	where := make(map[string]any)
	where[filter.Column] = filter.Value
	return where
}

func resultsToString(results *chromago.QueryResults) string {
	courses := ""

	for i := range results.Documents[0] {
		courses += "Document: " + results.Documents[0][i] + "\n"
		courses += "Metadata:\n" + getMetadata(results.Metadatas[0][i]) + "\n"
	}
	return courses
}

func getMetadata(metadata map[string]any) string {
	metadataStr := ""
	keys := []string{}
	for key := range metadata {
		keys = append(keys, key)
	}
	sort.Strings(keys)
	for _, key := range keys {
		metadataStr += fmt.Sprintf("  %-30s: %v\n", key, metadata[key])
	}
	return metadataStr
}
