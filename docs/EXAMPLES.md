# Examples

Practical examples demonstrating common use cases for the dotenv package.

## Table of Contents

1. [Web Server Configuration](#web-server-configuration)
2. [Database Connection](#database-connection)
3. [Microservices Configuration](#microservices-configuration)
4. [Multi-Environment Deployment](#multi-environment-deployment)
5. [Feature Flags](#feature-flags)
6. [Configuration Management Tool](#configuration-management-tool)
7. [Docker Integration](#docker-integration)
8. [CLI Application](#cli-application)

## Web Server Configuration

Complete example of a web server using struct-based configuration.

```go
package main

import (
    "fmt"
    "log"
    "net/http"
    "time"
    "github.com/nyxstack/dotenv"
)

type ServerConfig struct {
    Host            string        `env:"SERVER_HOST,default=0.0.0.0"`
    Port            int           `env:"SERVER_PORT,default=8080"`
    ReadTimeout     time.Duration `env:"SERVER_READ_TIMEOUT,default=30s"`
    WriteTimeout    time.Duration `env:"SERVER_WRITE_TIMEOUT,default=30s"`
    ShutdownTimeout time.Duration `env:"SERVER_SHUTDOWN_TIMEOUT,default=10s"`
    TLSEnabled      bool          `env:"SERVER_TLS_ENABLED,default=false"`
    TLSCert         string        `env:"SERVER_TLS_CERT"`
    TLSKey          string        `env:"SERVER_TLS_KEY"`
}

func main() {
    // Load environment configuration
    if err := dotenv.LoadAndApply(".env"); err != nil {
        log.Printf("Warning: .env file not found: %v", err)
    }

    // Parse configuration
    var config ServerConfig
    if err := dotenv.Unmarshal(&config); err != nil {
        log.Fatalf("Failed to load configuration: %v", err)
    }

    // Create server
    srv := &http.Server{
        Addr:         fmt.Sprintf("%s:%d", config.Host, config.Port),
        ReadTimeout:  config.ReadTimeout,
        WriteTimeout: config.WriteTimeout,
    }

    // Setup routes
    http.HandleFunc("/", handleRoot)
    http.HandleFunc("/health", handleHealth)

    // Start server
    log.Printf("Starting server on %s", srv.Addr)
    
    if config.TLSEnabled {
        if config.TLSCert == "" || config.TLSKey == "" {
            log.Fatal("TLS enabled but certificate/key not provided")
        }
        log.Fatal(srv.ListenAndServeTLS(config.TLSCert, config.TLSKey))
    } else {
        log.Fatal(srv.ListenAndServe())
    }
}

func handleRoot(w http.ResponseWriter, r *http.Request) {
    fmt.Fprintf(w, "Hello, World!")
}

func handleHealth(w http.ResponseWriter, r *http.Request) {
    w.WriteHeader(http.StatusOK)
    fmt.Fprintf(w, "OK")
}
```

`.env` file:
```bash
SERVER_HOST=0.0.0.0
SERVER_PORT=8080
SERVER_READ_TIMEOUT=30s
SERVER_WRITE_TIMEOUT=30s
SERVER_SHUTDOWN_TIMEOUT=10s
SERVER_TLS_ENABLED=false
```

## Database Connection

Example showing database configuration with connection pooling.

```go
package main

import (
    "database/sql"
    "fmt"
    "log"
    "time"
    _ "github.com/lib/pq"
    "github.com/nyxstack/dotenv"
)

type DatabaseConfig struct {
    Host            string        `env:"DB_HOST,default=localhost"`
    Port            int           `env:"DB_PORT,default=5432"`
    User            string        `env:"DB_USER,required"`
    Password        string        `env:"DB_PASSWORD,required"`
    Database        string        `env:"DB_NAME,required"`
    SSLMode         string        `env:"DB_SSLMODE,default=disable"`
    MaxConnections  int           `env:"DB_MAX_CONNECTIONS,default=25"`
    MaxIdleConns    int           `env:"DB_MAX_IDLE_CONNECTIONS,default=5"`
    ConnMaxLifetime time.Duration `env:"DB_CONN_MAX_LIFETIME,default=5m"`
}

func (c *DatabaseConfig) ConnectionString() string {
    return fmt.Sprintf(
        "host=%s port=%d user=%s password=%s dbname=%s sslmode=%s",
        c.Host, c.Port, c.User, c.Password, c.Database, c.SSLMode,
    )
}

func main() {
    // Load configuration
    if err := dotenv.LoadAndApply(".env"); err != nil {
        log.Fatal(err)
    }

    var dbConfig DatabaseConfig
    if err := dotenv.Unmarshal(&dbConfig); err != nil {
        log.Fatal(err)
    }

    // Connect to database
    db, err := sql.Open("postgres", dbConfig.ConnectionString())
    if err != nil {
        log.Fatalf("Failed to open database: %v", err)
    }
    defer db.Close()

    // Configure connection pool
    db.SetMaxOpenConns(dbConfig.MaxConnections)
    db.SetMaxIdleConns(dbConfig.MaxIdleConns)
    db.SetConnMaxLifetime(dbConfig.ConnMaxLifetime)

    // Test connection
    if err := db.Ping(); err != nil {
        log.Fatalf("Failed to ping database: %v", err)
    }

    log.Println("Successfully connected to database")
    
    // Run application
    runApplication(db)
}

func runApplication(db *sql.DB) {
    // Your application logic here
}
```

`.env` file:
```bash
DB_HOST=localhost
DB_PORT=5432
DB_USER=myapp
DB_PASSWORD=secretpassword
DB_NAME=myapp_db
DB_SSLMODE=disable
DB_MAX_CONNECTIONS=25
DB_MAX_IDLE_CONNECTIONS=5
DB_CONN_MAX_LIFETIME=5m
```

## Microservices Configuration

Configuration for a microservice with multiple dependencies.

```go
package main

import (
    "log"
    "time"
    "github.com/nyxstack/dotenv"
)

type AppConfig struct {
    // Service configuration
    ServiceName string `env:"SERVICE_NAME,required"`
    Environment string `env:"ENVIRONMENT,default=development"`
    Version     string `env:"VERSION,default=dev"`
    
    // HTTP Server
    HTTPPort int           `env:"HTTP_PORT,default=8080"`
    HTTPHost string        `env:"HTTP_HOST,default=0.0.0.0"`
    Timeout  time.Duration `env:"HTTP_TIMEOUT,default=30s"`
    
    // Database
    DatabaseURL string `env:"DATABASE_URL,required"`
    DBPoolSize  int    `env:"DB_POOL_SIZE,default=10"`
    
    // Redis Cache
    RedisURL     string        `env:"REDIS_URL,required"`
    RedisTTL     time.Duration `env:"REDIS_TTL,default=1h"`
    RedisEnabled bool          `env:"REDIS_ENABLED,default=true"`
    
    // Message Queue
    RabbitMQURL   string `env:"RABBITMQ_URL,required"`
    QueueName     string `env:"QUEUE_NAME,default=default"`
    
    // Logging
    LogLevel  string `env:"LOG_LEVEL,default=info"`
    LogFormat string `env:"LOG_FORMAT,default=json"`
    
    // Monitoring
    MetricsEnabled bool   `env:"METRICS_ENABLED,default=true"`
    MetricsPort    int    `env:"METRICS_PORT,default=9090"`
    TracingEnabled bool   `env:"TRACING_ENABLED,default=false"`
    JaegerURL      string `env:"JAEGER_URL"`
    
    // Feature Flags
    Features []string `env:"FEATURES"`
}

func main() {
    // Load configuration
    if err := dotenv.LoadAndApply(".env"); err != nil {
        log.Fatal(err)
    }

    var config AppConfig
    if err := dotenv.Unmarshal(&config); err != nil {
        log.Fatal(err)
    }

    // Initialize services
    log.Printf("Starting %s v%s in %s mode", 
        config.ServiceName, config.Version, config.Environment)
    
    // Setup logging
    setupLogging(config.LogLevel, config.LogFormat)
    
    // Connect to dependencies
    db := connectDatabase(config.DatabaseURL, config.DBPoolSize)
    defer db.Close()
    
    var cache *RedisClient
    if config.RedisEnabled {
        cache = connectRedis(config.RedisURL)
        defer cache.Close()
    }
    
    mq := connectRabbitMQ(config.RabbitMQURL)
    defer mq.Close()
    
    // Start metrics server
    if config.MetricsEnabled {
        go startMetricsServer(config.MetricsPort)
    }
    
    // Start tracing
    if config.TracingEnabled {
        initTracing(config.JaegerURL)
    }
    
    // Start HTTP server
    startHTTPServer(config, db, cache, mq)
}

// Stub functions
func setupLogging(level, format string)           {}
func connectDatabase(url string, pool int) *DB    { return nil }
func connectRedis(url string) *RedisClient        { return nil }
func connectRabbitMQ(url string) *MQ              { return nil }
func startMetricsServer(port int)                 {}
func initTracing(url string)                      {}
func startHTTPServer(c AppConfig, db *DB, cache *RedisClient, mq *MQ) {}

type DB struct{}
func (d *DB) Close() {}

type RedisClient struct{}
func (r *RedisClient) Close() {}

type MQ struct{}
func (m *MQ) Close() {}
```

`.env` file:
```bash
SERVICE_NAME=user-service
ENVIRONMENT=production
VERSION=1.0.0

HTTP_PORT=8080
HTTP_HOST=0.0.0.0
HTTP_TIMEOUT=30s

DATABASE_URL=postgresql://user:pass@localhost:5432/userdb
DB_POOL_SIZE=20

REDIS_URL=redis://localhost:6379
REDIS_TTL=1h
REDIS_ENABLED=true

RABBITMQ_URL=amqp://guest:guest@localhost:5672/
QUEUE_NAME=user-events

LOG_LEVEL=info
LOG_FORMAT=json

METRICS_ENABLED=true
METRICS_PORT=9090
TRACING_ENABLED=true
JAEGER_URL=http://localhost:14268/api/traces

FEATURES=auth,notifications,analytics
```

## Multi-Environment Deployment

Managing different environments with the same codebase.

```go
package main

import (
    "fmt"
    "log"
    "os"
    "github.com/nyxstack/dotenv"
)

type Config struct {
    Environment string `env:"ENVIRONMENT,default=development"`
    DatabaseURL string `env:"DATABASE_URL,required"`
    APIBaseURL  string `env:"API_BASE_URL,required"`
    Debug       bool   `env:"DEBUG,default=false"`
}

func loadEnvironmentConfig() (*Config, error) {
    // Determine which environment we're in
    env := os.Getenv("APP_ENV")
    if env == "" {
        env = "development"
    }

    // Try to load environment-specific file first
    envFile := fmt.Sprintf(".env.%s", env)
    if err := dotenv.LoadAndApply(envFile); err != nil {
        log.Printf("No %s file found, trying .env", envFile)
        
        // Fallback to default .env
        if err := dotenv.LoadAndApply(".env"); err != nil {
            return nil, fmt.Errorf("failed to load configuration: %w", err)
        }
    }

    // Load local overrides if they exist
    if err := dotenv.LoadAndApply(".env.local"); err == nil {
        log.Println("Loaded local configuration overrides")
    }

    // Parse into struct
    var config Config
    if err := dotenv.Unmarshal(&config); err != nil {
        return nil, err
    }

    return &config, nil
}

func main() {
    config, err := loadEnvironmentConfig()
    if err != nil {
        log.Fatal(err)
    }

    log.Printf("Running in %s environment", config.Environment)
    log.Printf("API Base URL: %s", config.APIBaseURL)
    log.Printf("Debug Mode: %t", config.Debug)
    
    // Start application with config
}
```

Create multiple environment files:

`.env.development`:
```bash
ENVIRONMENT=development
DATABASE_URL=postgresql://localhost:5432/myapp_dev
API_BASE_URL=http://localhost:8080
DEBUG=true
```

`.env.staging`:
```bash
ENVIRONMENT=staging
DATABASE_URL=postgresql://staging-db:5432/myapp_staging
API_BASE_URL=https://api-staging.example.com
DEBUG=false
```

`.env.production`:
```bash
ENVIRONMENT=production
DATABASE_URL=postgresql://prod-db:5432/myapp_prod
API_BASE_URL=https://api.example.com
DEBUG=false
```

## Feature Flags

Dynamic feature flag system using environment variables.

```go
package main

import (
    "log"
    "strings"
    "github.com/nyxstack/dotenv"
)

type FeatureFlags struct {
    EnabledFeatures []string `env:"ENABLED_FEATURES"`
}

type Features struct {
    Authentication bool
    UserProfiles   bool
    Analytics      bool
    BetaFeatures   bool
    ExperimentalUI bool
}

func loadFeatures() (*Features, error) {
    if err := dotenv.LoadAndApply(".env"); err != nil {
        log.Printf("Warning: %v", err)
    }

    var flags FeatureFlags
    if err := dotenv.Unmarshal(&flags); err != nil {
        return nil, err
    }

    // Convert to feature struct
    features := &Features{
        Authentication: contains(flags.EnabledFeatures, "auth"),
        UserProfiles:   contains(flags.EnabledFeatures, "profiles"),
        Analytics:      contains(flags.EnabledFeatures, "analytics"),
        BetaFeatures:   contains(flags.EnabledFeatures, "beta"),
        ExperimentalUI: contains(flags.EnabledFeatures, "experimental-ui"),
    }

    return features, nil
}

func contains(slice []string, item string) bool {
    for _, s := range slice {
        if strings.TrimSpace(s) == item {
            return true
        }
    }
    return false
}

func main() {
    features, err := loadFeatures()
    if err != nil {
        log.Fatal(err)
    }

    log.Println("Feature flags:")
    log.Printf("  Authentication: %t", features.Authentication)
    log.Printf("  User Profiles: %t", features.UserProfiles)
    log.Printf("  Analytics: %t", features.Analytics)
    log.Printf("  Beta Features: %t", features.BetaFeatures)
    log.Printf("  Experimental UI: %t", features.ExperimentalUI)

    // Use features in application
    if features.Authentication {
        enableAuthentication()
    }
    if features.Analytics {
        enableAnalytics()
    }
}

func enableAuthentication() { log.Println("Authentication enabled") }
func enableAnalytics()       { log.Println("Analytics enabled") }
```

`.env`:
```bash
ENABLED_FEATURES=auth,profiles,analytics
```

## Configuration Management Tool

CLI tool for managing .env files.

```go
package main

import (
    "flag"
    "fmt"
    "log"
    "os"
    "github.com/nyxstack/dotenv"
)

func main() {
    var (
        file      = flag.String("file", ".env", "Environment file to manage")
        get       = flag.String("get", "", "Get value of a key")
        set       = flag.String("set", "", "Set key=value")
        del       = flag.String("del", "", "Delete a key")
        list      = flag.Bool("list", false, "List all variables")
        validate  = flag.Bool("validate", false, "Validate file syntax")
    )
    flag.Parse()

    switch {
    case *get != "":
        getValue(*file, *get)
    case *set != "":
        setValue(*file, *set)
    case *del != "":
        deleteKey(*file, *del)
    case *list:
        listVars(*file)
    case *validate:
        validateFile(*file)
    default:
        flag.Usage()
    }
}

func getValue(file, key string) {
    env, err := dotenv.Load(file)
    if err != nil {
        log.Fatal(err)
    }

    if value, exists := env[key]; exists {
        fmt.Println(value)
    } else {
        log.Fatalf("Key %s not found", key)
    }
}

func setValue(file, keyValue string) {
    // Parse key=value
    parts := splitFirst(keyValue, "=")
    if len(parts) != 2 {
        log.Fatal("Invalid format. Use: key=value")
    }
    key, value := parts[0], parts[1]

    // Load existing
    env, err := dotenv.Load(file)
    if err != nil && !os.IsNotExist(err) {
        log.Fatal(err)
    }
    if env == nil {
        env = make(map[string]string)
    }

    // Set value
    env[key] = value

    // Write back
    if err := dotenv.WriteEnvFile(file, env); err != nil {
        log.Fatal(err)
    }

    fmt.Printf("Set %s=%s\n", key, value)
}

func deleteKey(file, key string) {
    env, err := dotenv.Load(file)
    if err != nil {
        log.Fatal(err)
    }

    if _, exists := env[key]; !exists {
        log.Fatalf("Key %s not found", key)
    }

    delete(env, key)

    if err := dotenv.WriteEnvFile(file, env); err != nil {
        log.Fatal(err)
    }

    fmt.Printf("Deleted %s\n", key)
}

func listVars(file string) {
    env, err := dotenv.Load(file)
    if err != nil {
        log.Fatal(err)
    }

    for key, value := range env {
        fmt.Printf("%s=%s\n", key, value)
    }
}

func validateFile(file string) {
    _, err := dotenv.Load(file)
    if err != nil {
        log.Fatalf("Validation failed: %v", err)
    }
    fmt.Println("File is valid")
}

func splitFirst(s, sep string) []string {
    idx := strings.Index(s, sep)
    if idx == -1 {
        return []string{s}
    }
    return []string{s[:idx], s[idx+1:]}
}
```

Usage:
```bash
# List all variables
go run config-tool.go -list

# Get a value
go run config-tool.go -get DATABASE_URL

# Set a value
go run config-tool.go -set "API_KEY=secret123"

# Delete a key
go run config-tool.go -del OLD_KEY

# Validate syntax
go run config-tool.go -validate
```

## Docker Integration

Using .env files with Docker containers.

```go
package main

import (
    "log"
    "github.com/nyxstack/dotenv"
)

type Config struct {
    DatabaseURL string `env:"DATABASE_URL,required"`
    RedisURL    string `env:"REDIS_URL,required"`
    Port        int    `env:"PORT,default=8080"`
}

func main() {
    // In Docker, env vars might come from multiple sources:
    // 1. .env file
    // 2. docker-compose environment
    // 3. Docker run -e flags
    
    // Try to load .env if it exists
    if err := dotenv.LoadAndApply(".env"); err != nil {
        log.Printf("No .env file found (this is OK in Docker): %v", err)
    }

    var config Config
    if err := dotenv.Unmarshal(&config); err != nil {
        log.Fatal(err)
    }

    log.Printf("Starting application on port %d", config.Port)
    // Start application
}
```

`Dockerfile`:
```dockerfile
FROM golang:1.24-alpine AS builder
WORKDIR /app
COPY go.* ./
RUN go mod download
COPY . .
RUN go build -o app .

FROM alpine:latest
WORKDIR /app
COPY --from=builder /app/app .
# .env file is optional - use Docker env vars instead
CMD ["./app"]
```

`docker-compose.yml`:
```yaml
version: '3.8'
services:
  app:
    build: .
    ports:
      - "8080:8080"
    environment:
      - DATABASE_URL=postgresql://postgres:password@db:5432/myapp
      - REDIS_URL=redis://redis:6379
      - PORT=8080
    env_file:
      - .env.docker  # Optional: load from file
    depends_on:
      - db
      - redis

  db:
    image: postgres:15-alpine
    environment:
      - POSTGRES_PASSWORD=password
      - POSTGRES_DB=myapp

  redis:
    image: redis:7-alpine
```

## CLI Application

Command-line tool with configuration support.

```go
package main

import (
    "flag"
    "fmt"
    "log"
    "os"
    "path/filepath"
    "github.com/nyxstack/dotenv"
)

type CLIConfig struct {
    APIKey    string `env:"API_KEY,required"`
    APIBaseURL string `env:"API_BASE_URL,default=https://api.example.com"`
    Timeout   int    `env:"TIMEOUT,default=30"`
    Verbose   bool   `env:"VERBOSE,default=false"`
}

func main() {
    // Define flags
    var (
        configFile = flag.String("config", "", "Path to .env file")
        verbose    = flag.Bool("verbose", false, "Verbose output")
        apiKey     = flag.String("api-key", "", "API key (overrides config)")
    )
    flag.Parse()

    // Load configuration
    config, err := loadConfig(*configFile)
    if err != nil {
        log.Fatal(err)
    }

    // Command-line flags override config file
    if *verbose {
        config.Verbose = true
    }
    if *apiKey != "" {
        config.APIKey = *apiKey
    }

    // Run command
    if flag.NArg() < 1 {
        fmt.Println("Usage: cli [options] <command>")
        os.Exit(1)
    }

    command := flag.Arg(0)
    if err := runCommand(command, config); err != nil {
        log.Fatal(err)
    }
}

func loadConfig(configFile string) (*CLIConfig, error) {
    // Determine config file location
    if configFile == "" {
        // Try default locations
        locations := []string{
            ".env",
            filepath.Join(os.Getenv("HOME"), ".myapp.env"),
            "/etc/myapp/.env",
        }

        for _, loc := range locations {
            if _, err := os.Stat(loc); err == nil {
                configFile = loc
                break
            }
        }
    }

    // Load config file if found
    if configFile != "" {
        if err := dotenv.LoadAndApply(configFile); err != nil {
            log.Printf("Warning: failed to load %s: %v", configFile, err)
        }
    }

    // Unmarshal into struct
    var config CLIConfig
    if err := dotenv.Unmarshal(&config); err != nil {
        return nil, fmt.Errorf("configuration error: %w", err)
    }

    return &config, nil
}

func runCommand(command string, config *CLIConfig) error {
    if config.Verbose {
        log.Printf("Running command: %s", command)
        log.Printf("API Base URL: %s", config.APIBaseURL)
    }

    switch command {
    case "status":
        return checkStatus(config)
    case "deploy":
        return deploy(config)
    default:
        return fmt.Errorf("unknown command: %s", command)
    }
}

func checkStatus(config *CLIConfig) error {
    fmt.Println("Checking status...")
    // Implementation here
    return nil
}

func deploy(config *CLIConfig) error {
    fmt.Println("Deploying...")
    // Implementation here
    return nil
}
```

Usage:
```bash
# Use default .env file
./cli status

# Specify config file
./cli -config=/path/to/.env status

# Override with flags
./cli -api-key=xyz123 -verbose deploy

# Set via environment
export API_KEY=abc123
./cli status
```
