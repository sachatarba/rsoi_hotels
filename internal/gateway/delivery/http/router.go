package http

import (
	"github.com/gin-gonic/gin"
	"github.com/sachatarba/rsoi_hotels/config"
	"github.com/sachatarba/rsoi_hotels/internal/gateway/delivery/http/v1"
	"github.com/sachatarba/rsoi_hotels/internal/gateway/domain/services"
	"github.com/sachatarba/rsoi_hotels/internal/loyalty/delivery/http/middleware"
	"log/slog"
	"net/http"
)

func NewRouter(app *gin.Engine, cfg *config.Config, service services.IGatewayService, logger *slog.Logger) {
	app.Use(middleware.Logger(logger))
	app.Use(gin.Recovery())

	app.GET("/manage/health", func(c *gin.Context) {
		c.Status(http.StatusOK)
	})

	apiV1Group := app.Group("/api/v1")
	{
		v1.NewGatewayRoutes(apiV1Group, service, logger)
	}
}
