package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"text/template"

	"github.com/sashabaranov/go-openai"
)

type ChatResponse struct {
	Response string `json:"response"`
}

func startWebServer(collectionsData CollectionData, chatHistory *[]openai.ChatCompletionMessage) {
	// Serve the main page
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		tmpl := template.Must(template.ParseFiles("html_files/index.html"))
		tmpl.Execute(w, nil)
	})

	// Handle chat requests
	http.HandleFunc("/chat", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}

		if !useTerminal {
			fmt.Printf("\n\nChat request received\n")
		}

		var request struct {
			Prompt string `json:"prompt"`
		}

		if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
			http.Error(w, "Invalid request body", http.StatusBadRequest)
			return
		}

		response := llmUserQuery(collectionsData, chatHistory, request.Prompt)

		json.NewEncoder(w).Encode(ChatResponse{
			Response: response,
		})
		if !useTerminal {
			fmt.Printf("Chat response generated and returned\n")
		}
	})

	if !useTerminal {
		fmt.Printf("Starting web server on http://localhost:8080\n")
	}

	if err := http.ListenAndServe(":8080", nil); err != nil {
		fmt.Printf("Error starting server: %v\n", err)
	}
}
