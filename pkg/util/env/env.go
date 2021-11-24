package env

import "os"

func GetEnv(key string) string {
	return os.Getenv(key)
}

func GetEnvWithFallback(key, fallback string) string {
	val := os.Getenv(key)
	if val != "" {
		return val
	}

	return fallback
}

func SetEnv(key, value string) {
	os.Setenv(key, value)
}
