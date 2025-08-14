package main

import (
	"strings"
	"testing"

	"github.com/sashabaranov/go-openai"
)

func TestPhil(t *testing.T) {
	collection := setUpCollections("Fall 2024 Class Schedule 08082024.csv")
	query := "What courses is Phil Peterson teaching in Fall 2024?"
	chatHistory := []openai.ChatCompletionMessage{}
	result := llmUserQuery(collection, &chatHistory, query)
	if !strings.Contains(result, "40646") {
		t.Errorf("result does not include CS272: %s'", result)
	}
}

func TestPHIL(t *testing.T) {
	collection := setUpCollections("Fall 2024 Class Schedule 08082024.csv")
	query := "Which philosophy courses are offered this semester?"
	chatHistory := []openai.ChatCompletionMessage{}
	result := llmUserQuery(collection, &chatHistory, query)
	if !strings.Contains(result, "PHIL") {
		t.Errorf("result does not include PHIL: %s'", result)
	}
}

func TestBio(t *testing.T) {
	collection := setUpCollections("Fall 2024 Class Schedule 08082024.csv")
	query := "Where does Bioinformatics meet?"
	chatHistory := []openai.ChatCompletionMessage{}
	result := llmUserQuery(collection, &chatHistory, query)
	if !strings.Contains(result, "Bioinformatics") {
		t.Errorf("empty result for query '%s'", result)
	}
}

func TestGuitar(t *testing.T) {
	collection := setUpCollections("Fall 2024 Class Schedule 08082024.csv")
	query := "Can I learn guitar this semester?"
	chatHistory := []openai.ChatCompletionMessage{}
	result := llmUserQuery(collection, &chatHistory, query)
	if !strings.Contains(result, "MUS121") {
		t.Errorf("empty result for query '%s'", result)
	}
}

func TestMultiple(t *testing.T) {
	collection := setUpCollections("Fall 2024 Class Schedule 08082024.csv")
	query := "I would like to take a Rhetoric course from Phil Choong. What can I take?"
	chatHistory := []openai.ChatCompletionMessage{}
	result := llmUserQuery(collection, &chatHistory, query)
	if !strings.Contains(result, "RHET") {
		t.Errorf("empty result for query '%s'", result)
	}
}
