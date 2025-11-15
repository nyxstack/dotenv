package dotenv

import (
	"fmt"
	"regexp"
	"strings"
)

// Tokenizer handles lexical analysis of .env content
type Tokenizer struct {
	content    string
	pos        int
	line       int
	col        int
	length     int
	exportMode bool // whether to handle "export KEY=value" syntax
}

// NewTokenizer creates a new tokenizer for the given content
func NewTokenizer(content string) *Tokenizer {
	return &Tokenizer{
		content:    content,
		pos:        0,
		line:       1,
		col:        1,
		length:     len(content),
		exportMode: true, // enable export handling by default
	}
}

// peek returns the character at the current position without advancing
func (t *Tokenizer) peek() byte {
	if t.pos >= t.length {
		return 0
	}
	return t.content[t.pos]
}

// peekNext returns the character at the next position without advancing
func (t *Tokenizer) peekNext() byte {
	if t.pos+1 >= t.length {
		return 0
	}
	return t.content[t.pos+1]
}

// advance moves to the next character
func (t *Tokenizer) advance() byte {
	if t.pos >= t.length {
		return 0
	}
	ch := t.content[t.pos]
	t.pos++
	if ch == '\n' {
		t.line++
		t.col = 1
	} else {
		t.col++
	}
	return ch
}

// skipWhitespace skips spaces and tabs (but not newlines)
func (t *Tokenizer) skipWhitespace() {
	for t.pos < t.length && (t.peek() == ' ' || t.peek() == '\t') {
		t.advance()
	}
}

// skipToNextLine skips to the next line
func (t *Tokenizer) skipToNextLine() {
	for t.pos < t.length && t.peek() != '\n' {
		t.advance()
	}
	if t.pos < t.length {
		t.advance() // consume the newline
	}
}

// isValidKeyChar checks if character is valid in a key name
func isValidKeyChar(ch byte) bool {
	return (ch >= 'A' && ch <= 'Z') || (ch >= 'a' && ch <= 'z') ||
		(ch >= '0' && ch <= '9') || ch == '_'
}

// isValidKeyStart checks if character is valid at the start of a key name
func isValidKeyStart(ch byte) bool {
	return (ch >= 'A' && ch <= 'Z') || (ch >= 'a' && ch <= 'z') || ch == '_'
}

// parseKey parses a key (identifier)
func (t *Tokenizer) parseKey() (string, error) {
	start := t.pos

	// First character must be letter or underscore
	if !isValidKeyStart(t.peek()) {
		return "", fmt.Errorf("invalid key name at line %d: keys must start with letter or underscore", t.line)
	}

	for t.pos < t.length && isValidKeyChar(t.peek()) {
		t.advance()
	}

	return t.content[start:t.pos], nil
}

// parseUnquotedValue parses an unquoted value until comment or newline
func (t *Tokenizer) parseUnquotedValue() (string, bool) {
	var result strings.Builder
	hasComment := false

	for t.pos < t.length {
		ch := t.peek()
		if ch == '\n' || ch == '\r' {
			break
		}
		if ch == '#' {
			hasComment = true
			break
		}
		result.WriteByte(t.advance())
	}

	// Only trim trailing whitespace, preserve leading whitespace
	value := strings.TrimRight(result.String(), " \t")
	return value, hasComment
} // parseQuotedValue parses a quoted value (single or double quotes)
func (t *Tokenizer) parseQuotedValue(quote byte) (string, error) {
	var result strings.Builder
	t.advance() // consume opening quote

	for t.pos < t.length {
		ch := t.peek()

		if ch == quote {
			t.advance() // consume closing quote
			return result.String(), nil
		}

		if ch == '\\' && quote == '"' {
			// Handle escapes only in double quotes
			t.advance() // consume backslash
			if t.pos >= t.length {
				return "", fmt.Errorf("unexpected end of file after escape at line %d", t.line)
			}

			escaped := t.advance()
			switch escaped {
			case 'n':
				result.WriteByte('\n')
			case 't':
				result.WriteByte('\t')
			case 'r':
				result.WriteByte('\r')
			case '\\':
				result.WriteByte('\\')
			case '"':
				result.WriteByte('"')
			case '\'':
				result.WriteByte('\'')
			default:
				// For unknown escapes, include both backslash and character
				result.WriteByte('\\')
				result.WriteByte(escaped)
			}
		} else {
			result.WriteByte(t.advance())
		}
	}

	return "", fmt.Errorf("unterminated quoted string at line %d", t.line)
}

// expandVariables expands ${VAR} and $VAR patterns in the value
func expandVariables(value string, env map[string]string) string {
	// First handle ${VAR} format
	varPattern := regexp.MustCompile(`\$\{([A-Za-z_][A-Za-z0-9_]*)\}`)
	result := varPattern.ReplaceAllStringFunc(value, func(match string) string {
		varName := match[2 : len(match)-1] // remove ${ and }
		if val, exists := env[varName]; exists {
			return val
		}
		return match // keep original if not found
	})

	// Then handle $VAR format (but not if already inside ${})
	simplePattern := regexp.MustCompile(`\$([A-Za-z_][A-Za-z0-9_]*)`)
	result = simplePattern.ReplaceAllStringFunc(result, func(match string) string {
		varName := match[1:] // remove $
		if val, exists := env[varName]; exists {
			return val
		}
		return match // keep original if not found
	})

	return result
}
