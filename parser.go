package dotenv

import (
	"fmt"
	"strings"
)

// Parser represents the .env parser with quote context tracking
type Parser struct {
	tokenizer *Tokenizer
}

// NewParser creates a new parser for the given content
func NewParser(content string) *Parser {
	return &Parser{
		tokenizer: NewTokenizer(content),
	}
}

// ParseState represents the current parsing state
type ParseState int

const (
	StateNormal ParseState = iota
	StateSingleQuoted
	StateDoubleQuoted
)

// LineResult holds the result of parsing a line
type LineResult struct {
	Key            string
	Value          string
	AllowExpansion bool
	Error          error
}

// ParseLine parses a single line and returns key, value, expansion flag, and any error
func (p *Parser) ParseLine() LineResult {
	// Skip leading whitespace
	p.tokenizer.skipWhitespace()

	// Check for empty line or comment
	if p.tokenizer.pos >= p.tokenizer.length ||
		p.tokenizer.peek() == '\n' ||
		p.tokenizer.peek() == '\r' ||
		p.tokenizer.peek() == '#' {
		p.tokenizer.skipToNextLine()
		return LineResult{Key: "", Value: "", AllowExpansion: false, Error: nil}
	}

	// Handle optional "export " prefix
	if p.tokenizer.exportMode && strings.HasPrefix(p.tokenizer.content[p.tokenizer.pos:], "export ") {
		p.tokenizer.pos += 7 // skip "export "
		p.tokenizer.col += 7
		p.tokenizer.skipWhitespace()
	}

	// Parse key
	key, err := p.tokenizer.parseKey()
	if err != nil {
		return LineResult{Key: "", Value: "", AllowExpansion: false, Error: err}
	}
	if key == "" {
		return LineResult{Key: "", Value: "", AllowExpansion: false,
			Error: fmt.Errorf("expected variable name at line %d", p.tokenizer.line)}
	}

	// Skip whitespace after key
	p.tokenizer.skipWhitespace()

	// Expect '=' assignment
	if p.tokenizer.peek() != '=' {
		return LineResult{Key: "", Value: "", AllowExpansion: false,
			Error: fmt.Errorf("expected '=' after variable name at line %d", p.tokenizer.line)}
	}
	p.tokenizer.advance() // consume '='

	// Skip whitespace after '='
	p.tokenizer.skipWhitespace()

	// Parse value and track quote type for variable expansion
	var value string
	var allowExpansion bool = true // default to allowing expansion

	ch := p.tokenizer.peek()
	if ch == '"' {
		// Double-quoted string - allow expansion
		value, err = p.tokenizer.parseQuotedValue('"')
		if err != nil {
			return LineResult{Key: "", Value: "", AllowExpansion: false, Error: err}
		}
		allowExpansion = true
	} else if ch == '\'' {
		// Single-quoted string - no variable expansion
		value, err = p.tokenizer.parseQuotedValue('\'')
		if err != nil {
			return LineResult{Key: "", Value: "", AllowExpansion: false, Error: err}
		}
		allowExpansion = false
	} else {
		// Unquoted value - allow expansion
		var hasComment bool
		value, hasComment = p.tokenizer.parseUnquotedValue()
		allowExpansion = true

		// Skip trailing comment if present
		if hasComment {
			p.tokenizer.skipToNextLine()
		}
	}

	// Skip any remaining whitespace and comments on the line
	p.tokenizer.skipWhitespace()
	if p.tokenizer.pos < p.tokenizer.length && p.tokenizer.peek() == '#' {
		p.tokenizer.skipToNextLine()
	}

	return LineResult{Key: key, Value: value, AllowExpansion: allowExpansion, Error: nil}
}

// ParseLineCompat provides backward compatibility with the old ParseLine signature
func (p *Parser) ParseLineCompat() (string, string, error) {
	result := p.ParseLine()
	return result.Key, result.Value, result.Error
}

// Parse parses the entire .env content and returns a map of environment variables
func (p *Parser) Parse() (map[string]string, error) {
	env := make(map[string]string)

	for p.tokenizer.pos < p.tokenizer.length {
		result := p.ParseLine()
		if result.Error != nil {
			return nil, result.Error
		}

		// Skip empty lines and comments
		if result.Key == "" {
			continue
		}

		value := result.Value

		// Expand variables only if expansion is allowed and value contains $
		if result.AllowExpansion && strings.Contains(value, "$") {
			value = expandVariables(value, env)
		}

		env[result.Key] = value
	}

	return env, nil
}
