package middleware

import (
	"os"
	"time"

	log "github.com/charmbracelet/log"
	"github.com/gin-gonic/gin"
)

var (
	LogGlobal = NewLogger(log.Options{
		ReportTimestamp: true,
		Prefix:          "global 🌏",
	})
)

type Logger struct {
	Log *log.Logger
}

func NewLogger(opts log.Options) *Logger {
	logger := log.NewWithOptions(os.Stdout, opts)
	return &Logger{Log: logger}
}

func (l *Logger) Middleware() func(c *gin.Context) {
	return func(c *gin.Context) {
		defer func(startTime time.Time) {
			requestMethod := c.Request.Method
			requestPath := c.Request.URL.Path
			statusCode := c.Writer.Status()
			l.Log.Info(requestMethod, "path", requestPath, "status-code", statusCode, "latency", time.Since(startTime))
		}(time.Now().UTC())
		c.Next()
	}
}
