package utils

import (
	"fmt"
	"os"
)

// SetFromEnvValOrFallback takes env variables and assigns to a string pointer.
// If the env variable does not exist, it will set a fall back.
// A program exit will occurr if a fallback is not set.
func SetFromEnvValOrFallback(dst *string, keys []string, fallback string) {
	for _, k := range keys {
		if v := os.Getenv(k); len(v) != 0 {
			*dst = v
			return
		}
	}

	if fallback == "" {
		fmt.Printf("Error in setting value from env: %s", keys[0])
		os.Exit(1)
	} else {
		*dst = fallback
	}
}
