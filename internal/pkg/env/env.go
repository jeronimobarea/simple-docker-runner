package env

import "os"

func GetEnvWithFallback(key, fallback string) string {
	if res, ok := os.LookupEnv(key); ok {
		return res
	}
	return fallback
}
