package http

import (
	"github.com/gin-gonic/gin"
	"github.com/sachatarba/rsoi_hotels/config"
	"github.com/sachatarba/rsoi_hotels/internal/loyalty/delivery/http/middleware"
	v1 "github.com/sachatarba/rsoi_hotels/internal/loyalty/delivery/http/v1"
	"github.com/sachatarba/rsoi_hotels/internal/loyalty/domain/services"
	"log/slog"
	"net/http"
)

func NewRouter(app *gin.Engine,
	cfg *config.Config, loyaltyService services.ILoyaltyService, logger *slog.Logger) {
	app.Use(middleware.Logger(logger))

	// todo: можно заменить на самописную
	app.Use(gin.Recovery())

	// метрики

	// сваггер

	apiV1Group := app.Group("/api/v1")
	{
		v1.NewLoyaltyRoutes(apiV1Group, loyaltyService, logger)
	}

	app.GET("/manage/health", func(c *gin.Context) {
		c.Status(http.StatusOK)
	})
}
