package token

import (
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"fmt"
	"time"

	"github.com/google/uuid"
)

// GenerateRandomToken returns a URL-safe token generated from n random bytes.
// Uses base64 URL encoding without padding.
func GenerateRandomToken(n int) (string, error) {
	b := make([]byte, n)
	if _, err := rand.Read(b); err != nil {
		return "", err
	}
	// base64.RawURLEncoding avoids '+' and '/' and trailing '='
	return base64.RawURLEncoding.EncodeToString(b), nil
}

// GenerateActivationToken returns a 32-byte random token suitable for activation links.
func GenerateActivationToken() (string, error) { return GenerateRandomToken(32) }

// GenerateActivationTokenFrom deterministically generates a URL-safe token from collaborator ID and timestamp.
// This ensures uniqueness based on the (id, time) pair without relying on retries.
func GenerateActivationTokenFrom(id uuid.UUID, ts time.Time) string {
	data := fmt.Sprintf("%s:%d", id, ts.UnixNano())
	sum := sha256.Sum256([]byte(data))
	return base64.RawURLEncoding.EncodeToString(sum[:])
}
