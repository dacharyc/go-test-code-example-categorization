package main

type RepoReport struct {
	LLMCategorizedCount    int                       `json:"llm_categorized_count"`
	StringMatchedCount     int                       `json:"string_matched_count"`
	CategoryLanguageCounts map[string]map[string]int `json:"category_language_counts"`
}
