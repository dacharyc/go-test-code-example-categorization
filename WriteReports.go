package main

import (
	"encoding/json"
	"fmt"
	"os"
)

func WriteSnippetReport(snippets []SnippetInfo, projectName string) {
	fmt.Println("Writing snippet report")
	snippetJsonData, marshallingErr := json.MarshalIndent(snippets, "", "  ")
	if marshallingErr != nil {
		fmt.Println("Error marshalling JSON:", marshallingErr)
		return
	}
	reportOutputDir := BaseReportOutputDir + projectName
	mkdirErr := os.MkdirAll(reportOutputDir, 0755)
	if mkdirErr != nil {
		fmt.Println("Error creating directory: ", mkdirErr)
		return
	}
	snippetDetailsFilepath := BaseReportOutputDir + projectName + "/snippets.json"
	writeReportErr := os.WriteFile(snippetDetailsFilepath, snippetJsonData, 0644)
	if writeReportErr != nil {
		fmt.Println("Error writing JSON to file:", writeReportErr)
		return
	}
	fmt.Println("Snippet report successfully written to", snippetDetailsFilepath)
}

func WriteCategoryCountsReport(counts map[string]map[string]int, projectName string) {
	categorySums := GetCategorySums(counts)
	repoData, jsonMarshallingErr := json.MarshalIndent(categorySums, "", "  ")
	if jsonMarshallingErr != nil {
		fmt.Println("Error marshalling JSON:", jsonMarshallingErr)
		return
	}
	fmt.Println("Writing category and language counts report")
	filePath := BaseReportOutputDir + projectName + "/language_category_counts.json"
	writeReportErr := os.WriteFile(filePath, repoData, 0644)
	if writeReportErr != nil {
		fmt.Println("Error writing JSON to file: ", writeReportErr)
		return
	}
	fmt.Println("Category and language counts report successfully written to", filePath)
}
