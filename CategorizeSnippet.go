package main

import (
	"context"
	"fmt"
	"github.com/tmc/langchaingo/llms"
	"github.com/tmc/langchaingo/llms/ollama"
	"github.com/tmc/langchaingo/prompts"
	"log"
)

func ProcessSnippet(contents string, lang string, llm *ollama.LLM, ctx context.Context) (string, int) {
	var category string
	attemptCounter := 1
	validCategories := []string{AtlasCliCommand, ApiMethodSignature, ExampleReturnObject, ExampleConfigurationObject, UsageExample}

	for attemptCounter <= 3 {
		category = LLMAssignCategory(contents, lang, llm, ctx)
		if containsString(validCategories, category) {
			return category, attemptCounter
		} else {
			attemptCounter++
			category = LLMAssignCategory(contents, lang, llm, ctx)
		}
	}
	return "Uncategorized", attemptCounter
}

func LLMAssignCategory(contents string, lang string, llm *ollama.LLM, ctx context.Context) string {
	var category string
	textShellOrJson := []string{JSON, SHELL, TEXT, XML, YAML}
	driversLanguagesMinusJS := []string{C, CPP, CSHARP, GO, JAVA, KOTLIN, PHP, PYTHON, RUBY, RUST, SCALA, SWIFT, TYPESCRIPT}
	if containsString(textShellOrJson, lang) {
		category = CategorizeTextShellOrJsonSnippet(string(contents), llm, ctx)
	} else if containsString(driversLanguagesMinusJS, lang) {
		category = CategorizeDriverLanguageSnippet(string(contents), llm, ctx)
	} else if lang == JAVASCRIPT {
		category = CategorizeJavaScriptSnippet(string(contents), llm, ctx)
	} else {
		fmt.Printf("unknown language: %v\n", lang)
	}
	return category
}

// CategorizeSnippet uses the LLM categorize the string passed into the func based on the prompt defined here
func CategorizeSnippet(contents string, llm *ollama.LLM, ctx context.Context) string {
	// To tweak the prompt for accuracy, edit this question
	const question = `I need to sort code examples into one of five categories. The five category options are:

		Atlas CLI Command
		API method signature
		Example return object
		Example configuration object
		Task-based usage

		Use these definitions for each category to help categorize the code example:

		Atlas CLI Command: One-line or only a few lines of code that shows an Atlas CLI command, typically starting with 'atlas'. If it is multiple lines with multiple 'atlas' blocks, it belongs in the Task-based usage category.
		API method signature: One-line or only a few lines of code that shows an API method signature. Code blocks showing 'main()' or other function declarations are not API method signatures - they are task-based usage. JSON blobs are not API method signatures.
		Example return object: Example object, typically represented in JSON, enumerating fields in the return object and their types. Typically includes an '_id' field and represents one or more example documents. Many pieces of JSON that look similar or repetitive in structure probably represent an example return object.
		Example configuration object: Example object, typically represented in JSON or YAML, enumerating required/optional parameters and their types.
		Task-based usage: Longer code snippet that establishes parameters, performs basic set up code, and includes the larger context to demonstrate how to accomplish a task. If an example shows parameters but does not show initializing parameters, it does not fit this category. JSON blobs do not fit in this category.
	
		Using these definitions, which category applies to this code example? Don't list an explanation, only list the category name.`
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

func CategorizeTextShellOrJsonSnippet(contents string, llm *ollama.LLM, ctx context.Context) string {
	// To tweak the prompt for accuracy, edit this question
	const questionTemplate = `I need to sort code examples into one of these categories:
	%s
	%s
	%s
	%s
	Use these definitions for each category to help categorize the code example:
	%s: One-line or only a few lines of code that shows an Atlas CLI command, typically starting with 'atlas'. If it is multiple lines with multiple 'atlas' blocks, it belongs in the Task-based usage category.
	%s: One-line or only a few lines of code that shows an API method signature. Code blocks showing 'main()' or other function declarations are not API method signatures - they are task-based usage. JSON blobs are not API method signatures.
	%s: Example object, typically represented in JSON, enumerating fields in the return object and their types. Typically includes an '_id' field and represents one or more example documents. Many pieces of JSON that look similar or repetitive in structure probably represent an example return object.
	%s: Example object, typically represented in JSON or YAML, enumerating required/optional parameters and their types.
	Using these definitions, which category applies to this code example? Don't list an explanation, only list the category name.`
	question := fmt.Sprintf(questionTemplate,
		AtlasCliCommand,
		ApiMethodSignature,
		ExampleReturnObject,
		ExampleConfigurationObject,
		AtlasCliCommand,
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

func CategorizeJavaScriptSnippet(contents string, llm *ollama.LLM, ctx context.Context) string {
	// To tweak the prompt for accuracy, edit this question
	const questionTemplate = `I need to sort code examples into one of these categories:
	%s
	%s
	%s
	%s
	%s
	Use these definitions for each category to help categorize the code example:
	%s: One-line or only a few lines of code that shows an Atlas CLI command, typically starting with 'atlas'. If it is multiple lines with multiple 'atlas' blocks, it belongs in the Task-based usage category.
	%s: One-line or only a few lines of code that shows an API method signature. Code blocks showing 'main()' or other function declarations are not API method signatures - they are task-based usage. JSON blobs are not API method signatures.
	%s: Example object, typically represented in JSON, enumerating fields in the return object and their types. Typically includes an '_id' field and represents one or more example documents. Many pieces of JSON that look similar or repetitive in structure probably represent an example return object.
	%s: Example object, typically represented in JSON or YAML, enumerating required/optional parameters and their types.
	%s: Longer code snippet that establishes parameters, performs basic set up code, and includes the larger context to demonstrate how to accomplish a task. If an example shows parameters but does not show initializing parameters, it does not fit this category. JSON blobs do not fit in this category.	
	Using these definitions, which category applies to this code example? Don't list an explanation, only list the category name.`
	question := fmt.Sprintf(questionTemplate,
		AtlasCliCommand,
		ApiMethodSignature,
		ExampleReturnObject,
		ExampleConfigurationObject,
		UsageExample,
		AtlasCliCommand,
		ApiMethodSignature,
		ExampleReturnObject,
		ExampleConfigurationObject,
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

func CategorizeDriverLanguageSnippet(contents string, llm *ollama.LLM, ctx context.Context) string {
	// To tweak the prompt for accuracy, edit this question
	const questionTemplate = `I need to sort code examples into one of these categories:
		%s
		%s
		%s
		Use these definitions for each category to help categorize the code example:
		%s: One-line or only a few lines of code that shows an API method signature. Code blocks showing 'main()' or other function declarations are not API method signatures - they are task-based usage. JSON blobs are not API method signatures.
		%s: Example object, typically represented in JSON or YAML, enumerating required/optional parameters and their types.
		%s: Longer code snippet that establishes parameters, performs basic set up code, and includes the larger context to demonstrate how to accomplish a task. If an example shows parameters but does not show initializing parameters, it does not fit this category. JSON blobs do not fit in this category.
		Using these definitions, which category applies to this code example? Don't list an explanation, only list the category name.`
	question := fmt.Sprintf(questionTemplate,
		ApiMethodSignature,
		ExampleConfigurationObject,
		UsageExample,
		ApiMethodSignature,
		ExampleConfigurationObject,
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
