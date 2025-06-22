package middleware

import (
	"fmt"
	"log/slog"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

const expectedDurationMs = 2000

func LoggingMiddleware(log *slog.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		op := getPackageAndHandlerName(c.HandlerName()) + "." + uuid.NewString() // add uuid to track every operation individually
		c.Set("op", op)
		log := log.With(slog.String("op", op))

		log.Info("Request received", slog.String("ip", c.ClientIP()), slog.String("user_agent", c.GetHeader("User-Agent")), slog.String("path", c.FullPath()))

		start := time.Now()
		log.Debug("Start", slog.Time("start_time", start))

		c.Next()

		duration := time.Since(start).Milliseconds()
		log.Debug("Finish", slog.Int64("duration_ms", duration), slog.Int("status", c.Writer.Status()))
		if duration > expectedDurationMs {
			log.Warn(fmt.Sprintf("Operation is taking more than expected duration(%d ms)", expectedDurationMs), slog.Int64("actual duration(ms)", duration))
		}
	}
}

// this func exists because default handler name looks like this github.com/Util787/junTask/internal/handlers.(*Handler).getAllUsers-fm
func getPackageAndHandlerName(fullName string) string {
	fullName = strings.TrimSuffix(fullName, "-fm")
	parts := strings.Split(fullName, "/")
	lastPart := parts[len(parts)-1] // "handlers.(*Handler).getAllUsers"
	packageAndMethodStr := strings.ReplaceAll(lastPart, ".(*Handler)", "")
	return packageAndMethodStr
}
