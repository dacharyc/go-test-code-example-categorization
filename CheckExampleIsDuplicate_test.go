package main

import (
	"os"
	"testing"
)

func TestCheckExampleIsDuplicateYes(t *testing.T) {
	firstExampleFilepath := RelSnippetsStartDirectory + "/other/insertOne.sh"
	firstExample, err := os.ReadFile(firstExampleFilepath)
	if err != nil {
		t.Errorf("failed to read file %v", err)
	}
	firstExampleHash := GetSnippetHash(string(firstExample))
	secondExampleFilepath := RelSnippetsStartDirectory + "/other/insertOneDuplicate.sh"
	secondExample, err := os.ReadFile(secondExampleFilepath)
	if err != nil {
		t.Errorf("failed to read file %v", err)
	}
	secondExampleHash := GetSnippetHash(string(secondExample))
	hashes := make(map[string]bool)
	hashes[firstExampleHash] = true
	isDuplicate := CheckExampleIsDuplicate(hashes, secondExampleHash)
	if !isDuplicate {
		t.Errorf("expected the example to be a duplicate but it wasn't.")
	}
}

func TestCheckExampleIsDuplicateNo(t *testing.T) {
	firstExampleFilepath := RelSnippetsStartDirectory + "/other/insertOne.sh"
	firstExample, err := os.ReadFile(firstExampleFilepath)
	if err != nil {
		t.Errorf("failed to read file %v", err)
	}
	firstExampleHash := GetSnippetHash(string(firstExample))
	secondExampleFilepath := RelSnippetsStartDirectory + "/other/returnExample.sh"
	secondExample, err := os.ReadFile(secondExampleFilepath)
	if err != nil {
		t.Errorf("failed to read file %v", err)
	}
	secondExampleHash := GetSnippetHash(string(secondExample))
	hashes := make(map[string]bool)
	hashes[firstExampleHash] = true
	isDuplicate := CheckExampleIsDuplicate(hashes, secondExampleHash)
	if isDuplicate {
		t.Errorf("expected the example to not be a duplicate but it was.")
	}
}
