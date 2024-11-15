package main

import (
	"testing"
)

func TestGetSnippetHash(t *testing.T) {
	got := GetSnippetHash("Hello World")
	expected := "872e4e50ce9990d8b041330c47c9ddd11bec6b503ae9386a99da8584e9bb12c4"
	if got != expected {
		t.Errorf("got %q, want %q", got, expected)
	}
}
