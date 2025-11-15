package dotenv

import (
	"fmt"
	"io"
	"os"
)

// Load loads environment variables from a .env file
func Load(filename string) (map[string]string, error) {
	data, err := os.ReadFile(filename)
	if err != nil {
		return nil, fmt.Errorf("failed to read file %s: %w", filename, err)
	}

	parser := NewParser(string(data))
	return parser.Parse()
}

// LoadFromReader loads environment variables from an io.Reader
func LoadFromReader(reader io.Reader) (map[string]string, error) {
	data, err := io.ReadAll(reader)
	if err != nil {
		return nil, fmt.Errorf("failed to read data: %w", err)
	}

	parser := NewParser(string(data))
	return parser.Parse()
}

// MustLoad loads environment variables and panics on error
func MustLoad(filename string) map[string]string {
	env, err := Load(filename)
	if err != nil {
		panic(err)
	}
	return env
}

// Apply applies the environment variables to the current process
func Apply(env map[string]string) error {
	for key, value := range env {
		if err := os.Setenv(key, value); err != nil {
			return fmt.Errorf("failed to set environment variable %s: %w", key, err)
		}
	}
	return nil
}

// LoadAndApply loads and applies environment variables from a file
func LoadAndApply(filename string) error {
	env, err := Load(filename)
	if err != nil {
		return err
	}
	return Apply(env)
}
