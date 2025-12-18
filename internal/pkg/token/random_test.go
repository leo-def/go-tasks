package token

import (
    "encoding/base64"
    "testing"
    "time"
    "github.com/google/uuid"
)

func TestGenerateRandomToken(t *testing.T) {
    tok, err := GenerateRandomToken(32)
    if err != nil || tok == "" {
        t.Fatalf("error: %v", err)
    }
    if _, err := base64.RawURLEncoding.DecodeString(tok); err != nil {
        t.Fatalf("not url-safe base64: %v", err)
    }
}

func TestGenerateActivationTokenFrom(t *testing.T) {
    id := uuid.New()
    ts := time.Unix(1700000000, 123)
    a := GenerateActivationTokenFrom(id, ts)
    b := GenerateActivationTokenFrom(id, ts)
    if a != b {
        t.Fatalf("not deterministic")
    }
}
