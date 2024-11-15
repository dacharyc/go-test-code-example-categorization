package main

import (
	"context"
	"github.com/tmc/langchaingo/llms/ollama"
	"log"
	"os"
	"path/filepath"
	"testing"
)

func TestCategorizeSnippetAPIMethod(t *testing.T) {
	llm, err := ollama.New(ollama.WithModel(MODEL))
	if err != nil {
		log.Fatalf("failed to connect to ollama: %v", err)
	}
	ctx := context.Background()
	testFilePath, err := filepath.Abs("../go-test-code-example-categorization/examples/other/insertOne.sh")
	if err != nil {
		log.Fatalf("failed to construct the file path: %v", err)
	}
	contents, err := os.ReadFile(testFilePath)
	if err != nil {
		log.Fatalf("failed to read the file at %v: %v", testFilePath, err)
	}
	got := CategorizeSnippet(string(contents), llm, ctx)
	if got != "API method signature" {
		t.Errorf("got %v, want %v", got, "API method signature")
	}
}

func TestCategorizeSnippetAPIMethodWithValues(t *testing.T) {
	llm, err := ollama.New(ollama.WithModel(MODEL))
	if err != nil {
		log.Fatalf("failed to connect to ollama: %v", err)
	}
	ctx := context.Background()
	testFilePath, err := filepath.Abs("../go-test-code-example-categorization/examples/other/api-method.go")
	if err != nil {
		log.Fatalf("failed to construct the file path: %v", err)
	}
	contents, err := os.ReadFile(testFilePath)
	if err != nil {
		log.Fatalf("failed to read the file at %v: %v", testFilePath, err)
	}
	got := CategorizeSnippet(string(contents), llm, ctx)
	expectation := "API method signature"
	if got != expectation {
		t.Errorf("got %v, want %v", got, expectation)
	}
}

func TestCategorizeConfigExample(t *testing.T) {
	llm, err := ollama.New(ollama.WithModel(MODEL))
	if err != nil {
		log.Fatalf("failed to connect to ollama: %v", err)
	}
	ctx := context.Background()
	testFilePath, err := filepath.Abs("../go-test-code-example-categorization/examples/other/configExample.yaml")
	if err != nil {
		log.Fatalf("failed to construct the file path: %v", err)
	}
	contents, err := os.ReadFile(testFilePath)
	if err != nil {
		log.Fatalf("failed to read the file at %v: %v", testFilePath, err)
	}
	got := CategorizeSnippet(string(contents), llm, ctx)
	expectation := "Example configuration object"
	if got != expectation {
		t.Errorf("got %v, want %v", got, expectation)
	}
}

func TestCategorizeSimpleReturnExample(t *testing.T) {
	llm, err := ollama.New(ollama.WithModel(MODEL))
	if err != nil {
		log.Fatalf("failed to connect to ollama: %v", err)
	}
	ctx := context.Background()
	testFilePath, err := filepath.Abs("../go-test-code-example-categorization/examples/other/returnExample.sh")
	if err != nil {
		log.Fatalf("failed to construct the file path: %v", err)
	}
	contents, err := os.ReadFile(testFilePath)
	if err != nil {
		log.Fatalf("failed to read the file at %v: %v", testFilePath, err)
	}
	got := CategorizeSnippet(string(contents), llm, ctx)
	expectation := "Example return object"
	if got != expectation {
		t.Errorf("got %v, want %v", got, expectation)
	}
}

// This test is currently failing - the LLMs seem to assess multi return examples as Task-based usage
// Should further tweak prompt until this passes
func TestCategorizeMultiReturnExample(t *testing.T) {
	llm, err := ollama.New(ollama.WithModel(MODEL))
	if err != nil {
		log.Fatalf("failed to connect to ollama: %v", err)
	}
	ctx := context.Background()
	testFilePath, err := filepath.Abs("../go-test-code-example-categorization/examples/other/runQueriesReturnExample.sh")
	if err != nil {
		log.Fatalf("failed to construct the file path: %v", err)
	}
	contents, err := os.ReadFile(testFilePath)
	if err != nil {
		log.Fatalf("failed to read the file at %v: %v", testFilePath, err)
	}
	got := CategorizeSnippet(string(contents), llm, ctx)
	expectation := "Example return object"
	if got != expectation {
		t.Errorf("got %v, want %v", got, expectation)
	}
}

func TestCategorizeTaskBasedUsage(t *testing.T) {
	llm, err := ollama.New(ollama.WithModel(MODEL))
	if err != nil {
		log.Fatalf("failed to connect to ollama: %v", err)
	}
	ctx := context.Background()
	testFilePath, err := filepath.Abs("../go-test-code-example-categorization/examples/manage-indexes/drop-index.go")
	if err != nil {
		log.Fatalf("failed to construct the file path: %v", err)
	}
	contents, err := os.ReadFile(testFilePath)
	if err != nil {
		log.Fatalf("failed to read the file at %v: %v", testFilePath, err)
	}
	got := CategorizeSnippet(string(contents), llm, ctx)
	expectation := "Task-based usage"
	if got != expectation {
		t.Errorf("got %v, want %v", got, expectation)
	}
}
