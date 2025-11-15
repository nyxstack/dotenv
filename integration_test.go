package dotenv

import (
	"strings"
	"testing"
)

func TestIntegration(t *testing.T) {
	// Test the complete workflow with a realistic .env file
	content := `# Production configuration
export NODE_ENV=production
export DEBUG=false

# Database settings  
DB_HOST=localhost
DB_PORT=5432
DB_NAME="my_app_prod"
DB_URL="postgresql://${DB_HOST}:${DB_PORT}/${DB_NAME}"

# API Configuration
API_BASE_URL="https://api.example.com"
API_KEY='secret-key-$NODE_ENV'  # Single quotes prevent expansion
API_TIMEOUT=30 # seconds

# File paths with variable expansion
LOG_DIR=/var/log/myapp
CONFIG_FILE="${LOG_DIR}/app.conf"
BACKUP_PATH="${LOG_DIR}/backups"

# Feature flags
ENABLE_CACHE=true
CACHE_TTL=3600

# Empty optional settings
OPTIONAL_PLUGIN=
OPTIONAL_CONFIG=""

# Values with special characters
SPECIAL_CHARS="!@#$%^&*()_+-={}[]|\\:;\"'<>,.?/"
ESCAPED_QUOTES="Say \"Hello\" and 'Goodbye'"
MULTILINE_TEXT="Line 1\nLine 2\tTabbed\nLine 3"
`

	parser := NewParser(content)
	env, err := parser.Parse()
	if err != nil {
		t.Fatalf("Parse failed: %v", err)
	}

	// Verify basic parsing
	if env["NODE_ENV"] != "production" {
		t.Errorf("Expected NODE_ENV=production, got %s", env["NODE_ENV"])
	}

	if env["DB_NAME"] != "my_app_prod" {
		t.Errorf("Expected DB_NAME=my_app_prod, got %s", env["DB_NAME"])
	}

	// Verify variable expansion
	expectedDBUrl := "postgresql://localhost:5432/my_app_prod"
	if env["DB_URL"] != expectedDBUrl {
		t.Errorf("Expected DB_URL=%s, got %s", expectedDBUrl, env["DB_URL"])
	}

	expectedConfigFile := "/var/log/myapp/app.conf"
	if env["CONFIG_FILE"] != expectedConfigFile {
		t.Errorf("Expected CONFIG_FILE=%s, got %s", expectedConfigFile, env["CONFIG_FILE"])
	}

	// Verify single quotes prevent expansion
	expectedAPIKey := "secret-key-$NODE_ENV"
	if env["API_KEY"] != expectedAPIKey {
		t.Errorf("Expected API_KEY=%s (no expansion), got %s", expectedAPIKey, env["API_KEY"])
	}

	// Verify inline comments are stripped
	if env["API_TIMEOUT"] != "30" {
		t.Errorf("Expected API_TIMEOUT=30 (comment stripped), got %s", env["API_TIMEOUT"])
	}

	// Verify empty values
	if env["OPTIONAL_PLUGIN"] != "" {
		t.Errorf("Expected OPTIONAL_PLUGIN to be empty, got %s", env["OPTIONAL_PLUGIN"])
	}

	if env["OPTIONAL_CONFIG"] != "" {
		t.Errorf("Expected OPTIONAL_CONFIG to be empty, got %s", env["OPTIONAL_CONFIG"])
	}

	// Verify escape sequences
	expectedMultiline := "Line 1\nLine 2\tTabbed\nLine 3"
	if env["MULTILINE_TEXT"] != expectedMultiline {
		t.Errorf("Expected MULTILINE_TEXT=%q, got %q", expectedMultiline, env["MULTILINE_TEXT"])
	}

	// Verify special characters are preserved
	expectedSpecial := "!@#$%^&*()_+-={}[]|\\:;\"'<>,.?/"
	if env["SPECIAL_CHARS"] != expectedSpecial {
		t.Errorf("Expected SPECIAL_CHARS=%q, got %q", expectedSpecial, env["SPECIAL_CHARS"])
	}

	// Test LoadFromReader
	reader := strings.NewReader("TEST_KEY=test_value\nTEST_KEY2=\"quoted value\"")
	readerEnv, err := LoadFromReader(reader)
	if err != nil {
		t.Fatalf("LoadFromReader failed: %v", err)
	}

	if readerEnv["TEST_KEY"] != "test_value" {
		t.Errorf("Expected TEST_KEY=test_value, got %s", readerEnv["TEST_KEY"])
	}

	if readerEnv["TEST_KEY2"] != "quoted value" {
		t.Errorf("Expected TEST_KEY2=quoted value, got %s", readerEnv["TEST_KEY2"])
	}

	t.Logf("Integration test passed! Parsed %d environment variables", len(env))
}
