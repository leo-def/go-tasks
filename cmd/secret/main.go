package main

import (
	"crypto/rand"
	"encoding/base64"
	"encoding/hex"
	"flag"
	"fmt"
	"os"
)

func main() {
	bytes := flag.Int("bytes", 48, "number of random bytes")
	format := flag.String("format", "base64url", "output format: base64url|base64|hex")
	env := flag.Bool("env", false, "print as JWT_SECRET=... line")
	flag.Parse()

	if *bytes <= 0 {
		fmt.Fprintln(os.Stderr, "bytes must be > 0")
		os.Exit(1)
	}

	b := make([]byte, *bytes)
	if _, err := rand.Read(b); err != nil {
		fmt.Fprintln(os.Stderr, "failed to read random bytes:", err)
		os.Exit(1)
	}

	var out string
	switch *format {
	case "base64url":
		out = base64.RawURLEncoding.EncodeToString(b)
	case "base64":
		out = base64.StdEncoding.EncodeToString(b)
	case "hex":
		out = hex.EncodeToString(b)
	default:
		fmt.Fprintln(os.Stderr, "invalid format; use base64url, base64, or hex")
		os.Exit(1)
	}

	if *env {
		fmt.Printf("JWT_SECRET=%s\n", out)
		return
	}
	fmt.Println(out)
}
