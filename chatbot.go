package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"github.com/sashabaranov/go-openai"
)

func startChatbot(collectionsData CollectionData) {

	chatHistory := []openai.ChatCompletionMessage{}

	if useTerminal {
		startTerminal(collectionsData, &chatHistory)
		return
	}

	go startWebServer(collectionsData, &chatHistory)
	select {}
}

func startTerminal(collectionsData CollectionData, chatHistory *[]openai.ChatCompletionMessage) {
	fmt.Println("Type your question (type 'bye' or 'q' to quit):")
	scanner := bufio.NewScanner(os.Stdin)
	for {
		fmt.Print(">>> ")
		if !scanner.Scan() {
			break
		}
		input := strings.TrimSpace(scanner.Text())

		if strings.ToLower(input) == "bye" || strings.ToLower(input) == "q" {
			break
		}
		response := llmUserQuery(collectionsData, chatHistory, input)
		fmt.Println(response)
		fmt.Println()
	}
	fmt.Println("Goodbye!")
}
