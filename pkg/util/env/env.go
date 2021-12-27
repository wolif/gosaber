package env

import "os"

func Get(key string, fallback ...string) string {
	if val := os.Getenv(key); val != "" {
		return val
	} else if len(fallback) > 0 {
		return fallback[0]
	}
	return ""
}

func Set(key, value string) {
	os.Setenv(key, value)
}
