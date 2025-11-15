/*
Package dotenv provides a comprehensive .env file parser with full grammar support.

This package implements the same features as mature dotenv loaders in other languages,
including variable expansion, quoted strings, inline comments, and escape sequences.

The package is organized into several files:
- token.go: Token types and definitions
- tokenizer.go: Lexical analysis and tokenization
- parser.go: Parsing logic and state machine
- env.go: Main API functions for external users

Basic usage:

	env, err := dotenv.Load(".env")
	if err != nil {
		log.Fatal(err)
	}

	// Apply to current process
	err = dotenv.Apply(env)
	if err != nil {
		log.Fatal(err)
	}

For more control, use the Parser directly:

	parser := dotenv.NewParser(content)
	env, err := parser.Parse()
*/
package dotenv
