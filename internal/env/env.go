package env

import (
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

func init() {
	loadDotEnv()
}

func loadDotEnv() {
	// Try common locations so env vars work whether app is started from repo root or cmd/server.
	for _, path := range []string{".env", "../.env", "../../.env"} {
		if err := godotenv.Load(path); err == nil {
			return
		}
	}
}

func GetString(key, fallback string) string {
	val, ok := os.LookupEnv(key)
	if !ok {
		return fallback
	}

	return val
}

func GetInt(key string, fallback int) int {
	val, ok := os.LookupEnv(key)
	if !ok {
		return fallback
	}

	valAsInt, err := strconv.Atoi(val)
	if err != nil {
		return fallback
	}

	return valAsInt
}
