package logger

import (
	"bufio"
	"encoding/json"
	"os"
	"testing"
)

type logEntry struct {
	Level   string         `json:"level"`
	Message string         `json:"message"`
	Fields  map[string]any `json:"fields"`
}

func TestInfoWritesAndDebugFiltered(t *testing.T) {
	r, w, _ := os.Pipe()
	old := os.Stdout
	os.Stdout = w
	defer func() { os.Stdout = old }()

	Debug("debug", nil)
	Info("hello", map[string]any{"a": 1})

	w.Close()
	rd := bufio.NewReader(r)
	line, _, err := rd.ReadLine()
	if err != nil || len(line) == 0 {
		t.Fatalf("no output captured: %v", err)
	}
	var e logEntry
	if err := json.Unmarshal(line, &e); err != nil {
		t.Fatalf("unmarshal: %v", err)
	}
	if e.Level != "info" || e.Message != "hello" || e.Fields["a"].(float64) != 1 {
		t.Fatalf("unexpected log: %+v", e)
	}
}
