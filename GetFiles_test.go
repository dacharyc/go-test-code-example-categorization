package main

import (
	"fmt"
	"os"
	"strings"
	"testing"
)

func TestStringArrayContainsPathsFromMultipleDirectories(t *testing.T) {
	exampleFilePaths := GetFiles()
	expectedManageIndexFilepath := "go-test-code-example-categorization/examples/manage-indexes/create-index-basic.go"
	containsFileFromManageIndexesDir := false
	expectedRunQueriesFilepath := "go-test-code-example-categorization/examples/run-queries/ann-basic.go"
	containsFileFromRunQueriesDir := false
	for _, exampleFilePath := range exampleFilePaths {
		if strings.Contains(exampleFilePath, expectedManageIndexFilepath) {
			containsFileFromManageIndexesDir = true
		}
		if strings.Contains(exampleFilePath, expectedRunQueriesFilepath) {
			containsFileFromRunQueriesDir = true
		}
	}

	if !containsFileFromManageIndexesDir || !containsFileFromRunQueriesDir {
		t.Error("Expected to find manage indexes and run queries files, but was missing one or both file paths")
	}
}

func TestStringArrayDoesNotContainDirectoryPath(t *testing.T) {
	exampleFilePaths := GetFiles()
	// Function to check if a path is a file
	isFile := func(path string) (bool, error) {
		info, err := os.Stat(path)
		if err != nil {
			return false, err
		}
		return !info.IsDir(), nil
	}
	// Check if all paths are files
	allPathsAreFiles := true
	for _, filePath := range exampleFilePaths {
		file, err := isFile(filePath)
		if err != nil {
			fmt.Printf("Error checking path %s: %v\n", filePath, err)
			allPathsAreFiles = false
			break
		}
		if !file {
			allPathsAreFiles = false
			break
		}
	}
	if !allPathsAreFiles {
		t.Error("found a directory in the array of paths")
	}
}
