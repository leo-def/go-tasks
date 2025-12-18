package env

import (
	"os"
	"testing"
)

func TestGetEnvFallback(t *testing.T) {
	os.Unsetenv("X_TEST_ENV")
	if v := GetEnv("X_TEST_ENV", "fallback"); v != "fallback" {
		t.Fatalf("expected fallback, got %s", v)
	}
}

func TestGetEnvSet(t *testing.T) {
	os.Setenv("X_TEST_ENV", "value")
	defer os.Unsetenv("X_TEST_ENV")
	if v := GetEnv("X_TEST_ENV", "fallback"); v != "value" {
		t.Fatalf("expected value, got %s", v)
	}
}

func TestGetIntEnv(t *testing.T) {
	os.Setenv("X_INT", "10")
	defer os.Unsetenv("X_INT")
	if v := GetIntEnv("X_INT", 5); v != 10 {
		t.Fatalf("expected 10, got %d", v)
	}
	os.Setenv("X_INT", "abc")
	if v := GetIntEnv("X_INT", 7); v != 7 {
		t.Fatalf("expected 7, got %d", v)
	}
}
