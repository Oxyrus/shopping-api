package utils

import "os"

// GetEnv provides a fallback mechanism to retrieve environment
// variables and provide a fallback value in case the
// environment variable isn't found.
func GetEnv(key, fallback string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}
	return fallback
}
