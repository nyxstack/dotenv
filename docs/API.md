# API Reference

Complete API documentation for the dotenv package.

## Core Functions

### Loading Functions

#### `Load(filename string) (map[string]string, error)`
Load environment variables from a file.

```go
env, err := dotenv.Load(".env")
if err != nil {
    log.Fatal(err)
}
```

**Parameters:**
- `filename`: Path to the .env file

**Returns:**
- `map[string]string`: Environment variables as key-value pairs
- `error`: Error if file cannot be read or parsed

---

#### `LoadFromReader(reader io.Reader) (map[string]string, error)`
Load environment variables from any io.Reader.

```go
env, err := dotenv.LoadFromReader(strings.NewReader("KEY=value"))
```

**Parameters:**
- `reader`: Any io.Reader containing .env content

**Returns:**
- `map[string]string`: Environment variables
- `error`: Parse or read error

---

#### `MustLoad(filename string) map[string]string`
Load environment variables and panic on error. Use for initialization where failure should halt execution.

```go
env := dotenv.MustLoad(".env")
```

**Parameters:**
- `filename`: Path to the .env file

**Returns:**
- `map[string]string`: Environment variables

**Panics:** If file cannot be loaded or parsed

---

#### `Apply(env map[string]string) error`
Apply environment variables to the current process.

```go
env := map[string]string{"KEY": "value"}
err := dotenv.Apply(env)
```

**Parameters:**
- `env`: Map of environment variables to set

**Returns:**
- `error`: Error if any variable cannot be set

---

#### `LoadAndApply(filename string) error`
Convenience function that loads and applies environment variables in one call.

```go
err := dotenv.LoadAndApply(".env")
```

**Parameters:**
- `filename`: Path to the .env file

**Returns:**
- `error`: Load or apply error

---

## Struct-Based Configuration

### `Unmarshal(v interface{}) error`
Populate a struct with environment variables using `env` tags.

```go
type Config struct {
    DatabaseURL string        `env:"DATABASE_URL,required"`
    Port        int           `env:"PORT,default=8080"`
    Debug       bool          `env:"DEBUG,default=false"`
    Timeout     time.Duration `env:"TIMEOUT,default=30s"`
    Features    []string      `env:"FEATURES"`
}

var config Config
err := dotenv.Unmarshal(&config)
```

**Parameters:**
- `v`: Pointer to struct with `env` tags

**Returns:**
- `error`: Type conversion or validation error

**Tag Format:** `env:"VARIABLE_NAME[,option1][,option2]"`

**Tag Options:**
- `required` - Field must have a value in environment
- `default=value` - Default value if environment variable not set

**Supported Types:**
- `string`, `[]string` (comma-separated values)
- `int`, `int8`, `int16`, `int32`, `int64`
- `uint`, `uint8`, `uint16`, `uint32`, `uint64`
- `float32`, `float64`
- `bool` (accepts: true/false, 1/0, yes/no, on/off)
- `time.Duration` (e.g., "1h", "30m", "45s")

---

### `UnmarshalWithPrefix(v interface{}, prefix string) error`
Same as Unmarshal but adds a prefix to all environment variable names.

```go
type Config struct {
    Host string `env:"HOST"`
    Port int    `env:"PORT"`
}

var config Config
// Looks for MYAPP_HOST, MYAPP_PORT
err := dotenv.UnmarshalWithPrefix(&config, "MYAPP_")
```

**Parameters:**
- `v`: Pointer to struct with `env` tags
- `prefix`: String prefix to add to all variable names

**Returns:**
- `error`: Type conversion or validation error

---

### `Marshal(v interface{}) (map[string]string, error)`
Convert a struct with `env` tags to a map of environment variables.

```go
type Config struct {
    DatabaseURL string `env:"DATABASE_URL"`
    Port        int    `env:"PORT"`
    Debug       bool   `env:"DEBUG"`
}

config := Config{
    DatabaseURL: "postgresql://localhost:5432/myapp",
    Port:        8080,
    Debug:       true,
}

env, err := dotenv.Marshal(&config)
// Returns: map[string]string{
//   "DATABASE_URL": "postgresql://localhost:5432/myapp",
//   "PORT": "8080", 
//   "DEBUG": "true",
// }
```

**Parameters:**
- `v`: Struct or pointer to struct with `env` tags

**Returns:**
- `map[string]string`: Environment variables
- `error`: Marshaling error

**Note:** Only non-zero values are included in the output map.

---

### `MarshalWithPrefix(v interface{}, prefix string) (map[string]string, error)`
Same as Marshal but adds a prefix to all environment variable names.

```go
// Creates MYAPP_DATABASE_URL, MYAPP_PORT, etc.
env, err := dotenv.MarshalWithPrefix(&config, "MYAPP_")
```

**Parameters:**
- `v`: Struct or pointer to struct with `env` tags
- `prefix`: String prefix to add to all variable names

**Returns:**
- `map[string]string`: Prefixed environment variables
- `error`: Marshaling error

---

### `MarshalToFile(filename string, v interface{}) error`
Marshal a struct directly to a .env file with proper formatting and quoting.

```go
config := Config{Port: 8080, Debug: true}
err := dotenv.MarshalToFile("config.env", &config)
```

**Parameters:**
- `filename`: Path where .env file will be written
- `v`: Struct or pointer to struct with `env` tags

**Returns:**
- `error`: Marshaling or file write error

**Output Format:**
- Keys are sorted alphabetically
- Values with special characters are automatically quoted
- Proper escape sequences applied

---

### `MarshalToFileWithPrefix(filename string, v interface{}, prefix string) error`
Marshal a struct to a .env file with prefixed keys.

```go
err := dotenv.MarshalToFileWithPrefix("db.env", &dbConfig, "DB_")
```

**Parameters:**
- `filename`: Path where .env file will be written
- `v`: Struct or pointer to struct with `env` tags
- `prefix`: String prefix to add to all variable names

**Returns:**
- `error`: Marshaling or file write error

---

### `WriteEnvFile(filename string, env map[string]string) error`
Write a map of environment variables to a .env file with proper quoting.

```go
env := map[string]string{
    "KEY1": "simple_value",
    "KEY2": "value with spaces",
    "KEY3": "value with \"quotes\" and\nnewlines",
}
err := dotenv.WriteEnvFile("output.env", env)
```

**Parameters:**
- `filename`: Path where .env file will be written
- `env`: Map of environment variables

**Returns:**
- `error`: File write error

**Features:**
- Automatic quoting for values with spaces/special characters
- Proper escape sequences (`\n`, `\t`, `\"`, etc.)
- Sorted keys for consistent output

---

## Typed Environment Variable Functions

Type-safe functions for accessing environment variables with automatic conversion and default values.

### String

#### `Env(key string, defaultValue ...string) string`
Get environment variable as string.

```go
host := dotenv.Env("DATABASE_HOST", "localhost")
```

**Parameters:**
- `key`: Environment variable name
- `defaultValue`: Optional default if variable not set

**Returns:** String value or default

---

### Integer Types

#### `EnvInt(key string, defaultValue ...int) int`
```go
port := dotenv.EnvInt("PORT", 8080)
```

#### `EnvInt8(key string, defaultValue ...int8) int8`
#### `EnvInt16(key string, defaultValue ...int16) int16`
#### `EnvInt32(key string, defaultValue ...int32) int32`
#### `EnvInt64(key string, defaultValue ...int64) int64`

```go
timeout := dotenv.EnvInt64("TIMEOUT_MS", 30000)
```

**Parameters:**
- `key`: Environment variable name
- `defaultValue`: Optional default if variable not set or parse fails

**Returns:** Integer value or default (0 if no default provided)

---

### Unsigned Integer Types

#### `EnvUint(key string, defaultValue ...uint) uint`
#### `EnvUint8(key string, defaultValue ...uint8) uint8`
#### `EnvUint16(key string, defaultValue ...uint16) uint16`
#### `EnvUint32(key string, defaultValue ...uint32) uint32`
#### `EnvUint64(key string, defaultValue ...uint64) uint64`

```go
maxConn := dotenv.EnvUint32("MAX_CONNECTIONS", 100)
```

**Parameters:**
- `key`: Environment variable name
- `defaultValue`: Optional default if variable not set or parse fails

**Returns:** Unsigned integer value or default

---

### Floating Point Types

#### `EnvFloat32(key string, defaultValue ...float32) float32`
#### `EnvFloat64(key string, defaultValue ...float64) float64`

```go
pi := dotenv.EnvFloat64("PI", 3.14159)
rate := dotenv.EnvFloat32("RATE", 0.05)
```

**Parameters:**
- `key`: Environment variable name
- `defaultValue`: Optional default if variable not set or parse fails

**Returns:** Float value or default

---

### Boolean Type

#### `EnvBool(key string, defaultValue ...bool) bool`
Get environment variable as boolean.

```go
debug := dotenv.EnvBool("DEBUG", false)
```

**Accepted Values** (case-insensitive):
- True: `true`, `1`, `yes`, `on`
- False: `false`, `0`, `no`, `off`

**Parameters:**
- `key`: Environment variable name
- `defaultValue`: Optional default if variable not set or parse fails

**Returns:** Boolean value or default

---

### Duration Type

#### `EnvDuration(key string, defaultValue ...time.Duration) time.Duration`
Get environment variable as time.Duration.

```go
timeout := dotenv.EnvDuration("TIMEOUT", 10*time.Second)
```

**Accepted Formats:**
- `"1h"` - hours
- `"30m"` - minutes
- `"45s"` - seconds
- `"100ms"` - milliseconds
- `"1h30m"` - combined

**Parameters:**
- `key`: Environment variable name
- `defaultValue`: Optional default if variable not set or parse fails

**Returns:** Duration value or default

---

## Environment Management Functions

### `HasEnv(key string) bool`
Check if an environment variable is set.

```go
if dotenv.HasEnv("API_KEY") {
    // Variable exists
}
```

**Parameters:**
- `key`: Environment variable name

**Returns:** `true` if variable exists, `false` otherwise

---

### `SetEnv(key, value string) error`
Set an environment variable in the current process.

```go
err := dotenv.SetEnv("NEW_VAR", "value")
```

**Parameters:**
- `key`: Environment variable name
- `value`: Value to set

**Returns:** Error if variable cannot be set

---

### `UnsetEnv(key string) error`
Unset an environment variable in the current process.

```go
err := dotenv.UnsetEnv("OLD_VAR")
```

**Parameters:**
- `key`: Environment variable name

**Returns:** Error if variable cannot be unset

---

## Parser API

For advanced use cases requiring more control over the parsing process.

### `NewParser(content string) *Parser`
Create a new parser for the given .env content.

```go
parser := dotenv.NewParser(content)
env, err := parser.Parse()
```

**Parameters:**
- `content`: String containing .env file content

**Returns:** `*Parser` instance

---

### `Parse() (map[string]string, error)`
Parse the entire content and return environment variables.

```go
parser := dotenv.NewParser(content)
env, err := parser.Parse()
```

**Returns:**
- `map[string]string`: Parsed environment variables
- `error`: Parse error with line number

**Features:**
- Single-pass parsing
- Variable expansion
- Detailed error messages

---

### `ParseLine() LineResult`
Parse a single line of .env content.

```go
parser := dotenv.NewParser(content)
for parser.tokenizer.pos < parser.tokenizer.length {
    result := parser.ParseLine()
    if result.Error != nil {
        return result.Error
    }
    if result.Key != "" {
        fmt.Printf("%s=%s\n", result.Key, result.Value)
    }
}
```

**Returns:** `LineResult` struct containing:
- `Key`: Variable name (empty for blank lines/comments)
- `Value`: Variable value
- `AllowExpansion`: Whether variable expansion is allowed
- `Error`: Parse error if any

**Use Cases:**
- Custom parsing logic
- Line-by-line processing
- Error recovery

---

## Types

### `LineResult`
Result of parsing a single line.

```go
type LineResult struct {
    Key            string  // Variable name
    Value          string  // Variable value
    AllowExpansion bool    // Whether expansion is allowed
    Error          error   // Parse error
}
```

---

### `Parser`
Parser instance for .env content.

```go
type Parser struct {
    tokenizer *Tokenizer
}
```

---

## Error Handling

All functions that can fail return errors. Common error types:

1. **File Errors**: File not found, permission denied
2. **Parse Errors**: Syntax errors with line numbers
3. **Type Errors**: Failed type conversion
4. **Validation Errors**: Required fields missing

```go
env, err := dotenv.Load(".env")
if err != nil {
    log.Printf("Failed to load .env: %v", err)
    // Handle error
}
```

Error messages include line numbers for parse errors:
```
unterminated quoted string at line 5
expected '=' after variable name at line 10
```
