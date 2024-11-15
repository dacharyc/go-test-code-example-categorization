package main

import "testing"

func TestStripWhitespaceWithSpaces(t *testing.T) {
	got := StripWhitespace("Hello World    ")
	expected := "HelloWorld"
	if got != expected {
		t.Errorf("got %q want %q", got, expected)
	}
}

func TestStripWhitespaceWithNewline(t *testing.T) {
	got := StripWhitespace("Hello World\n")
	expected := "HelloWorld"
	if got != expected {
		t.Errorf("got %q want %q", got, expected)
	}
}

func TestStripWhitespaceWithTab(t *testing.T) {
	got := StripWhitespace("Hello World	")
	expected := "HelloWorld"
	if got != expected {
		t.Errorf("got %q want %q", got, expected)
	}
}
