package dotenv

import (
	"strings"
	"testing"
)

func TestBasicParsing(t *testing.T) {
	content := `KEY1=value1
KEY2=value2
KEY3=`
	parser := NewParser(content)
	env, err := parser.Parse()
	if err != nil {
		t.Fatalf("Parse failed: %v", err)
	}

	expected := map[string]string{
		"KEY1": "value1",
		"KEY2": "value2",
		"KEY3": "",
	}

	for k, v := range expected {
		if env[k] != v {
			t.Errorf("Expected %s=%s, got %s=%s", k, v, k, env[k])
		}
	}
}

func TestExportSyntax(t *testing.T) {
	content := `export KEY1=value1
export KEY2="value2"
KEY3=value3`
	parser := NewParser(content)
	env, err := parser.Parse()
	if err != nil {
		t.Fatalf("Parse failed: %v", err)
	}

	expected := map[string]string{
		"KEY1": "value1",
		"KEY2": "value2",
		"KEY3": "value3",
	}

	for k, v := range expected {
		if env[k] != v {
			t.Errorf("Expected %s=%s, got %s=%s", k, v, k, env[k])
		}
	}
}

func TestUnquotedValuesWithSpaces(t *testing.T) {
	content := `KEY1=some value with spaces
KEY2=value with trailing spaces   
KEY3=    value with leading spaces`
	parser := NewParser(content)
	env, err := parser.Parse()
	if err != nil {
		t.Fatalf("Parse failed: %v", err)
	}

	expected := map[string]string{
		"KEY1": "some value with spaces",
		"KEY2": "value with trailing spaces", // trailing spaces should be trimmed
		"KEY3": "value with leading spaces",  // leading spaces after = are skipped by whitespace parsing
	}

	for k, v := range expected {
		if env[k] != v {
			t.Errorf("Expected %s=%q, got %s=%q", k, v, k, env[k])
		}
	}
}

func TestQuotedStrings(t *testing.T) {
	content := `DOUBLE_QUOTED="value with spaces"
SINGLE_QUOTED='value with spaces'
DOUBLE_WITH_ESCAPES="line1\nline2\ttab"
SINGLE_NO_ESCAPES='line1\nline2\ttab'
EMPTY_DOUBLE=""
EMPTY_SINGLE=''`
	parser := NewParser(content)
	env, err := parser.Parse()
	if err != nil {
		t.Fatalf("Parse failed: %v", err)
	}

	expected := map[string]string{
		"DOUBLE_QUOTED":       "value with spaces",
		"SINGLE_QUOTED":       "value with spaces",
		"DOUBLE_WITH_ESCAPES": "line1\nline2\ttab",
		"SINGLE_NO_ESCAPES":   "line1\\nline2\\ttab", // escapes should be literal in single quotes
		"EMPTY_DOUBLE":        "",
		"EMPTY_SINGLE":        "",
	}

	for k, v := range expected {
		if env[k] != v {
			t.Errorf("Expected %s=%q, got %s=%q", k, v, k, env[k])
		}
	}
}

func TestInlineComments(t *testing.T) {
	content := `KEY1=value # this is a comment
KEY2="quoted value" # this is also a comment
KEY3='single quoted' # comment after single quotes
KEY4="string with # hash inside"
KEY5='string with # hash inside single quotes'
# Full line comment
KEY6=value_without_comment`
	parser := NewParser(content)
	env, err := parser.Parse()
	if err != nil {
		t.Fatalf("Parse failed: %v", err)
	}

	expected := map[string]string{
		"KEY1": "value",
		"KEY2": "quoted value",
		"KEY3": "single quoted",
		"KEY4": "string with # hash inside",
		"KEY5": "string with # hash inside single quotes",
		"KEY6": "value_without_comment",
	}

	for k, v := range expected {
		if env[k] != v {
			t.Errorf("Expected %s=%q, got %s=%q", k, v, k, env[k])
		}
	}
}

func TestEscapeSequences(t *testing.T) {
	content := `NEWLINE="line1\nline2"
TAB="col1\tcol2"
BACKSLASH="path\\to\\file"
QUOTE="say \"hello\""
UNKNOWN_ESCAPE="test\x"`
	parser := NewParser(content)
	env, err := parser.Parse()
	if err != nil {
		t.Fatalf("Parse failed: %v", err)
	}

	expected := map[string]string{
		"NEWLINE":        "line1\nline2",
		"TAB":            "col1\tcol2",
		"BACKSLASH":      "path\\to\\file",
		"QUOTE":          "say \"hello\"",
		"UNKNOWN_ESCAPE": "test\\x", // unknown escapes should preserve backslash
	}

	for k, v := range expected {
		if env[k] != v {
			t.Errorf("Expected %s=%q, got %s=%q", k, v, k, env[k])
		}
	}
}

func TestVariableExpansion(t *testing.T) {
	content := `BASE_DIR=/app
HOME_DIR=/home/user
PATH_ORIG=/usr/bin
PATH="${HOME_DIR}/bin:${PATH_ORIG}"
CONFIG_PATH="${HOME_DIR}/config"
NESTED="${PATH}:${BASE_DIR}/bin"
NO_EXPAND='$HOME_DIR/literal'`

	parser := NewParser(content)
	env, err := parser.Parse()
	if err != nil {
		t.Fatalf("Parse failed: %v", err)
	}

	expected := map[string]string{
		"BASE_DIR":    "/app",
		"HOME_DIR":    "/home/user",
		"PATH":        "/home/user/bin:/usr/bin",          // Should expand HOME_DIR and PATH_ORIG
		"CONFIG_PATH": "/home/user/config",                // Should expand HOME_DIR
		"NESTED":      "/home/user/bin:/usr/bin:/app/bin", // Should expand PATH and BASE_DIR
		"NO_EXPAND":   "$HOME_DIR/literal",                // Single quotes should prevent expansion
		"PATH_ORIG":   "/usr/bin",
	}

	for k, v := range expected {
		if env[k] != v {
			t.Errorf("Expected %s=%q, got %s=%q", k, v, k, env[k])
		}
	}
}

func TestEmptyValues(t *testing.T) {
	content := `EMPTY1=
EMPTY2=""
EMPTY3=''
WHITESPACE_ONLY=   `
	parser := NewParser(content)
	env, err := parser.Parse()
	if err != nil {
		t.Fatalf("Parse failed: %v", err)
	}

	expected := map[string]string{
		"EMPTY1":          "",
		"EMPTY2":          "",
		"EMPTY3":          "",
		"WHITESPACE_ONLY": "", // unquoted whitespace-only should be trimmed to empty
	}

	for k, v := range expected {
		if env[k] != v {
			t.Errorf("Expected %s=%q, got %s=%q", k, v, k, env[k])
		}
	}
}

func TestErrorCases(t *testing.T) {
	testCases := []struct {
		name    string
		content string
		wantErr bool
	}{
		{
			name:    "missing equals",
			content: "KEY_WITHOUT_EQUALS",
			wantErr: true,
		},
		{
			name:    "unterminated double quote",
			content: `KEY="unterminated`,
			wantErr: true,
		},
		{
			name:    "unterminated single quote",
			content: `KEY='unterminated`,
			wantErr: true,
		},
		{
			name:    "invalid key",
			content: "123INVALID=value",
			wantErr: true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			parser := NewParser(tc.content)
			_, err := parser.Parse()
			if tc.wantErr && err == nil {
				t.Errorf("Expected error but got none")
			}
			if !tc.wantErr && err != nil {
				t.Errorf("Unexpected error: %v", err)
			}
		})
	}
}

func TestComplexRealWorldExample(t *testing.T) {
	content := `# Database configuration
DB_HOST=localhost
DB_PORT=5432
DB_NAME=myapp
DB_USER=user
DB_PASSWORD="p@ssw0rd with spaces"

# API Configuration
export API_KEY="abc123def456"
API_URL="https://api.example.com/v1"
API_TIMEOUT=30

# File paths
LOG_DIR=/var/log/myapp
CONFIG_FILE="${LOG_DIR}/config.json"
BACKUP_PATH='$HOME/backups' # Single quotes = literal

# Feature flags
ENABLE_DEBUG=true
ENABLE_CACHE=false

# Empty and commented
OPTIONAL_KEY=
# COMMENTED_OUT=value

# Multi-line simulation with escapes
MESSAGE="Line 1\nLine 2\tTabbed"
`

	parser := NewParser(content)
	env, err := parser.Parse()
	if err != nil {
		t.Fatalf("Parse failed: %v", err)
	}

	// Verify some key entries
	expectedValues := map[string]string{
		"DB_HOST":      "localhost",
		"DB_PASSWORD":  "p@ssw0rd with spaces",
		"API_KEY":      "abc123def456",
		"CONFIG_FILE":  "/var/log/myapp/config.json",
		"BACKUP_PATH":  "$HOME/backups",
		"ENABLE_DEBUG": "true",
		"OPTIONAL_KEY": "",
		"MESSAGE":      "Line 1\nLine 2\tTabbed",
	}

	for k, v := range expectedValues {
		if env[k] != v {
			t.Errorf("Expected %s=%q, got %s=%q", k, v, k, env[k])
		}
	}

	// Verify commented out key doesn't exist
	if _, exists := env["COMMENTED_OUT"]; exists {
		t.Errorf("COMMENTED_OUT should not exist in parsed environment")
	}
}

func TestLoadFromReader(t *testing.T) {
	content := "KEY1=value1\nKEY2=value2"
	reader := strings.NewReader(content)

	env, err := LoadFromReader(reader)
	if err != nil {
		t.Fatalf("LoadFromReader failed: %v", err)
	}

	expected := map[string]string{
		"KEY1": "value1",
		"KEY2": "value2",
	}

	for k, v := range expected {
		if env[k] != v {
			t.Errorf("Expected %s=%s, got %s=%s", k, v, k, env[k])
		}
	}
}
