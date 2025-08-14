package main

import (
	"context"
	"fmt"

	"github.com/sashabaranov/go-openai"
)

func llmUserQuery(collectionData CollectionData, chatHistory *[]openai.ChatCompletionMessage, query string) string {
	if !useTerminal {
		fmt.Printf("Question: %s\n", query)
	}
	// -------- First LLM call --------
	*chatHistory = append(*chatHistory, []openai.ChatCompletionMessage{
		{
			Role:    openai.ChatMessageRoleSystem,
			Content: systemPrompt,
		},
		{
			Role:    openai.ChatMessageRoleUser,
			Content: query,
		},
	}...)

	firstMessage := getLLMMessage(collectionData, *chatHistory)

	// -------- No tool call --------
	if len(firstMessage.ToolCalls) == 0 {
		if !useTerminal {
			fmt.Printf("No Tool Calls\n")
		}
		*chatHistory = append(*chatHistory, openai.ChatCompletionMessage{
			Role:    openai.ChatMessageRoleAssistant,
			Content: "Relevant Courses I got by querying the vector database of all the courses:\n" + getRelavantCourses(collectionData, query),
		})
	} else {
		// -------- Handle tool calls --------
		if !useTerminal {
			fmt.Printf("%v Tool Calls\n", len(firstMessage.ToolCalls))
		}

		var toolCalls []openai.ToolCall
		var toolResponses []openai.ChatCompletionMessage
		for _, toolCall := range firstMessage.ToolCalls {
			if !useTerminal {
				fmt.Printf("Tool: %v\n", toolCall.Function.Name)
			}

			handler := handlerMap[toolCall.Function.Name]
			if handler == nil {
				fmt.Printf("Handler not found for tool: %v\n", toolCall.Function.Name)
				continue
			}

			toolOutput := handler(collectionData, toolCall.Function.Arguments)

			if toolOutput == "" {
				if !useTerminal {
					fmt.Printf("Incorrect Tool called, skipping to next tool\n")
				}
				continue
			}

			toolCalls = append(toolCalls, toolCall)
			toolResponses = append(toolResponses, openai.ChatCompletionMessage{
				Role:       openai.ChatMessageRoleTool,
				Name:       toolCall.Function.Name,
				ToolCallID: toolCall.ID,
				Content:    toolOutput,
			})
		}

		if len(toolCalls) > 0 {
			if !useTerminal {
				fmt.Printf("Correct Tool Calls: %v\n", len(toolCalls))
			}
			*chatHistory = append(*chatHistory, openai.ChatCompletionMessage{
				Role:      openai.ChatMessageRoleAssistant,
				ToolCalls: toolCalls,
			})
			*chatHistory = append(*chatHistory, toolResponses...)
		} else {
			if !useTerminal {
				fmt.Printf("No correct tool calls, falling back to relevant courses\n")
			}
			*chatHistory = append(*chatHistory, openai.ChatCompletionMessage{
				Role:    openai.ChatMessageRoleAssistant,
				Content: "Relevant Courses I got by querying the vector database of all the courses:\n" + getRelavantCourses(collectionData, query),
			})
		}
	}

	// -------- Second LLM call --------
	secondMessage := getLLMMessage(collectionData, *chatHistory)
	*chatHistory = append(*chatHistory, secondMessage)

	if !useTerminal && printResponse {
		fmt.Printf("Response: %s\n", secondMessage.Content)
	}

	return secondMessage.Content
}

func getLLMMessage(collectionData CollectionData, messages []openai.ChatCompletionMessage) openai.ChatCompletionMessage {
	resp, err := collectionData.openaiClient.CreateChatCompletion(
		context.Background(),
		openai.ChatCompletionRequest{
			Model:    openai.GPT4oMini,
			Messages: messages,
			Tools:    getAllTools(),
		},
	)
	if err != nil {
		fmt.Printf("ChatCompletion error: %v\n", err)
		return openai.ChatCompletionMessage{}
	}
	return resp.Choices[0].Message
}
