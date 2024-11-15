package main

import (
	"context"
	"fmt"
	"github.com/tmc/langchaingo/llms/ollama"
	"log"
	"os"
)

func main() {
	files := GetFiles()

	// To change the model, use a different model's string name here
	llm, err := ollama.New(ollama.WithModel(MODEL))
	if err != nil {
		log.Fatalf("failed to connect to ollama: %v", err)
	}
	ctx := context.Background()
	hashes := make(map[string]bool)

	for _, file := range files {
		contents, err := os.ReadFile(file)
		if err != nil {
			fmt.Printf("failed to read file: %v\n", err)
			return
		}
		category := CategorizeSnippet(string(contents), llm, ctx)
		snippetHash := GetSnippetHash(string(contents))
		isDuplicate := CheckExampleIsDuplicate(hashes, snippetHash)
		if !isDuplicate {
			hashes[snippetHash] = true
		}

		// To generate an artifact of running this code, change this print statement to some logic - either
		// write to a file to create a CSV, or insert the info into a MongoDB collection (for our purposes)
		fmt.Printf("Filepath: %v, Category: %v, Duplicate: %v\n", file, category, isDuplicate)
	}
}
