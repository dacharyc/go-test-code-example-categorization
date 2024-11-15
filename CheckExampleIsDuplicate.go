package main

func CheckExampleIsDuplicate(hashes map[string]bool, snippetHash string) bool {
	exists := hashes[snippetHash]
	if exists {
		return true
	} else {
		return false
	}
}
