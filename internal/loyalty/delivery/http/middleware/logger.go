package middleware

import (
	"github.com/gin-gonic/gin"
	"log/slog"
	"strings"
)

func buildLogMessage(c *gin.Context) string {
	var stringBuilder strings.Builder

	stringBuilder.WriteString(c.RemoteIP())
	stringBuilder.WriteString("-")
	stringBuilder.WriteString(c.Request.RequestURI)
	stringBuilder.WriteString("-")
	stringBuilder.WriteString(c.Request.Method)
	stringBuilder.WriteString("-")

	return stringBuilder.String()
}

func Logger(logger *slog.Logger) func(c *gin.Context) {
	return func(ctx *gin.Context) {
		ctx.Next()

		logger.Info(buildLogMessage(ctx))
	}
}
