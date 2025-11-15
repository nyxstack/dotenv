package dotenv

import (
	"fmt"
	"os"
	"reflect"
	"strconv"
	"strings"
	"time"
)

// Marshal converts a struct with `env` tags to environment variable format
func Marshal(v interface{}) (map[string]string, error) {
	return MarshalWithPrefix(v, "")
}

// MarshalWithPrefix converts a struct to environment variables with a prefix
func MarshalWithPrefix(v interface{}, prefix string) (map[string]string, error) {
	rv := reflect.ValueOf(v)

	// Handle pointer to struct
	if rv.Kind() == reflect.Ptr {
		if rv.IsNil() {
			return nil, fmt.Errorf("marshal source cannot be nil pointer")
		}
		rv = rv.Elem()
	}

	if rv.Kind() != reflect.Struct {
		return nil, fmt.Errorf("marshal source must be a struct or pointer to struct")
	}

	rt := rv.Type()
	env := make(map[string]string)

	for i := 0; i < rv.NumField(); i++ {
		field := rv.Field(i)
		fieldType := rt.Field(i)

		// Skip unexported fields
		if !field.CanInterface() {
			continue
		}

		// Get env tag
		envTag := fieldType.Tag.Get("env")
		if envTag == "" {
			continue
		}

		// Parse tag to get key name (ignore other options)
		parts := strings.Split(envTag, ",")
		envKey := parts[0]

		// Add prefix if specified
		if prefix != "" {
			envKey = prefix + envKey
		}

		// Convert field value to string
		value, err := fieldToString(field)
		if err != nil {
			return nil, fmt.Errorf("failed to marshal field %s: %w", fieldType.Name, err)
		}

		// Only add non-empty values (skip zero values)
		if value != "" {
			env[envKey] = value
		}
	}

	return env, nil
}

// fieldToString converts a reflect.Value to its string representation
func fieldToString(field reflect.Value) (string, error) {
	switch field.Kind() {
	case reflect.String:
		return field.String(), nil

	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		// Handle time.Duration specially
		if field.Type() == reflect.TypeOf(time.Duration(0)) {
			duration := time.Duration(field.Int())
			return duration.String(), nil
		}
		return strconv.FormatInt(field.Int(), 10), nil

	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		return strconv.FormatUint(field.Uint(), 10), nil

	case reflect.Float32, reflect.Float64:
		return strconv.FormatFloat(field.Float(), 'g', -1, 64), nil

	case reflect.Bool:
		return strconv.FormatBool(field.Bool()), nil

	case reflect.Slice:
		// Handle slices by joining with comma
		if field.Type().Elem().Kind() == reflect.String {
			var parts []string
			for i := 0; i < field.Len(); i++ {
				parts = append(parts, field.Index(i).String())
			}
			return strings.Join(parts, ","), nil
		}
		return "", fmt.Errorf("unsupported slice type: %s", field.Type())

	default:
		return "", fmt.Errorf("unsupported field type: %s", field.Type())
	}
}

// MarshalToFile writes a struct to a .env file
func MarshalToFile(filename string, v interface{}) error {
	return MarshalToFileWithPrefix(filename, v, "")
}

// MarshalToFileWithPrefix writes a struct to a .env file with prefix
func MarshalToFileWithPrefix(filename string, v interface{}, prefix string) error {
	env, err := MarshalWithPrefix(v, prefix)
	if err != nil {
		return err
	}

	return WriteEnvFile(filename, env)
}

// WriteEnvFile writes a map of environment variables to a .env file
func WriteEnvFile(filename string, env map[string]string) error {
	var lines []string

	// Sort keys for consistent output
	var keys []string
	for key := range env {
		keys = append(keys, key)
	}

	// Simple sort (could use sort.Strings but avoiding extra import)
	for i := 0; i < len(keys); i++ {
		for j := i + 1; j < len(keys); j++ {
			if keys[i] > keys[j] {
				keys[i], keys[j] = keys[j], keys[i]
			}
		}
	}

	for _, key := range keys {
		value := env[key]

		// Quote values that contain spaces or special characters
		if needsQuoting(value) {
			value = quoteValue(value)
		}

		lines = append(lines, fmt.Sprintf("%s=%s", key, value))
	}

	content := strings.Join(lines, "\n")
	if content != "" {
		content += "\n" // Add final newline
	}

	return os.WriteFile(filename, []byte(content), 0644)
}

// needsQuoting determines if a value needs to be quoted
func needsQuoting(value string) bool {
	if value == "" {
		return false
	}

	// Quote if contains spaces, quotes, or special characters
	for _, ch := range value {
		switch ch {
		case ' ', '\t', '\n', '\r', '"', '\'', '\\', '#', '$':
			return true
		}
	}

	return false
}

// quoteValue properly quotes and escapes a value
func quoteValue(value string) string {
	// Use double quotes and escape necessary characters
	escaped := strings.ReplaceAll(value, "\\", "\\\\")  // Escape backslashes first
	escaped = strings.ReplaceAll(escaped, "\"", "\\\"") // Escape quotes
	escaped = strings.ReplaceAll(escaped, "\n", "\\n")  // Escape newlines
	escaped = strings.ReplaceAll(escaped, "\t", "\\t")  // Escape tabs
	escaped = strings.ReplaceAll(escaped, "\r", "\\r")  // Escape carriage returns

	return fmt.Sprintf("\"%s\"", escaped)
}

// Unmarshal populates a struct with environment variables based on `env` tags
func Unmarshal(v interface{}) error {
	return UnmarshalWithPrefix(v, "")
}

// UnmarshalWithPrefix populates a struct with environment variables using a prefix
func UnmarshalWithPrefix(v interface{}, prefix string) error {
	rv := reflect.ValueOf(v)
	if rv.Kind() != reflect.Ptr || rv.Elem().Kind() != reflect.Struct {
		return fmt.Errorf("unmarshal target must be a pointer to struct")
	}

	rv = rv.Elem()
	rt := rv.Type()

	for i := 0; i < rv.NumField(); i++ {
		field := rv.Field(i)
		fieldType := rt.Field(i)

		// Skip unexported fields
		if !field.CanSet() {
			continue
		}

		// Get env tag
		envTag := fieldType.Tag.Get("env")
		if envTag == "" {
			continue
		}

		// Parse tag options (e.g., "KEY,required,default=value")
		parts := strings.Split(envTag, ",")
		envKey := parts[0]

		// Add prefix if specified
		if prefix != "" {
			envKey = prefix + envKey
		}

		// Parse options
		var defaultValue string
		var required bool

		for _, part := range parts[1:] {
			part = strings.TrimSpace(part)
			if part == "required" {
				required = true
			} else if strings.HasPrefix(part, "default=") {
				defaultValue = part[8:] // len("default=") = 8
			}
		}

		// Get environment variable
		envValue, exists := os.LookupEnv(envKey)
		if !exists {
			if required {
				return fmt.Errorf("required environment variable %s is not set", envKey)
			}
			if defaultValue != "" {
				envValue = defaultValue
			} else {
				continue // Skip if no value and not required
			}
		}

		// Set field value with type conversion
		if err := setFieldValue(field, envValue, envKey); err != nil {
			return err
		}
	}

	return nil
}

// setFieldValue converts and sets a field value from a string
func setFieldValue(field reflect.Value, value string, envKey string) error {
	switch field.Kind() {
	case reflect.String:
		field.SetString(value)

	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		// Handle time.Duration specially
		if field.Type() == reflect.TypeOf(time.Duration(0)) {
			duration, err := time.ParseDuration(value)
			if err != nil {
				return fmt.Errorf("failed to parse duration for %s: %w", envKey, err)
			}
			field.SetInt(int64(duration))
		} else {
			intVal, err := strconv.ParseInt(value, 10, 64)
			if err != nil {
				return fmt.Errorf("failed to parse int for %s: %w", envKey, err)
			}
			field.SetInt(intVal)
		}

	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		uintVal, err := strconv.ParseUint(value, 10, 64)
		if err != nil {
			return fmt.Errorf("failed to parse uint for %s: %w", envKey, err)
		}
		field.SetUint(uintVal)

	case reflect.Float32, reflect.Float64:
		floatVal, err := strconv.ParseFloat(value, 64)
		if err != nil {
			return fmt.Errorf("failed to parse float for %s: %w", envKey, err)
		}
		field.SetFloat(floatVal)

	case reflect.Bool:
		boolVal, err := strconv.ParseBool(value)
		if err != nil {
			return fmt.Errorf("failed to parse bool for %s: %w", envKey, err)
		}
		field.SetBool(boolVal)

	case reflect.Slice:
		// Handle slices by splitting on comma
		if field.Type().Elem().Kind() == reflect.String {
			parts := strings.Split(value, ",")
			slice := reflect.MakeSlice(field.Type(), len(parts), len(parts))
			for i, part := range parts {
				slice.Index(i).SetString(strings.TrimSpace(part))
			}
			field.Set(slice)
		} else {
			return fmt.Errorf("unsupported slice type for %s", envKey)
		}

	default:
		return fmt.Errorf("unsupported field type %s for %s", field.Type(), envKey)
	}

	return nil
}
