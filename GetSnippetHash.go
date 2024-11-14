package main

import (
	"crypto/sha256"
	"encoding/hex"
	"io"
	"log"
)

// GetSnippetHash removes whitespace from the string code example, and then creates a sha256 representation of it to
// save memory and use when calculating whether the code example is a duplicate
func GetSnippetHash(contents string) string {
	whitespaceStrippedString := SpaceMap(contents)
	hasher := sha256.New()
	_, err := io.WriteString(hasher, whitespaceStrippedString)
	if err != nil {
		log.Fatalf("failed to hash contents: %q\n, %q\n", contents, err)
	}
	return hex.EncodeToString(hasher.Sum(nil))
}
