# Dotenv Package - Quick Reference for AI Agents

## Package Overview
`github.com/nyxstack/dotenv` - A comprehensive Go package for parsing and managing `.env` files with full grammar support including variable expansion, quoted strings, and escape sequences.

## Installation
```go
import "github.com/nyxstack/dotenv"
```

## Core Concepts

### 1. Basic Loading & Applying
```go
// Load .env file into a map
env, err := dotenv.Load(".env")

// Apply to current process
err := dotenv.Apply(env)

// Or combine both
err := dotenv.LoadAndApply(".env")
```

### 2. Struct-Based Configuration (Recommended)
```go
type Config struct {
    DatabaseURL string        `env:"DATABASE_URL,required"`
    Port        int           `env:"PORT,default=8080"`
    Debug       bool          `env:"DEBUG,default=false"`
    Timeout     time.Duration `env:"TIMEOUT,default=30s"`
    Features    []string      `env:"FEATURES"`  // comma-separated
}

// Load from environment into struct
var config Config
err := dotenv.Unmarshal(&config)

// Save struct to .env file
err := dotenv.MarshalToFile("config.env", &config)
```

**Tag Options:**
- `required` - Field must have a value
- `default=value` - Default if not set
- Slices split on comma

**Supported Types:**
- Basic: `string`, `bool`
- Integers: `int`, `int8/16/32/64`, `uint`, `uint8/16/32/64`
- Floats: `float32`, `float64`
- Special: `time.Duration`, `[]string`

### 3. Typed Environment Access
```go
// With type safety and defaults
host := dotenv.Env("DB_HOST", "localhost")
port := dotenv.EnvInt("DB_PORT", 5432)
ssl := dotenv.EnvBool("DB_SSL", false)
timeout := dotenv.EnvDuration("TIMEOUT", 30*time.Second)
```

**Available Functions:**
`Env`, `EnvInt`, `EnvInt8/16/32/64`, `EnvUint`, `EnvUint8/16/32/64`, `EnvFloat32/64`, `EnvBool`, `EnvDuration`

### 4. Environment Management
```go
// Check existence
if dotenv.HasEnv("API_KEY") { }

// Set/Unset
dotenv.SetEnv("KEY", "value")
dotenv.UnsetEnv("KEY")
```

## .env File Syntax

### Basic Format
```bash
KEY=value
KEY2=              # empty value
export KEY3=value  # export syntax supported
```

### Quoted Strings
```bash
# Double quotes: allow escapes and variable expansion
MESSAGE="Hello\nWorld"
PATH="$HOME/config"

# Single quotes: literal (no escapes/expansion)
REGEX='^\d{3}-\d{3}$'
LITERAL='$HOME stays literal'
```

### Variable Expansion (in double quotes only)
```bash
HOME=/home/user
CONFIG="${HOME}/config"    # ${VAR} syntax
BACKUP="$HOME/backups"     # $VAR syntax
```

### Comments
```bash
# Full line comment
KEY=value # inline comment
PASSWORD="secret#123"  # # inside quotes is literal
```

### Unquoted Values
```bash
APP_NAME=My Application Name  # spaces preserved
```

## Common Usage Patterns

### Pattern 1: Simple Load and Use
```go
dotenv.LoadAndApply(".env")
dbHost := os.Getenv("DB_HOST")
```

### Pattern 2: Type-Safe Config
```go
dotenv.LoadAndApply(".env")
port := dotenv.EnvInt("PORT", 8080)
debug := dotenv.EnvBool("DEBUG", false)
```

### Pattern 3: Struct-Based (Best for complex config)
```go
type Config struct {
    DatabaseURL string `env:"DATABASE_URL,required"`
    Port        int    `env:"PORT,default=8080"`
}

dotenv.LoadAndApply(".env")
var cfg Config
dotenv.Unmarshal(&cfg)
```

### Pattern 4: Multi-Environment
```go
// Create environment-specific files
prodCfg := Config{DatabaseURL: "prod-db:5432"}
dotenv.MarshalToFileWithPrefix("prod.env", &prodCfg, "PROD_")

// Load with prefix
dotenv.LoadAndApply("prod.env")
dotenv.UnmarshalWithPrefix(&cfg, "PROD_")
```

### Pattern 5: Export Config to File
```go
config := Config{Port: 8080, Debug: true}
dotenv.MarshalToFile("output.env", &config)
```

## Key Features for Agents

1. **Full .env Grammar**: Supports all standard .env syntax including export, quotes, escapes, expansion
2. **Variable Expansion**: `$VAR` and `${VAR}` in double-quoted strings
3. **Type Safety**: Struct tags with validation and defaults
4. **Bidirectional**: Load from files OR save structs to files
5. **Error Handling**: Detailed error messages with line numbers

## Quick Decision Tree

**Need to:**
- Just load variables? → `dotenv.LoadAndApply()`
- Type-safe access? → `dotenv.EnvInt()`, `dotenv.EnvBool()`, etc.
- Complex config? → Define struct with `env` tags + `Unmarshal()`
- Save config? → `dotenv.MarshalToFile()`
- Multi-env? → Use prefix variants: `MarshalToFileWithPrefix()`, `UnmarshalWithPrefix()`

## Error Handling

All functions return errors. Common error types:
- File not found
- Parse errors (with line numbers)
- Type conversion errors
- Required field missing

```go
if err := dotenv.LoadAndApply(".env"); err != nil {
    log.Fatal(err)  // or handle gracefully
}
```

## Best Practices

1. **Use struct-based config** for applications with multiple settings
2. **Use `required` tag** for critical configuration
3. **Provide defaults** for optional settings
4. **Load early** in application startup
5. **Don't commit** `.env` files with secrets (use `.env.example`)

## Advanced: Direct Parser Access

For custom parsing logic:
```go
parser := dotenv.NewParser(content)
env, err := parser.Parse()

// Or line-by-line
result := parser.ParseLine()
if result.Error != nil { }
```

## Documentation

- Full API: See `docs/API.md`
- Usage Guide: See `docs/USAGE.md`
- Examples: See `docs/EXAMPLES.md`
