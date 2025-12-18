package httpx

import (
    "net/http/httptest"
    "testing"
    "github.com/gin-gonic/gin"
)

func TestResolveUUIDParam(t *testing.T) {
    w := httptest.NewRecorder()
    c, _ := gin.CreateTestContext(w)
    c.Params = append(c.Params, gin.Param{Key: "id", Value: "123e4567-e89b-12d3-a456-426614174000"})
    if _, err := ResolveUUIDParam(c, "id"); err != nil {
        t.Fatalf("unexpected error: %v", err)
    }
    c.Params = []gin.Param{{Key: "id", Value: "bad"}}
    if _, err := ResolveUUIDParam(c, "id"); err == nil {
        t.Fatalf("expected error")
    }
}

func TestResolveUintParam(t *testing.T) {
    w := httptest.NewRecorder()
    c, _ := gin.CreateTestContext(w)
    c.Params = append(c.Params, gin.Param{Key: "n", Value: "42"})
    if v, err := ResolveUintParam(c, "n"); err != nil || v != 42 {
        t.Fatalf("unexpected: %v %d", err, v)
    }
    c.Params = []gin.Param{{Key: "n", Value: "x"}}
    if _, err := ResolveUintParam(c, "n"); err == nil {
        t.Fatalf("expected error")
    }
}
