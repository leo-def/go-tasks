package logger

import (
    "encoding/json"
    "os"
    "strings"
    "time"

    "github.com/gin-gonic/gin"
)

// Simple leveled logger with optional JSON fields.
// Levels: debug < info < warn < error

var levelOrder = map[string]int{
    "debug": 1,
    "info":  2,
    "warn":  3,
    "error": 4,
}

var currentLevel = func() int {
    lv := strings.ToLower(os.Getenv("LOG_LEVEL"))
    if lv == "" {
        lv = "info"
    }
    if v, ok := levelOrder[lv]; ok {
        return v
    }
    return levelOrder["info"]
}()

type entry struct {
    Time    string                 `json:"time"`
    Level   string                 `json:"level"`
    Message string                 `json:"message"`
    Fields  map[string]interface{} `json:"fields,omitempty"`
}

func enabled(l string) bool { return levelOrder[l] >= currentLevel }

func logJSON(l, msg string, fields map[string]interface{}) {
    if !enabled(l) {
        return
    }
    e := entry{
        Time:    time.Now().UTC().Format(time.RFC3339Nano),
        Level:   l,
        Message: msg,
        Fields:  fields,
    }
    b, _ := json.Marshal(e)
    os.Stdout.Write(append(b, '\n'))
}

func Debug(msg string, fields map[string]interface{}) { logJSON("debug", msg, fields) }
func Info(msg string, fields map[string]interface{})  { logJSON("info", msg, fields) }
func Warn(msg string, fields map[string]interface{})  { logJSON("warn", msg, fields) }
func Error(msg string, fields map[string]interface{}) { logJSON("error", msg, fields) }

// WithContext enriches logs with HTTP request context for tracing.
func WithContext(c *gin.Context, l, msg string, fields map[string]interface{}) {
    if fields == nil {
        fields = map[string]interface{}{}
    }
    // Safe, minimal request metadata (no sensitive payloads)
    fields["method"] = c.Request.Method
    fields["path"] = c.Request.URL.Path
    fields["client_ip"] = c.ClientIP()
    if rid := c.Request.Header.Get("X-Request-ID"); rid != "" {
        fields["request_id"] = rid
    }
    logJSON(l, msg, fields)
}