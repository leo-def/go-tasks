package env

import (
	"os"
	"strconv"
)

func GetEnv(key, fallback string) string {
	value := os.Getenv(key)
	if value == "" {
		return fallback
	}
	return value
}

func GetIntEnv(key string, fallback int) int {
	strvalue := strconv.Itoa(fallback)
	valueStr := GetEnv(key, strvalue)
	val, err := strconv.Atoi(valueStr)
	if err != nil {
		return fallback
	}
	return val
}
