package jwttoken

import (
    "os"
    "testing"

    "github.com/google/uuid"
)

func TestGenerateAndParseToken(t *testing.T) {
    os.Setenv("JWT_SECRET", "testsecret")
    defer os.Unsetenv("JWT_SECRET")
    os.Setenv("JWT_HOUR_EXPIRATION", "1")
    defer os.Unsetenv("JWT_HOUR_EXPIRATION")

    s := NewService()
    in := &AuthData{
        Id:       uuid.New(),
        Username: "tester",
        Role:     "ADMIN",
        Collaborator: CollaboratorData{
            ID:        uuid.New(),
            CompanyID: uuid.New(),
            Role:      "OWNER",
        },
    }
    tok, err := s.GenerateToken(in)
    if err != nil || tok == "" {
        t.Fatalf("generate error: %v", err)
    }
    out, err := s.ParseToken(tok)
    if err != nil {
        t.Fatalf("parse error: %v", err)
    }
    if out.Username != in.Username || out.Role != in.Role || out.Collaborator.CompanyID != in.Collaborator.CompanyID {
        t.Fatalf("mismatch: in=%+v out=%+v", in, out)
    }
}

func TestMissingSecret(t *testing.T) {
    os.Unsetenv("JWT_SECRET")
    s := NewService()
    _, err := s.GenerateToken(&AuthData{})
    if err != ErrJWTSecretNotSet {
        t.Fatalf("expected ErrJWTSecretNotSet, got %v", err)
    }
    _, err = s.ParseToken("x")
    if err != ErrJWTSecretNotSet {
        t.Fatalf("expected ErrJWTSecretNotSet, got %v", err)
    }
}
