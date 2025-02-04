package main

type RepoReport struct {
	TotalCodeBlocks        int                       `json:"total_code_blocks"`
	LLMCategorizedCount    int                       `json:"llm_categorized_count"`
	StringMatchedCount     int                       `json:"string_matched_count"`
	AccuracyEstimate       float64                   `json:"accuracy_estimate"`
	CategoryLanguageCounts map[string]map[string]int `json:"category_language_counts"`
}
