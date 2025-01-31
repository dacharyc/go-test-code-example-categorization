package main

import (
	"context"
	"github.com/tmc/langchaingo/llms"
	"github.com/tmc/langchaingo/llms/ollama"
	"github.com/tmc/langchaingo/prompts"
	"log"
)

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
	const question = `I need to sort code examples into one of these categories:

		Atlas CLI Command
		API method signature
		Example return object
		Example configuration object

		These categories are defined as:

		Atlas CLI Command: One-line or only a few lines of code that shows an Atlas CLI command, typically starting with 'atlas'. If it is multiple lines with multiple 'atlas' blocks, it belongs in the Task-based usage category.
		API method signature: One-line or only a few lines of code that shows an API method signature. Code blocks showing 'main()' or other function declarations are not API method signatures - they are task-based usage. JSON blobs are not API method signatures.
		Example return object: Example object, typically represented in JSON, enumerating fields in the return object and their types. Typically includes an '_id' field and represents one or more example documents. Many pieces of JSON that look similar or repetitive in structure probably represent an example return object.
		Example configuration object: Example object, typically represented in JSON or YAML, enumerating required/optional parameters and their types.
	
		Using these definitions, which category name from [Atlas CLI Command, API method signature, Example return object, Example configuration object] applies to this code example? Don't list an explanation, only list the category name.`
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
	const question = `I need to sort code examples into one of these categories:

		Atlas CLI Command
		API method signature
		Example return object
		Example configuration object
		Task-based usage

		These categories are defined as:

		Atlas CLI Command: One-line or only a few lines of code that shows an Atlas CLI command, typically starting with 'atlas'. If it is multiple lines with multiple 'atlas' blocks, it belongs in the Task-based usage category.
		API method signature: One-line or only a few lines of code that shows an API method signature. Code blocks showing 'main()' or other function declarations are not API method signatures - they are task-based usage. JSON blobs are not API method signatures.
		Example return object: Example object, typically represented in JSON, enumerating fields in the return object and their types. Typically includes an '_id' field and represents one or more example documents. Many pieces of JSON that look similar or repetitive in structure probably represent an example return object.
		Example configuration object: Example object, typically represented in JSON or YAML, enumerating required/optional parameters and their types.
		Task-based usage: Longer code snippet that establishes parameters, performs basic set up code, and includes the larger context to demonstrate how to accomplish a task. If an example shows parameters but does not show initializing parameters, it does not fit this category. JSON blobs do not fit in this category.
	
		Using these definitions, which category name from [Atlas CLI Command, API method signature, Example return object, Example configuration object, Task-based usage] applies to this code example? Don't list an explanation, only list the category name.`
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
	const question = `I need to sort code examples into one of these categories:

		API method signature
		Example configuration object
		Task-based usage

		These categories are defined as:

		API method signature: One-line or only a few lines of code that shows an API method signature. Code blocks showing 'main()' or other function declarations are not API method signatures - they are task-based usage. JSON blobs are not API method signatures.
		Example configuration object: Example object, typically represented in JSON or YAML, enumerating required/optional parameters and their types.
		Task-based usage: Longer code snippet that establishes parameters, performs basic set up code, and includes the larger context to demonstrate how to accomplish a task. If an example shows parameters but does not show initializing parameters, it does not fit this category. JSON blobs do not fit in this category.
	
		Using these definitions, which category name from [API method signature, Example configuration object, Task-based usage] applies to this code example? Don't list an explanation, only list the category name.`
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
