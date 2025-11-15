# Dotenv Parser

A comprehensive Go package for parsing `.env` files with full grammar support, implementing the same features as mature dotenv loaders in other languages.

## Installation

```bash
go get github.com/nyxstack/dotenv
```

Import in your Go code:

```go
import "github.com/nyxstack/dotenv"
```

## Features

- ✅ **Full .env grammar support**
- ✅ **Export syntax**: `export KEY=value`
- ✅ **Quoted strings**: Double quotes with escapes, single quotes literal
- ✅ **Variable expansion**: `$VAR` and `${VAR}` syntax
- ✅ **Inline comments**: `KEY=value # comment`
- ✅ **Escape sequences**: `\n`, `\t`, `\"`, `\\`, etc.
- ✅ **Unquoted values with spaces**: `KEY=some value`
- ✅ **Empty values**: `KEY=`
- ✅ **Comprehensive error handling**

## Comparison with Other Libraries

| Feature | This Package | godotenv | joho/godotenv |
|---------|-------------|-----------|---------------|
| Variable expansion | ✅ | ❌ | ❌ |
| Inline comments | ✅ | ❌ | ❌ |
| Export syntax | ✅ | ❌ | ❌ |
| Proper escape sequences | ✅ | ❌ | ❌ |
| Single quote literals | ✅ | ❌ | ❌ |
| Detailed error messages | ✅ | ❌ | ❌ |
| Full .env compatibility | ✅ | ❌ | ❌ |

## Quick Start

### Basic Usage
```go
package main

import (
    "fmt"
    "log"
    "dotenv"
)

func main() {
    // Load from file
    env, err := dotenv.Load(".env")
    if err != nil {
        log.Fatal(err)
    }
    
    // Access values
    fmt.Println("Database URL:", env["DATABASE_URL"])
    
    // Apply to current process
    err = dotenv.Apply(env)
    if err != nil {
        log.Fatal(err)
    }
}
```

### Struct-Based Configuration
```go
package main

import (
    "fmt"
    "log"
    "time"
    "dotenv"
)

type Config struct {
    DatabaseURL     string        `env:"DATABASE_URL,required"`
    Port           int           `env:"PORT,default=8080"`
    Debug          bool          `env:"DEBUG,default=false"`
    RequestTimeout time.Duration `env:"REQUEST_TIMEOUT,default=30s"`
    Features       []string      `env:"FEATURES"`
}

func main() {
    // Load .env file first
    dotenv.LoadAndApply(".env")
    
    // Populate struct from environment
    var config Config
    err := dotenv.Unmarshal(&config)
    if err != nil {
        log.Fatal(err)
    }
    
    fmt.Printf("Server running on port %d\n", config.Port)
    fmt.Printf("Debug mode: %t\n", config.Debug)
    fmt.Printf("Features: %v\n", config.Features)
}
```

### Typed Environment Access
```go
package main

import (
    "fmt"
    "time"
    "dotenv"
)

func main() {
    // Load configuration with type safety and defaults
    host := dotenv.Env("DATABASE_HOST", "localhost")
    port := dotenv.EnvInt("DATABASE_PORT", 5432)
    ssl := dotenv.EnvBool("DATABASE_SSL", false)
    timeout := dotenv.EnvDuration("DATABASE_TIMEOUT", 30*time.Second)
    maxConn := dotenv.EnvInt64("MAX_CONNECTIONS", 100)
    
    fmt.Printf("Connecting to %s:%d (SSL: %t, Timeout: %v, MaxConn: %d)\n", 
        host, port, ssl, timeout, maxConn)
}
```

### Configuration Management Workflow

The package supports a complete configuration management lifecycle:

```go
// 1. Define configuration struct
type AppConfig struct {
    DatabaseURL string        `env:"DATABASE_URL,required"`
    Port        int           `env:"PORT,default=8080"`
    Debug       bool          `env:"DEBUG,default=false"`
    Features    []string      `env:"FEATURES"`
}

// 2. Create config with application defaults
config := AppConfig{
    DatabaseURL: "postgresql://localhost:5432/app",
    Port:        3000,
    Debug:       true,
    Features:    []string{"auth", "logging"},
}

// 3. Export config to .env file
err := dotenv.MarshalToFile("app.env", &config)

// 4. Later, load config from environment
dotenv.LoadAndApply("app.env")
var loadedConfig AppConfig
err = dotenv.Unmarshal(&loadedConfig)

// 5. Modify and save back
loadedConfig.Port = 8080
loadedConfig.Features = append(loadedConfig.Features, "metrics")
err = dotenv.MarshalToFile("app.env", &loadedConfig)
```

### Multi-Environment Configuration

```go
// Development config
devConfig := Config{DatabaseURL: "localhost:5432/app_dev"}
dotenv.MarshalToFileWithPrefix("dev.env", &devConfig, "DEV_")

// Production config  
prodConfig := Config{DatabaseURL: "prod-db:5432/app"}
dotenv.MarshalToFileWithPrefix("prod.env", &prodConfig, "PROD_")

// Load based on environment
env := os.Getenv("ENVIRONMENT")
prefix := strings.ToUpper(env) + "_"
err := dotenv.LoadAndApply(env + ".env")
err = dotenv.UnmarshalWithPrefix(&config, prefix)
```

## Supported .env Grammar

### Basic key=value pairs
```bash
KEY1=value1
KEY2=value2
KEY3=  # empty value
```

### Export syntax
```bash
export NODE_ENV=production
export DEBUG=true
```

### Quoted strings
```bash
# Double quotes allow escapes
MESSAGE="Hello\nWorld\t!"
PATH="C:\\Program Files\\App"

# Single quotes are literal
LITERAL='No escapes\n here'
REGEX='^\d{3}-\d{3}-\d{4}$'
```

### Variable expansion
```bash
HOME_DIR=/home/user
CONFIG_PATH="${HOME_DIR}/config"    # ${VAR} syntax
BACKUP_PATH="$HOME_DIR/backups"     # $VAR syntax

# Expansion only works in double quotes
NO_EXPAND='$HOME/literal'           # stays literal
```

### Inline comments
```bash
MAX_CONNECTIONS=100 # Database connection limit
API_URL="https://api.example.com"  # Production API

# Comments inside quotes are preserved
PASSWORD="secret#123!"  # The # is part of the password
```

### Unquoted values with spaces
```bash
APP_NAME=My Application Name
LOG_FORMAT=%timestamp% %level%: %message%
```

## Documentation

For comprehensive API documentation, usage examples, and best practices, see:

- **[Usage Guide](docs/USAGE.md)** - Complete guide to using the package
- **[API Reference](docs/API.md)** - Detailed API documentation
- **[Examples](docs/EXAMPLES.md)** - Practical examples and patterns
- **[AI Agents Guide](AGENTS.md)** - Quick reference for AI assistants

## License

MIT License - see LICENSE file for details.

## Contributing

Contributions are welcome! Please feel free to submit a Pull Request.