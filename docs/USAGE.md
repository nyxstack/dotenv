# Usage Guide

Comprehensive guide for using the dotenv package in your Go applications.

## Table of Contents

1. [Getting Started](#getting-started)
2. [Basic Usage](#basic-usage)
3. [Struct-Based Configuration](#struct-based-configuration)
4. [Typed Environment Access](#typed-environment-access)
5. [Configuration Management Workflows](#configuration-management-workflows)
6. [Multi-Environment Setup](#multi-environment-setup)
7. [Best Practices](#best-practices)

## Getting Started

### Installation

```bash
go get github.com/nyxstack/dotenv
```

### Import

```go
import "github.com/nyxstack/dotenv"
```

### Quick Start

Create a `.env` file:
```bash
DATABASE_URL=postgresql://localhost:5432/myapp
PORT=8080
DEBUG=true
```

Load it in your application:
```go
package main

import (
    "fmt"
    "log"
    "github.com/nyxstack/dotenv"
)

func main() {
    // Load and apply .env file
    err := dotenv.LoadAndApply(".env")
    if err != nil {
        log.Fatal(err)
    }
    
    // Access using standard os package
    fmt.Println("Port:", os.Getenv("PORT"))
}
```

## Basic Usage

### Loading Environment Files

#### Load into a Map

```go
// Load .env file into a map
env, err := dotenv.Load(".env")
if err != nil {
    log.Fatal(err)
}

// Access values
dbURL := env["DATABASE_URL"]
port := env["PORT"]
```

#### Load from Reader

```go
import "strings"

content := "KEY=value\nNAME=myapp"
env, err := dotenv.LoadFromReader(strings.NewReader(content))
```

#### Apply to Process

```go
// Load into map
env, err := dotenv.Load(".env")
if err != nil {
    log.Fatal(err)
}

// Apply to current process
err = dotenv.Apply(env)
if err != nil {
    log.Fatal(err)
}

// Now available via os.Getenv
value := os.Getenv("KEY")
```

#### Combined Load and Apply

```go
// One-step convenience function
err := dotenv.LoadAndApply(".env")
if err != nil {
    log.Fatal(err)
}
```

#### Must Load (Panic on Error)

```go
// For initialization code where failure should halt execution
env := dotenv.MustLoad(".env")
```

## Struct-Based Configuration

The recommended approach for applications with multiple configuration options.

### Basic Struct Setup

```go
package main

import (
    "fmt"
    "log"
    "time"
    "github.com/nyxstack/dotenv"
)

type Config struct {
    // Required field - must be set
    DatabaseURL string `env:"DATABASE_URL,required"`
    
    // Optional field with default
    Port int `env:"PORT,default=8080"`
    
    // Boolean flag
    Debug bool `env:"DEBUG,default=false"`
    
    // Duration type
    Timeout time.Duration `env:"TIMEOUT,default=30s"`
    
    // String slice (comma-separated)
    Features []string `env:"FEATURES"`
    
    // Optional field (no default)
    APIKey string `env:"API_KEY"`
}

func main() {
    // Load .env file first
    err := dotenv.LoadAndApply(".env")
    if err != nil {
        log.Fatal(err)
    }
    
    // Unmarshal into struct
    var config Config
    err = dotenv.Unmarshal(&config)
    if err != nil {
        log.Fatal(err)
    }
    
    fmt.Printf("Server running on port %d\n", config.Port)
    fmt.Printf("Debug mode: %t\n", config.Debug)
    fmt.Printf("Features: %v\n", config.Features)
}
```

### Supported Field Types

```go
type AllTypes struct {
    // String types
    Name    string   `env:"NAME"`
    Items   []string `env:"ITEMS"` // comma-separated
    
    // Integer types
    Count   int      `env:"COUNT"`
    Age     int8     `env:"AGE"`
    Year    int16    `env:"YEAR"`
    ID      int32    `env:"ID"`
    BigNum  int64    `env:"BIG_NUM"`
    
    // Unsigned integer types
    UCount  uint     `env:"UCOUNT"`
    Byte    uint8    `env:"BYTE"`
    Small   uint16   `env:"SMALL"`
    Medium  uint32   `env:"MEDIUM"`
    Large   uint64   `env:"LARGE"`
    
    // Float types
    Price   float32  `env:"PRICE"`
    Rate    float64  `env:"RATE"`
    
    // Boolean
    Enabled bool     `env:"ENABLED"`
    
    // Duration
    Timeout time.Duration `env:"TIMEOUT"`
}
```

### Tag Options

```go
type Config struct {
    // Required: Must be set, error if missing
    APIKey string `env:"API_KEY,required"`
    
    // Default: Use this value if not set
    Port int `env:"PORT,default=8080"`
    
    // Both: Required but with fallback
    Host string `env:"HOST,required,default=localhost"`
    
    // Multiple options (comma-separated in .env)
    Tags []string `env:"TAGS"`
}
```

### Marshaling Structs to Files

```go
// Create configuration
config := Config{
    DatabaseURL: "postgresql://localhost:5432/myapp",
    Port:        8080,
    Debug:       true,
    Timeout:     30 * time.Second,
    Features:    []string{"auth", "logging", "metrics"},
}

// Save to .env file
err := dotenv.MarshalToFile("app.env", &config)
if err != nil {
    log.Fatal(err)
}
```

Output `app.env`:
```bash
DATABASE_URL=postgresql://localhost:5432/myapp
DEBUG=true
FEATURES=auth,logging,metrics
PORT=8080
TIMEOUT=30s
```

### Marshal to Map

```go
config := Config{Port: 8080, Debug: true}
env, err := dotenv.Marshal(&config)
// Returns: map[string]string{"PORT": "8080", "DEBUG": "true"}
```

## Typed Environment Access

Type-safe functions for direct environment variable access with defaults.

### String Values

```go
host := dotenv.Env("DATABASE_HOST", "localhost")
name := dotenv.Env("APP_NAME") // empty string if not set
```

### Integer Values

```go
port := dotenv.EnvInt("PORT", 8080)
maxRetries := dotenv.EnvInt32("MAX_RETRIES", 3)
timeout := dotenv.EnvInt64("TIMEOUT_MS", 30000)
```

### Unsigned Integer Values

```go
maxConnections := dotenv.EnvUint32("MAX_CONNECTIONS", 100)
bufferSize := dotenv.EnvUint64("BUFFER_SIZE", 4096)
```

### Float Values

```go
rate := dotenv.EnvFloat64("RATE", 0.05)
precision := dotenv.EnvFloat32("PRECISION", 0.001)
```

### Boolean Values

```go
debug := dotenv.EnvBool("DEBUG", false)
enableCache := dotenv.EnvBool("ENABLE_CACHE", true)
```

Accepts: `true`, `false`, `1`, `0`, `yes`, `no`, `on`, `off` (case-insensitive)

### Duration Values

```go
timeout := dotenv.EnvDuration("TIMEOUT", 10*time.Second)
interval := dotenv.EnvDuration("CHECK_INTERVAL", 1*time.Minute)
```

Accepts: `"1h"`, `"30m"`, `"45s"`, `"100ms"`, `"1h30m"`, etc.

### Complete Example

```go
package main

import (
    "fmt"
    "time"
    "github.com/nyxstack/dotenv"
)

func main() {
    dotenv.LoadAndApply(".env")
    
    // Database configuration
    dbHost := dotenv.Env("DB_HOST", "localhost")
    dbPort := dotenv.EnvInt("DB_PORT", 5432)
    dbSSL := dotenv.EnvBool("DB_SSL", false)
    dbTimeout := dotenv.EnvDuration("DB_TIMEOUT", 30*time.Second)
    
    // Server configuration
    serverPort := dotenv.EnvInt("SERVER_PORT", 8080)
    debugMode := dotenv.EnvBool("DEBUG", false)
    
    fmt.Printf("Connecting to %s:%d (SSL: %t, Timeout: %v)\n",
        dbHost, dbPort, dbSSL, dbTimeout)
    fmt.Printf("Server starting on port %d (Debug: %t)\n",
        serverPort, debugMode)
}
```

## Configuration Management Workflows

### Workflow 1: Application Initialization

```go
package main

import (
    "log"
    "github.com/nyxstack/dotenv"
)

type AppConfig struct {
    DatabaseURL string `env:"DATABASE_URL,required"`
    Port        int    `env:"PORT,default=8080"`
    LogLevel    string `env:"LOG_LEVEL,default=info"`
}

func main() {
    // 1. Load .env file
    if err := dotenv.LoadAndApply(".env"); err != nil {
        log.Fatal(err)
    }
    
    // 2. Unmarshal into struct
    var config AppConfig
    if err := dotenv.Unmarshal(&config); err != nil {
        log.Fatal(err)
    }
    
    // 3. Use configuration
    startServer(config)
}
```

### Workflow 2: Configuration Generation

```go
// Generate default configuration file
func generateDefaultConfig() error {
    config := AppConfig{
        DatabaseURL: "postgresql://localhost:5432/myapp",
        Port:        8080,
        LogLevel:    "info",
    }
    
    return dotenv.MarshalToFile(".env.example", &config)
}
```

### Workflow 3: Configuration Migration

```go
// Load existing config, modify, and save
func migrateConfig() error {
    // Load existing
    dotenv.LoadAndApply("old.env")
    
    var config AppConfig
    if err := dotenv.Unmarshal(&config); err != nil {
        return err
    }
    
    // Modify
    config.Port = 9000
    config.LogLevel = "debug"
    
    // Save to new location
    return dotenv.MarshalToFile("new.env", &config)
}
```

### Workflow 4: Dynamic Configuration

```go
func updateConfig(key, value string) error {
    // Load current config
    env, err := dotenv.Load(".env")
    if err != nil {
        return err
    }
    
    // Update value
    env[key] = value
    
    // Write back
    return dotenv.WriteEnvFile(".env", env)
}
```

## Multi-Environment Setup

### Environment-Specific Files

```go
package main

import (
    "log"
    "os"
    "github.com/nyxstack/dotenv"
)

type Config struct {
    DatabaseURL string `env:"DATABASE_URL,required"`
    Port        int    `env:"PORT,default=8080"`
    Debug       bool   `env:"DEBUG,default=false"`
}

func loadConfig() (*Config, error) {
    // Determine environment
    env := os.Getenv("APP_ENV")
    if env == "" {
        env = "development"
    }
    
    // Load environment-specific file
    filename := ".env." + env
    if err := dotenv.LoadAndApply(filename); err != nil {
        return nil, err
    }
    
    // Unmarshal
    var config Config
    if err := dotenv.Unmarshal(&config); err != nil {
        return nil, err
    }
    
    return &config, nil
}

func main() {
    config, err := loadConfig()
    if err != nil {
        log.Fatal(err)
    }
    
    // Use config
    startApp(config)
}
```

### Using Prefixes

```go
type DatabaseConfig struct {
    Host     string `env:"HOST"`
    Port     int    `env:"PORT"`
    Database string `env:"NAME"`
}

// Create separate configs with prefixes
func setupMultiDB() error {
    // Primary database
    primaryDB := DatabaseConfig{
        Host:     "primary-db.example.com",
        Port:     5432,
        Database: "primary",
    }
    dotenv.MarshalToFileWithPrefix("primary.env", &primaryDB, "PRIMARY_")
    
    // Replica database
    replicaDB := DatabaseConfig{
        Host:     "replica-db.example.com",
        Port:     5432,
        Database: "replica",
    }
    dotenv.MarshalToFileWithPrefix("replica.env", &replicaDB, "REPLICA_")
    
    return nil
}

// Load with prefixes
func loadMultiDB() error {
    dotenv.LoadAndApply("databases.env")
    
    var primary, replica DatabaseConfig
    
    // Load with different prefixes
    if err := dotenv.UnmarshalWithPrefix(&primary, "PRIMARY_"); err != nil {
        return err
    }
    if err := dotenv.UnmarshalWithPrefix(&replica, "REPLICA_"); err != nil {
        return err
    }
    
    // Use configurations
    connectPrimary(primary)
    connectReplica(replica)
    
    return nil
}
```

### Layered Configuration

```go
// Load multiple files with precedence
func loadLayeredConfig() error {
    // 1. Load defaults
    dotenv.LoadAndApply(".env.defaults")
    
    // 2. Load environment-specific (overrides defaults)
    env := os.Getenv("APP_ENV")
    dotenv.LoadAndApply(".env." + env)
    
    // 3. Load local overrides (overrides everything)
    dotenv.LoadAndApply(".env.local")
    
    return nil
}
```

## Best Practices

### 1. Security

```go
// DON'T commit .env files with secrets
// Add to .gitignore:
// .env
// .env.local
// .env.*.local

// DO provide .env.example
config := Config{
    DatabaseURL: "postgresql://localhost:5432/myapp",
    APIKey:      "", // User must provide
}
dotenv.MarshalToFile(".env.example", &config)
```

### 2. Validation

```go
type Config struct {
    APIKey string `env:"API_KEY,required"`
    Port   int    `env:"PORT,default=8080"`
}

func loadAndValidate() (*Config, error) {
    dotenv.LoadAndApply(".env")
    
    var config Config
    if err := dotenv.Unmarshal(&config); err != nil {
        return nil, err
    }
    
    // Additional validation
    if config.Port < 1024 || config.Port > 65535 {
        return nil, fmt.Errorf("invalid port: %d", config.Port)
    }
    
    return &config, nil
}
```

### 3. Error Handling

```go
func initConfig() (*Config, error) {
    // Attempt to load .env, but don't fail if not found
    if err := dotenv.LoadAndApply(".env"); err != nil {
        log.Printf("Warning: .env file not found: %v", err)
        // Continue with environment variables only
    }
    
    var config Config
    if err := dotenv.Unmarshal(&config); err != nil {
        return nil, fmt.Errorf("config validation failed: %w", err)
    }
    
    return &config, nil
}
```

### 4. Type Safety

```go
// PREFER: Struct-based configuration
type Config struct {
    Port int `env:"PORT,default=8080"`
}
var config Config
dotenv.Unmarshal(&config)

// OVER: Manual type conversion
portStr := os.Getenv("PORT")
port, _ := strconv.Atoi(portStr)

// OR USE: Typed helpers
port := dotenv.EnvInt("PORT", 8080)
```

### 5. Defaults

```go
type Config struct {
    // Provide sensible defaults
    Host         string        `env:"HOST,default=localhost"`
    Port         int           `env:"PORT,default=8080"`
    Timeout      time.Duration `env:"TIMEOUT,default=30s"`
    MaxRetries   int           `env:"MAX_RETRIES,default=3"`
    
    // No default for secrets
    APIKey       string        `env:"API_KEY,required"`
    DatabaseURL  string        `env:"DATABASE_URL,required"`
}
```

### 6. Documentation

```go
// Document your configuration
type Config struct {
    // Server configuration
    Host string `env:"HOST,default=localhost"` // Server hostname
    Port int    `env:"PORT,default=8080"`      // Server port
    
    // Database configuration
    DatabaseURL string `env:"DATABASE_URL,required"` // PostgreSQL connection string
    
    // Feature flags
    EnableCache bool `env:"ENABLE_CACHE,default=true"` // Enable response caching
}
```

### 7. Testing

```go
func TestConfigLoading(t *testing.T) {
    // Set test environment variables
    os.Setenv("PORT", "9000")
    os.Setenv("DEBUG", "true")
    defer os.Unsetenv("PORT")
    defer os.Unsetenv("DEBUG")
    
    var config Config
    err := dotenv.Unmarshal(&config)
    if err != nil {
        t.Fatal(err)
    }
    
    if config.Port != 9000 {
        t.Errorf("Expected port 9000, got %d", config.Port)
    }
}
```

### 8. Initialization Order

```go
func main() {
    // 1. Load environment files first
    if err := dotenv.LoadAndApply(".env"); err != nil {
        log.Printf("No .env file: %v", err)
    }
    
    // 2. Unmarshal into config struct
    var config Config
    if err := dotenv.Unmarshal(&config); err != nil {
        log.Fatal(err)
    }
    
    // 3. Initialize logging
    setupLogging(config.LogLevel)
    
    // 4. Initialize dependencies
    db := initDatabase(config.DatabaseURL)
    
    // 5. Start application
    startServer(config, db)
}
```

### 9. Environment Detection

```go
func getEnvironment() string {
    // Check various sources
    if env := dotenv.Env("APP_ENV"); env != "" {
        return env
    }
    if env := dotenv.Env("ENVIRONMENT"); env != "" {
        return env
    }
    if env := dotenv.Env("GO_ENV"); env != "" {
        return env
    }
    return "development"
}
```

### 10. Configuration Reloading

```go
import "os/signal"
import "syscall"

func watchConfig(configFile string) {
    sigs := make(chan os.Signal, 1)
    signal.Notify(sigs, syscall.SIGHUP)
    
    for {
        <-sigs
        log.Println("Reloading configuration...")
        
        if err := dotenv.LoadAndApply(configFile); err != nil {
            log.Printf("Failed to reload config: %v", err)
        }
    }
}
```
