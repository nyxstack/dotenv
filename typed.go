package dotenv

import (
	"os"
	"strconv"
	"time"
)

// Env returns the value of an environment variable as a string
// Returns the default value if the variable is not set
func Env(key string, defaultValue ...string) string {
	value, exists := os.LookupEnv(key)
	if !exists && len(defaultValue) > 0 {
		return defaultValue[0]
	}
	return value
}

// EnvInt returns the value of an environment variable as an int
// Returns the default value if the variable is not set or cannot be parsed
func EnvInt(key string, defaultValue ...int) int {
	value, exists := os.LookupEnv(key)
	if !exists {
		if len(defaultValue) > 0 {
			return defaultValue[0]
		}
		return 0
	}

	intVal, err := strconv.Atoi(value)
	if err != nil {
		if len(defaultValue) > 0 {
			return defaultValue[0]
		}
		return 0
	}
	return intVal
}

// EnvInt8 returns the value of an environment variable as an int8
func EnvInt8(key string, defaultValue ...int8) int8 {
	value, exists := os.LookupEnv(key)
	if !exists {
		if len(defaultValue) > 0 {
			return defaultValue[0]
		}
		return 0
	}

	intVal, err := strconv.ParseInt(value, 10, 8)
	if err != nil {
		if len(defaultValue) > 0 {
			return defaultValue[0]
		}
		return 0
	}
	return int8(intVal)
}

// EnvInt16 returns the value of an environment variable as an int16
func EnvInt16(key string, defaultValue ...int16) int16 {
	value, exists := os.LookupEnv(key)
	if !exists {
		if len(defaultValue) > 0 {
			return defaultValue[0]
		}
		return 0
	}

	intVal, err := strconv.ParseInt(value, 10, 16)
	if err != nil {
		if len(defaultValue) > 0 {
			return defaultValue[0]
		}
		return 0
	}
	return int16(intVal)
}

// EnvInt32 returns the value of an environment variable as an int32
func EnvInt32(key string, defaultValue ...int32) int32 {
	value, exists := os.LookupEnv(key)
	if !exists {
		if len(defaultValue) > 0 {
			return defaultValue[0]
		}
		return 0
	}

	intVal, err := strconv.ParseInt(value, 10, 32)
	if err != nil {
		if len(defaultValue) > 0 {
			return defaultValue[0]
		}
		return 0
	}
	return int32(intVal)
}

// EnvInt64 returns the value of an environment variable as an int64
func EnvInt64(key string, defaultValue ...int64) int64 {
	value, exists := os.LookupEnv(key)
	if !exists {
		if len(defaultValue) > 0 {
			return defaultValue[0]
		}
		return 0
	}

	intVal, err := strconv.ParseInt(value, 10, 64)
	if err != nil {
		if len(defaultValue) > 0 {
			return defaultValue[0]
		}
		return 0
	}
	return intVal
}

// EnvUint returns the value of an environment variable as a uint
func EnvUint(key string, defaultValue ...uint) uint {
	value, exists := os.LookupEnv(key)
	if !exists {
		if len(defaultValue) > 0 {
			return defaultValue[0]
		}
		return 0
	}

	uintVal, err := strconv.ParseUint(value, 10, 64)
	if err != nil {
		if len(defaultValue) > 0 {
			return defaultValue[0]
		}
		return 0
	}
	return uint(uintVal)
}

// EnvUint8 returns the value of an environment variable as a uint8
func EnvUint8(key string, defaultValue ...uint8) uint8 {
	value, exists := os.LookupEnv(key)
	if !exists {
		if len(defaultValue) > 0 {
			return defaultValue[0]
		}
		return 0
	}

	uintVal, err := strconv.ParseUint(value, 10, 8)
	if err != nil {
		if len(defaultValue) > 0 {
			return defaultValue[0]
		}
		return 0
	}
	return uint8(uintVal)
}

// EnvUint16 returns the value of an environment variable as a uint16
func EnvUint16(key string, defaultValue ...uint16) uint16 {
	value, exists := os.LookupEnv(key)
	if !exists {
		if len(defaultValue) > 0 {
			return defaultValue[0]
		}
		return 0
	}

	uintVal, err := strconv.ParseUint(value, 10, 16)
	if err != nil {
		if len(defaultValue) > 0 {
			return defaultValue[0]
		}
		return 0
	}
	return uint16(uintVal)
}

// EnvUint32 returns the value of an environment variable as a uint32
func EnvUint32(key string, defaultValue ...uint32) uint32 {
	value, exists := os.LookupEnv(key)
	if !exists {
		if len(defaultValue) > 0 {
			return defaultValue[0]
		}
		return 0
	}

	uintVal, err := strconv.ParseUint(value, 10, 32)
	if err != nil {
		if len(defaultValue) > 0 {
			return defaultValue[0]
		}
		return 0
	}
	return uint32(uintVal)
}

// EnvUint64 returns the value of an environment variable as a uint64
func EnvUint64(key string, defaultValue ...uint64) uint64 {
	value, exists := os.LookupEnv(key)
	if !exists {
		if len(defaultValue) > 0 {
			return defaultValue[0]
		}
		return 0
	}

	uintVal, err := strconv.ParseUint(value, 10, 64)
	if err != nil {
		if len(defaultValue) > 0 {
			return defaultValue[0]
		}
		return 0
	}
	return uintVal
}

// EnvFloat32 returns the value of an environment variable as a float32
func EnvFloat32(key string, defaultValue ...float32) float32 {
	value, exists := os.LookupEnv(key)
	if !exists {
		if len(defaultValue) > 0 {
			return defaultValue[0]
		}
		return 0
	}

	floatVal, err := strconv.ParseFloat(value, 32)
	if err != nil {
		if len(defaultValue) > 0 {
			return defaultValue[0]
		}
		return 0
	}
	return float32(floatVal)
}

// EnvFloat64 returns the value of an environment variable as a float64
func EnvFloat64(key string, defaultValue ...float64) float64 {
	value, exists := os.LookupEnv(key)
	if !exists {
		if len(defaultValue) > 0 {
			return defaultValue[0]
		}
		return 0
	}

	floatVal, err := strconv.ParseFloat(value, 64)
	if err != nil {
		if len(defaultValue) > 0 {
			return defaultValue[0]
		}
		return 0
	}
	return floatVal
}

// EnvBool returns the value of an environment variable as a bool
// Accepts: true, false, 1, 0, yes, no, on, off (case insensitive)
func EnvBool(key string, defaultValue ...bool) bool {
	value, exists := os.LookupEnv(key)
	if !exists {
		if len(defaultValue) > 0 {
			return defaultValue[0]
		}
		return false
	}

	boolVal, err := strconv.ParseBool(value)
	if err != nil {
		if len(defaultValue) > 0 {
			return defaultValue[0]
		}
		return false
	}
	return boolVal
}

// EnvDuration returns the value of an environment variable as a time.Duration
// Accepts formats like "1h", "30m", "45s", etc.
func EnvDuration(key string, defaultValue ...time.Duration) time.Duration {
	value, exists := os.LookupEnv(key)
	if !exists {
		if len(defaultValue) > 0 {
			return defaultValue[0]
		}
		return 0
	}

	duration, err := time.ParseDuration(value)
	if err != nil {
		if len(defaultValue) > 0 {
			return defaultValue[0]
		}
		return 0
	}
	return duration
}

// SetEnv sets an environment variable
func SetEnv(key, value string) error {
	return os.Setenv(key, value)
}

// UnsetEnv unsets an environment variable
func UnsetEnv(key string) error {
	return os.Unsetenv(key)
}

// HasEnv checks if an environment variable is set
func HasEnv(key string) bool {
	_, exists := os.LookupEnv(key)
	return exists
}
