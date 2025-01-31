package main

import (
	"context"
	"fmt"
	"github.com/tmc/langchaingo/llms/ollama"
	"log"
	"os"
	"path/filepath"
	"strings"
	"time"
)

func main() {
	startTime := time.Now()
	files := GetFiles()
	projectName := "cloud-docs"

	// To change the model, use a different model's string name here
	llm, err := ollama.New(ollama.WithModel(MODEL))
	if err != nil {
		log.Fatalf("failed to connect to ollama: %v", err)
	}
	ctx := context.Background()
	//hashes := make(map[string]bool)

	var snippets []SnippetInfo
	counts := make(map[string]map[string]int)
	filesProcessed := 0

	LogStartInfoToConsole(startTime)

	for _, file := range files {
		contents, err := os.ReadFile(file)
		if err != nil {
			fmt.Printf("failed to read file: %v\n", err)
			return
		}
		category := CategorizeSnippet(string(contents), llm, ctx)
		//snippetHash := GetSnippetHash(string(contents))
		//isDuplicate := CheckExampleIsDuplicate(hashes, snippetHash)
		//if !isDuplicate {
		//	hashes[snippetHash] = true
		//}

		// Find the starting index of "cloud-docs"
		startIndex := strings.Index(file, projectName)
		pagePath := file[startIndex:]
		ext := filepath.Ext(file)
		lang := GetLangFromExtension(ext)
		details := SnippetInfo{
			Page:     pagePath,
			Category: category,
			Language: lang,
			//Duplicate: isDuplicate,
		}
		snippets = append(snippets, details)
		if _, exists := counts[details.Category]; !exists {
			counts[details.Category] = make(map[string]int)
		}
		// Increment the language count for the specific category
		counts[details.Category][details.Language]++
		filesProcessed++
		if filesProcessed%100 == 0 {
			fmt.Println("Processed ", filesProcessed, " snippets")
		}
	}
	LogFinishInfoToConsole(startTime, filesProcessed)

	WriteSnippetReport(snippets, projectName)

	WriteCategoryCountsReport(counts, projectName)
}
