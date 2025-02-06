package main

import (
	"context"
	"fmt"
	"github.com/tmc/langchaingo/llms"
	"github.com/tmc/langchaingo/llms/ollama"
	"github.com/tmc/langchaingo/prompts"
	"log"
	"strings"
)

func HasStringMatchPrefix(contents string, langCategory string) (string, bool) {
	// These prefixes relate to usage examples
	importPrefix := "import "
	fromPrefix := "from "
	namespacePrefix := "namespace "
	packagePrefix := "package "
	usingPrefix := "using "
	mongoConnectionStringPrefix := "mongodb://"

	// These prefixes relate to command-line commands that *aren't* MongoDB specific, such as other tools, package managers, etc.
	mkdir := "mkdir "
	cd := "cd "
	docker := "docker "
	dockerCompose := "docker-compose "
	brew := "brew "
	yum := "yum "
	apt := "apt-"
	npm := "npm "
	pip := "pip "
	goRun := "go run "
	node := "node "
	dotnet := "dotnet "
	export := "export "
	jq := "jq "
	vi := "vi "

	usageExamplePrefixes := []string{importPrefix, fromPrefix, namespacePrefix, packagePrefix, usingPrefix, mongoConnectionStringPrefix}
	nonMongoPrefixes := []string{mkdir, cd, docker, dockerCompose, dockerCompose, brew, yum, apt, npm, pip, goRun, node, dotnet, export, jq, vi}

	if langCategory == SHELL {
		if strings.HasPrefix(contents, "atlas ") {
			return AtlasCliCommand, true
		} else if strings.HasPrefix(contents, "mongosh ") {
			return MongoshCommand, true
		} else {
			for _, prefix := range nonMongoPrefixes {
				if strings.HasPrefix(contents, prefix) {
					return NonMongoCommand, true
				}
			}
			return "Uncategorized", false
		}
	} else if langCategory == TEXT {
		if strings.HasPrefix(contents, "atlas ") {
			return AtlasCliCommand, true
		} else if strings.HasPrefix(contents, "mongosh ") {
			return MongoshCommand, true
		} else {
			for _, prefix := range nonMongoPrefixes {
				if strings.HasPrefix(contents, prefix) {
					return NonMongoCommand, true
				}
			}
			for _, prefix := range usageExamplePrefixes {
				if strings.HasPrefix(contents, prefix) {
					return UsageExample, true
				}
			}
			return "Uncategorized", false
		}
	} else {
		for _, prefix := range usageExamplePrefixes {
			if strings.HasPrefix(contents, prefix) {
				return UsageExample, true
			}
		}
		return "Uncategorized", false
	}
}

func ContainsString(contents string) (string, bool) {
	// These strings are typically included in usage examples
	aggregationExample := ".aggregate"
	mongoConnectionStringPrefix := "mongodb://"

	// These strings are typically included in return objects
	errorString := "error"
	warningString := "warning"
	deprecatedString := "deprecated"
	idString := "_id"

	// Some of the examples can be quite long. For the current case, we only care if `.aggregate` appears near the beginning of the example
	substringLengthToCheck := 50
	usageExampleSubstringsToEvaluate := []string{aggregationExample, mongoConnectionStringPrefix}
	returnObjectStringsToEvaluate := []string{errorString, warningString, deprecatedString, idString}

	if substringLengthToCheck < len(contents) {
		substring := contents[:substringLengthToCheck]
		for _, exampleString := range usageExampleSubstringsToEvaluate {
			if strings.Contains(substring, exampleString) {
				return UsageExample, true
			}
		}
		for _, exampleString := range returnObjectStringsToEvaluate {
			if strings.Contains(substring, exampleString) {
				return ExampleReturnObject, true
			}
		}
	} else {
		for _, exampleString := range usageExampleSubstringsToEvaluate {
			if strings.Contains(contents, exampleString) {
				return UsageExample, true
			}
		}
		for _, exampleString := range returnObjectStringsToEvaluate {
			if strings.Contains(contents, exampleString) {
				return ExampleReturnObject, true
			}
		}
	}
	return "Uncategorized", false
}

// CheckForStringMatch The bool we return from this func represents whether the string matching was successful
func CheckForStringMatch(contents string, lang string) (string, bool) {
	langCategory := GetLanguageCategory(lang)
	category, hasPrefix := HasStringMatchPrefix(contents, lang)
	if hasPrefix {
		return category, hasPrefix
	} else if langCategory != JSON_LIKE {
		// Currently, the only string matching we need to do is for an '.aggregate' method call, which should not appear in a
		// JSON-like example
		thisCategory, containsExampleString := ContainsString(contents)
		if containsExampleString {
			return thisCategory, containsExampleString
		} else {
			return "Uncategorized", false
		}
	} else {
		return "Uncategorized", false
	}
}

func ProcessSnippet(contents string, lang string, llm *ollama.LLM, ctx context.Context) (string, bool) {
	var category string
	validCategories := []string{AtlasCliCommand, ApiMethodSignature, ExampleReturnObject, ExampleConfigurationObject, MongoshCommand, NonMongoCommand, UsageExample}

	/* If the start characters of the code example match a pattern we have defined for a given category,
	 * return the category - no need to get the LLM involved.
	 */
	category, stringMatchSuccessful := CheckForStringMatch(contents, lang)
	if stringMatchSuccessful {
		/* The bool we are returning from this func represents whether the LLM categorized the snippet
		 * If we have successfully used string matching to categorize the snippet, the LLM does not process it, so we
		 * return false here
		 */
		return category, false
	} else {
		category = LLMAssignCategory(contents, lang, llm, ctx)
		/* I initially implemented this loop to ask the LLM to try again to categorize code examples that it couldn't categorize
		 * I found that even after retrying, the LLM cannot categorize "uncategorized" examples based on our current definitions
		 * Removing this loop for now
		 */
		//for attemptCounter < 3 {
		//	category = LLMAssignCategory(contents, lang, llm, ctx)
		//	if containsString(validCategories, category) {
		//		return category, attemptCounter
		//	} else {
		//		attemptCounter++
		//	}
		//}
		//return "Uncategorized", attemptCounter
		if containsString(validCategories, category) {
			return category, true
		} else {
			return "Uncategorized", true
		}
	}
}

func GetLanguageCategory(lang string) string {
	jsonLike := []string{JSON, XML, YAML}
	driversLanguagesMinusJS := []string{C, CPP, CSHARP, GO, JAVA, KOTLIN, PHP, PYTHON, RUBY, RUST, SCALA, SWIFT, TYPESCRIPT}
	if containsString([]string{SHELL}, lang) {
		return SHELL
	} else if containsString(jsonLike, lang) {
		return JSON_LIKE
	} else if containsString(driversLanguagesMinusJS, lang) {
		return DRIVERS_MINUS_JS
	} else if lang == JAVASCRIPT {
		return JAVASCRIPT
	} else if lang == TEXT {
		return TEXT
	} else {
		return "Unknown language"
	}
}

func LLMAssignCategory(contents string, lang string, llm *ollama.LLM, ctx context.Context) string {
	var category string
	langCategory := GetLanguageCategory(lang)
	if langCategory == JSON_LIKE {
		category = CategorizeJsonLikeSnippet(contents, llm, ctx)
	} else if langCategory == DRIVERS_MINUS_JS {
		category = CategorizeDriverLanguageSnippet(contents, llm, ctx)
	} else if langCategory == JAVASCRIPT || langCategory == TEXT {
		//category = CategorizeTextSnippet(contents, llm, ctx)
		category = CategorizeDriverLanguageSnippet(contents, llm, ctx)
	} else if langCategory == SHELL {
		category = CategorizeShellSnippet(contents, llm, ctx)
	} else {
		fmt.Printf("unknown language: %v\n", lang)
	}
	return category
}

func CategorizeJsonLikeSnippet(contents string, llm *ollama.LLM, ctx context.Context) string {
	// To tweak the prompt for accuracy, edit this question
	const questionTemplate = `I need to sort code examples into one of these categories:
	%s
	%s
	Use these definitions for each category to help categorize the code example:
	%s: An example object, typically represented in JSON, enumerating fields in a return object and their types. Typically includes an '_id' field and represents one or more example documents. Many pieces of JSON that look similar or repetitive in structure.
	%s: Example configuration object, typically represented in JSON or YAML, enumerating required/optional parameters and their types. If it shows an '_id' field, it is a return object, not a configuration object.
	Using these definitions, which category applies to this code example? Don't list an explanation, only list the category name.`
	question := fmt.Sprintf(questionTemplate,
		ExampleReturnObject,
		ExampleConfigurationObject,
		ExampleReturnObject,
		ExampleConfigurationObject,
	)
	template := prompts.NewPromptTemplate(
		`Use the following pieces of context to answer the question at the end.
			Context: {{.contents}}
			Question: {{.question}}`,
		[]string{"contents", "question"},
	)
	prompt, err := template.Format(map[string]any{
		"contents": contents,
		"question": question,
	})
	if err != nil {
		log.Fatalf("failed to create a prompt from the template: %q\n, %q\n, %q\n, %q\n", template, contents, question, err)
	}
	completion, err := llms.GenerateFromSinglePrompt(ctx, llm, prompt)
	if err != nil {
		log.Fatalf("failed to generate a response from the given prompt: %q", prompt)
	}
	return completion
}

func CategorizeShellSnippet(contents string, llm *ollama.LLM, ctx context.Context) string {
	// To tweak the prompt for accuracy, edit this question
	const questionTemplate = `I need to sort code examples into one of these categories:
	%s
	%s
	%s
	%s
	%s
	Use these definitions for each category to help categorize the code example:
	%s: One line or only a few lines of code that demonstrate popular command-line commands, such as 'docker ', 'go run', 'jq ', 'vi ', 'mkdir ', 'npm ', 'cd ' or other common command-line command invocations. If it starts with 'atlas ' it does not belong in this category - it is an Atlas CLI Command. If it starts with 'mongosh ' it does not belong in this category - it is a 'mongosh command'.
	%s: One-line or only a few lines of code that shows an Atlas CLI command. Must start with 'atlas ' at the beginning of the snippet or after a comment. If it does not start with 'atlas ' or has 'atlas ' anywhere that does not immediately follow a newline, it is not an Atlas CLI Command.
	%s: One-line or only a few lines of code that shows an mongosh function call, similar to 'db.methodName(arguments)' or 'collection.methodName(arguments)'.
	%s: Two variants: one is an example object, typically represented in JSON, enumerating fields in the return object and their types. Typically includes an '_id' field and represents one or more example documents. Many pieces of JSON that look similar or repetitive in structure. The second variant looks like text that has been logged to console, such as an error message or status information. May resemble "Backup completed." "Restore completed." or other short status messages.
	%s: Example object, typically represented in JSON or YAML, enumerating required/optional parameters and their types. If it shows an '_id' field, it is a return object, not a configuration object.
	Using these definitions, which category applies to this code example? Don't list an explanation, only list the category name.`
	question := fmt.Sprintf(questionTemplate,
		NonMongoCommand,
		AtlasCliCommand,
		MongoshCommand,
		ExampleReturnObject,
		ExampleConfigurationObject,
		NonMongoCommand,
		AtlasCliCommand,
		MongoshCommand,
		ExampleReturnObject,
		ExampleConfigurationObject,
	)
	template := prompts.NewPromptTemplate(
		`Use the following pieces of context to answer the question at the end.
			Context: {{.contents}}
			Question: {{.question}}`,
		[]string{"contents", "question"},
	)
	prompt, err := template.Format(map[string]any{
		"contents": contents,
		"question": question,
	})
	if err != nil {
		log.Fatalf("failed to create a prompt from the template: %q\n, %q\n, %q\n, %q\n", template, contents, question, err)
	}
	completion, err := llms.GenerateFromSinglePrompt(ctx, llm, prompt)
	if err != nil {
		log.Fatalf("failed to generate a response from the given prompt: %q", prompt)
	}
	return completion
}

func CategorizeTextSnippet(contents string, llm *ollama.LLM, ctx context.Context) string {
	// To tweak the prompt for accuracy, edit this question
	const questionTemplate = `I need to sort code examples into one of these categories:
	%s
	%s
	%s
	%s
	%s
	%s
	Use these definitions for each category to help categorize the code example:
	%s: One line or only a few lines of code that demonstrate popular command-line commands, such as 'docker ', 'go run', 'jq ', 'vi ', 'mkdir ', 'npm ', 'cd ' or other common command-line command invocations. If it starts with 'atlas ' it does not belong in this category - it is an Atlas CLI Command. If it starts with 'mongosh ' it does not belong in this category - it is a 'mongosh command'.
	%s: One-line or only a few lines of code that shows an Atlas CLI command. Must start with 'atlas ' at the beginning of the snippet or after a comment. If it does not start with 'atlas ' or has 'atlas ' anywhere that does not immediately follow a newline, it is not an Atlas CLI command.
	%s: One-line or only a few lines of code that shows an mongosh function call, similar to 'db.methodName(arguments)' or 'collection.methodName(arguments)'.	
	%s: One line that shows an API method signature, such as 'object.methodName(parameter1, parameter2)'. Code blocks showing 'main()' or other function declarations are not API method signatures - they are task-based usage. JSON blobs are not API method signatures.
	%s: Two variants: one is an example object, typically represented in JSON, enumerating fields in the return object and their types. Typically includes an '_id' field and represents one or more example documents. Many pieces of JSON that look similar or repetitive in structure. The second variant looks like text that has been logged to console, such as an error message or status information. May resemble "Backup completed." "Restore completed." or other short status messages.
	%s: Example object, typically represented in JSON or YAML, enumerating required/optional parameters and their types. If it shows an '_id' field, it is a return object, not a configuration object.
	Using these definitions, which category applies to this code example? Don't list an explanation, only list the category name.`
	question := fmt.Sprintf(questionTemplate,
		NonMongoCommand,
		AtlasCliCommand,
		MongoshCommand,
		ApiMethodSignature,
		ExampleReturnObject,
		ExampleConfigurationObject,
		NonMongoCommand,
		AtlasCliCommand,
		MongoshCommand,
		ApiMethodSignature,
		ExampleReturnObject,
		ExampleConfigurationObject,
	)
	template := prompts.NewPromptTemplate(
		`Use the following pieces of context to answer the question at the end.
			Context: {{.contents}}
			Question: {{.question}}`,
		[]string{"contents", "question"},
	)
	prompt, err := template.Format(map[string]any{
		"contents": contents,
		"question": question,
	})
	if err != nil {
		log.Fatalf("failed to create a prompt from the template: %q\n, %q\n, %q\n, %q\n", template, contents, question, err)
	}
	completion, err := llms.GenerateFromSinglePrompt(ctx, llm, prompt)
	if err != nil {
		log.Fatalf("failed to generate a response from the given prompt: %q", prompt)
	}
	return completion
}

func CategorizeDriverLanguageSnippet(contents string, llm *ollama.LLM, ctx context.Context) string {
	// To tweak the prompt for accuracy, edit this question
	const questionTemplate = `I need to sort code examples into one of these categories:
		%s
		%s
		Use these definitions for each category to help categorize the code example:
		%s: One line that shows an API method signature, such as 'object.methodName(parameter1, parameter2)'. Code blocks showing 'main()' or other function declarations are not API method signatures - they are task-based usage. JSON blobs are not API method signatures.
		%s: Longer code snippet that establishes parameters, performs basic set up code, and includes the larger context to demonstrate how to accomplish a task. If an example shows parameters but does not show initializing parameters, it does not fit this category. JSON blobs do not fit in this category.
		Using these definitions, which category applies to this code example? Don't list an explanation, only list the category name.`
	question := fmt.Sprintf(questionTemplate,
		ApiMethodSignature,
		UsageExample,
		ApiMethodSignature,
		UsageExample,
	)
	template := prompts.NewPromptTemplate(
		`Use the following pieces of context to answer the question at the end.
			Context: {{.contents}}
			Question: {{.question}}`,
		[]string{"contents", "question"},
	)
	prompt, err := template.Format(map[string]any{
		"contents": contents,
		"question": question,
	})
	if err != nil {
		log.Fatalf("failed to create a prompt from the template: %q\n, %q\n, %q\n, %q\n", template, contents, question, err)
	}
	completion, err := llms.GenerateFromSinglePrompt(ctx, llm, prompt)
	if err != nil {
		log.Fatalf("failed to generate a response from the given prompt: %q", prompt)
	}
	return completion
}
