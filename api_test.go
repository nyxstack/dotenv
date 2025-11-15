package dotenv

import (
	"os"
	"testing"
	"time"
)

// Test struct for unmarshaling
type Config struct {
	DatabaseURL    string        `env:"DATABASE_URL"`
	Port           int           `env:"PORT,default=8080"`
	Debug          bool          `env:"DEBUG,default=false"`
	Timeout        time.Duration `env:"TIMEOUT,default=30s"`
	MaxConnections int64         `env:"MAX_CONNECTIONS,required"`
	Features       []string      `env:"FEATURES"`
	APIKey         string        `env:"API_KEY,required"`
	LogLevel       string        `env:"LOG_LEVEL,default=info"`
	FloatValue     float64       `env:"FLOAT_VALUE,default=3.14"`
	UintValue      uint32        `env:"UINT_VALUE,default=100"`
	unexported     string        `env:"UNEXPORTED"` // Should be ignored
}

func TestUnmarshal(t *testing.T) {
	// Set up test environment variables
	testVars := map[string]string{
		"DATABASE_URL":    "postgresql://localhost:5432/test",
		"DEBUG":           "true",
		"TIMEOUT":         "45s",
		"MAX_CONNECTIONS": "150",
		"FEATURES":        "auth,logging,metrics",
		"API_KEY":         "secret123",
		"FLOAT_VALUE":     "2.718",
		"UINT_VALUE":      "200",
	}

	for key, value := range testVars {
		os.Setenv(key, value)
		defer os.Unsetenv(key)
	}

	var config Config
	err := Unmarshal(&config)
	if err != nil {
		t.Fatalf("Unmarshal failed: %v", err)
	}

	// Verify values
	if config.DatabaseURL != "postgresql://localhost:5432/test" {
		t.Errorf("Expected DatabaseURL=%s, got %s", "postgresql://localhost:5432/test", config.DatabaseURL)
	}

	if config.Port != 8080 { // Should use default
		t.Errorf("Expected Port=%d, got %d", 8080, config.Port)
	}

	if config.Debug != true {
		t.Errorf("Expected Debug=%t, got %t", true, config.Debug)
	}

	if config.Timeout != 45*time.Second {
		t.Errorf("Expected Timeout=%v, got %v", 45*time.Second, config.Timeout)
	}

	if config.MaxConnections != 150 {
		t.Errorf("Expected MaxConnections=%d, got %d", 150, config.MaxConnections)
	}

	expectedFeatures := []string{"auth", "logging", "metrics"}
	if len(config.Features) != len(expectedFeatures) {
		t.Errorf("Expected Features length=%d, got %d", len(expectedFeatures), len(config.Features))
	}
	for i, feature := range expectedFeatures {
		if config.Features[i] != feature {
			t.Errorf("Expected Features[%d]=%s, got %s", i, feature, config.Features[i])
		}
	}

	if config.APIKey != "secret123" {
		t.Errorf("Expected APIKey=%s, got %s", "secret123", config.APIKey)
	}

	if config.LogLevel != "info" { // Should use default
		t.Errorf("Expected LogLevel=%s, got %s", "info", config.LogLevel)
	}

	if config.FloatValue != 2.718 {
		t.Errorf("Expected FloatValue=%f, got %f", 2.718, config.FloatValue)
	}

	if config.UintValue != 200 {
		t.Errorf("Expected UintValue=%d, got %d", uint32(200), config.UintValue)
	}
}

func TestUnmarshalWithPrefix(t *testing.T) {
	// Set up test environment variables with prefix
	testVars := map[string]string{
		"APP_DATABASE_URL":    "postgresql://localhost:5432/app",
		"APP_API_KEY":         "app_secret456",
		"APP_MAX_CONNECTIONS": "75",
	}

	for key, value := range testVars {
		os.Setenv(key, value)
		defer os.Unsetenv(key)
	}

	var config Config
	err := UnmarshalWithPrefix(&config, "APP_")
	if err != nil {
		t.Fatalf("UnmarshalWithPrefix failed: %v", err)
	}

	if config.DatabaseURL != "postgresql://localhost:5432/app" {
		t.Errorf("Expected DatabaseURL=%s, got %s", "postgresql://localhost:5432/app", config.DatabaseURL)
	}

	if config.APIKey != "app_secret456" {
		t.Errorf("Expected APIKey=%s, got %s", "app_secret456", config.APIKey)
	}

	if config.MaxConnections != 75 {
		t.Errorf("Expected MaxConnections=%d, got %d", 75, config.MaxConnections)
	}
}

func TestUnmarshalRequired(t *testing.T) {
	var config Config

	// Should fail because required fields are missing
	err := Unmarshal(&config)
	if err == nil {
		t.Fatal("Expected error for missing required fields")
	}

	// Set required fields
	os.Setenv("MAX_CONNECTIONS", "100")
	os.Setenv("API_KEY", "test_key")
	defer func() {
		os.Unsetenv("MAX_CONNECTIONS")
		os.Unsetenv("API_KEY")
	}()

	// Should succeed now
	err = Unmarshal(&config)
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}
}

func TestEnvTypedFunctions(t *testing.T) {
	// Test string
	os.Setenv("TEST_STRING", "hello world")
	defer os.Unsetenv("TEST_STRING")

	if Env("TEST_STRING") != "hello world" {
		t.Errorf("Expected 'hello world', got '%s'", Env("TEST_STRING"))
	}

	if Env("NONEXISTENT", "default") != "default" {
		t.Errorf("Expected 'default', got '%s'", Env("NONEXISTENT", "default"))
	}

	// Test int
	os.Setenv("TEST_INT", "42")
	defer os.Unsetenv("TEST_INT")

	if EnvInt("TEST_INT") != 42 {
		t.Errorf("Expected 42, got %d", EnvInt("TEST_INT"))
	}

	if EnvInt("NONEXISTENT_INT", 100) != 100 {
		t.Errorf("Expected 100, got %d", EnvInt("NONEXISTENT_INT", 100))
	}

	// Test int64
	os.Setenv("TEST_INT64", "9223372036854775807")
	defer os.Unsetenv("TEST_INT64")

	if EnvInt64("TEST_INT64") != 9223372036854775807 {
		t.Errorf("Expected 9223372036854775807, got %d", EnvInt64("TEST_INT64"))
	}

	// Test bool
	os.Setenv("TEST_BOOL", "true")
	defer os.Unsetenv("TEST_BOOL")

	if EnvBool("TEST_BOOL") != true {
		t.Errorf("Expected true, got %t", EnvBool("TEST_BOOL"))
	}

	if EnvBool("NONEXISTENT_BOOL", false) != false {
		t.Errorf("Expected false, got %t", EnvBool("NONEXISTENT_BOOL", false))
	}

	// Test float64
	os.Setenv("TEST_FLOAT", "3.14159")
	defer os.Unsetenv("TEST_FLOAT")

	if EnvFloat64("TEST_FLOAT") != 3.14159 {
		t.Errorf("Expected 3.14159, got %f", EnvFloat64("TEST_FLOAT"))
	}

	// Test duration
	os.Setenv("TEST_DURATION", "5m30s")
	defer os.Unsetenv("TEST_DURATION")

	expected := 5*time.Minute + 30*time.Second
	if EnvDuration("TEST_DURATION") != expected {
		t.Errorf("Expected %v, got %v", expected, EnvDuration("TEST_DURATION"))
	}

	// Test uint types
	os.Setenv("TEST_UINT32", "4294967295")
	defer os.Unsetenv("TEST_UINT32")

	if EnvUint32("TEST_UINT32") != 4294967295 {
		t.Errorf("Expected 4294967295, got %d", EnvUint32("TEST_UINT32"))
	}
}

func TestEnvSetUnset(t *testing.T) {
	key := "TEST_SET_UNSET"
	value := "test_value"

	// Test SetEnv
	err := SetEnv(key, value)
	if err != nil {
		t.Fatalf("SetEnv failed: %v", err)
	}

	// Verify it was set
	if os.Getenv(key) != value {
		t.Errorf("Expected %s, got %s", value, os.Getenv(key))
	}

	// Test HasEnv
	if !HasEnv(key) {
		t.Error("Expected HasEnv to return true")
	}

	// Test UnsetEnv
	err = UnsetEnv(key)
	if err != nil {
		t.Fatalf("UnsetEnv failed: %v", err)
	}

	// Verify it was unset
	if HasEnv(key) {
		t.Error("Expected HasEnv to return false after unset")
	}
}

func TestEnvErrorHandling(t *testing.T) {
	// Test invalid int
	os.Setenv("INVALID_INT", "not_a_number")
	defer os.Unsetenv("INVALID_INT")

	if EnvInt("INVALID_INT", 999) != 999 {
		t.Errorf("Expected default value 999 for invalid int, got %d", EnvInt("INVALID_INT", 999))
	}

	// Test invalid bool
	os.Setenv("INVALID_BOOL", "not_a_bool")
	defer os.Unsetenv("INVALID_BOOL")

	if EnvBool("INVALID_BOOL", true) != true {
		t.Errorf("Expected default value true for invalid bool, got %t", EnvBool("INVALID_BOOL", true))
	}

	// Test invalid duration
	os.Setenv("INVALID_DURATION", "not_a_duration")
	defer os.Unsetenv("INVALID_DURATION")

	defaultDuration := 10 * time.Second
	if EnvDuration("INVALID_DURATION", defaultDuration) != defaultDuration {
		t.Errorf("Expected default value %v for invalid duration, got %v",
			defaultDuration, EnvDuration("INVALID_DURATION", defaultDuration))
	}
}

func TestMarshal(t *testing.T) {
	config := Config{
		DatabaseURL:    "postgresql://localhost:5432/test",
		Port:           3000,
		Debug:          true,
		Timeout:        45 * time.Second,
		MaxConnections: 150,
		Features:       []string{"auth", "logging", "metrics"},
		APIKey:         "secret123",
		LogLevel:       "debug",
		FloatValue:     2.718,
		UintValue:      200,
	}

	env, err := Marshal(&config)
	if err != nil {
		t.Fatalf("Marshal failed: %v", err)
	}

	// Verify marshaled values
	expected := map[string]string{
		"DATABASE_URL":    "postgresql://localhost:5432/test",
		"PORT":            "3000",
		"DEBUG":           "true",
		"TIMEOUT":         "45s",
		"MAX_CONNECTIONS": "150",
		"FEATURES":        "auth,logging,metrics",
		"API_KEY":         "secret123",
		"LOG_LEVEL":       "debug",
		"FLOAT_VALUE":     "2.718",
		"UINT_VALUE":      "200",
	}

	for key, expectedValue := range expected {
		if env[key] != expectedValue {
			t.Errorf("Expected %s=%s, got %s=%s", key, expectedValue, key, env[key])
		}
	}
}

func TestMarshalWithPrefix(t *testing.T) {
	config := Config{
		DatabaseURL: "postgresql://localhost:5432/app",
		APIKey:      "app_secret456",
		Port:        8080,
	}

	env, err := MarshalWithPrefix(&config, "APP_")
	if err != nil {
		t.Fatalf("MarshalWithPrefix failed: %v", err)
	}

	expected := map[string]string{
		"APP_DATABASE_URL": "postgresql://localhost:5432/app",
		"APP_API_KEY":      "app_secret456",
		"APP_PORT":         "8080",
	}

	for key, expectedValue := range expected {
		if env[key] != expectedValue {
			t.Errorf("Expected %s=%s, got %s=%s", key, expectedValue, key, env[key])
		}
	}
}

func TestMarshalToFile(t *testing.T) {
	config := Config{
		DatabaseURL:    "postgresql://localhost:5432/test",
		Port:           8080,
		Debug:          true,
		APIKey:         "test_key_123",
		LogLevel:       "info",
		Features:       []string{"auth", "cache"},
		FloatValue:     3.14159,
		MaxConnections: 100,
	}

	filename := "test_output.env"
	defer os.Remove(filename)

	err := MarshalToFile(filename, &config)
	if err != nil {
		t.Fatalf("MarshalToFile failed: %v", err)
	}

	// Read back and verify
	loadedEnv, err := Load(filename)
	if err != nil {
		t.Fatalf("Failed to load marshaled file: %v", err)
	}

	// Verify key values
	if loadedEnv["DATABASE_URL"] != "postgresql://localhost:5432/test" {
		t.Errorf("Expected DATABASE_URL=postgresql://localhost:5432/test, got %s", loadedEnv["DATABASE_URL"])
	}

	if loadedEnv["PORT"] != "8080" {
		t.Errorf("Expected PORT=8080, got %s", loadedEnv["PORT"])
	}

	if loadedEnv["DEBUG"] != "true" {
		t.Errorf("Expected DEBUG=true, got %s", loadedEnv["DEBUG"])
	}

	if loadedEnv["FEATURES"] != "auth,cache" {
		t.Errorf("Expected FEATURES=auth,cache, got %s", loadedEnv["FEATURES"])
	}
}

func TestMarshalUnmarshalRoundTrip(t *testing.T) {
	original := Config{
		DatabaseURL:    "postgresql://user:pass@localhost:5432/db",
		Port:           9000,
		Debug:          false,
		Timeout:        2 * time.Minute,
		MaxConnections: 75,
		Features:       []string{"feature1", "feature2", "feature3"},
		APIKey:         "round_trip_key",
		LogLevel:       "warn",
		FloatValue:     1.23456,
		UintValue:      500,
	}

	// Marshal to map
	envMap, err := Marshal(&original)
	if err != nil {
		t.Fatalf("Marshal failed: %v", err)
	}

	// Apply to environment
	for key, value := range envMap {
		os.Setenv(key, value)
		defer os.Unsetenv(key)
	}

	// Unmarshal back to struct
	var restored Config
	err = Unmarshal(&restored)
	if err != nil {
		t.Fatalf("Unmarshal failed: %v", err)
	}

	// Compare structs
	if restored.DatabaseURL != original.DatabaseURL {
		t.Errorf("DatabaseURL mismatch: expected %s, got %s", original.DatabaseURL, restored.DatabaseURL)
	}

	if restored.Port != original.Port {
		t.Errorf("Port mismatch: expected %d, got %d", original.Port, restored.Port)
	}

	if restored.Debug != original.Debug {
		t.Errorf("Debug mismatch: expected %t, got %t", original.Debug, restored.Debug)
	}

	if restored.Timeout != original.Timeout {
		t.Errorf("Timeout mismatch: expected %v, got %v", original.Timeout, restored.Timeout)
	}

	if len(restored.Features) != len(original.Features) {
		t.Errorf("Features length mismatch: expected %d, got %d", len(original.Features), len(restored.Features))
	} else {
		for i, feature := range original.Features {
			if restored.Features[i] != feature {
				t.Errorf("Feature[%d] mismatch: expected %s, got %s", i, feature, restored.Features[i])
			}
		}
	}
}

func TestWriteEnvFile(t *testing.T) {
	env := map[string]string{
		"SIMPLE_KEY":      "simple_value",
		"KEY_WITH_SPACES": "value with spaces",
		"KEY_WITH_QUOTES": "value with \"quotes\"",
		"MULTILINE_KEY":   "line1\nline2\ttab",
		"EMPTY_KEY":       "",
		"SPECIAL_CHARS":   "!@#$%^&*()_+-={}[]|\\:;\"'<>,.?/",
	}

	filename := "test_write.env"
	defer os.Remove(filename)

	err := WriteEnvFile(filename, env)
	if err != nil {
		t.Fatalf("WriteEnvFile failed: %v", err)
	}

	// Read back and parse
	loadedEnv, err := Load(filename)
	if err != nil {
		t.Fatalf("Failed to load written file: %v", err)
	}

	// Verify values that should be preserved
	if loadedEnv["SIMPLE_KEY"] != "simple_value" {
		t.Errorf("SIMPLE_KEY mismatch: expected 'simple_value', got '%s'", loadedEnv["SIMPLE_KEY"])
	}

	if loadedEnv["KEY_WITH_SPACES"] != "value with spaces" {
		t.Errorf("KEY_WITH_SPACES mismatch: expected 'value with spaces', got '%s'", loadedEnv["KEY_WITH_SPACES"])
	}

	if loadedEnv["KEY_WITH_QUOTES"] != "value with \"quotes\"" {
		t.Errorf("KEY_WITH_QUOTES mismatch: expected 'value with \"quotes\"', got '%s'", loadedEnv["KEY_WITH_QUOTES"])
	}

	if loadedEnv["MULTILINE_KEY"] != "line1\nline2\ttab" {
		t.Errorf("MULTILINE_KEY mismatch: expected 'line1\\nline2\\ttab', got '%s'", loadedEnv["MULTILINE_KEY"])
	}
}

func TestMarshalErrors(t *testing.T) {
	// Test nil pointer
	var nilPtr *Config
	_, err := Marshal(nilPtr)
	if err == nil {
		t.Error("Expected error for nil pointer")
	}

	// Test non-struct
	_, err = Marshal("not a struct")
	if err == nil {
		t.Error("Expected error for non-struct")
	}

	// Test struct with unsupported field type
	type BadStruct struct {
		UnsupportedField map[string]string `env:"UNSUPPORTED"`
	}

	bad := BadStruct{
		UnsupportedField: map[string]string{"key": "value"},
	}

	_, err = Marshal(&bad)
	if err == nil {
		t.Error("Expected error for unsupported field type")
	}
}
