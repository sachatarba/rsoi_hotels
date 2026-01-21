package http

import (
	"github.com/gin-gonic/gin"
	"github.com/sachatarba/rsoi_hotels/config"
	"github.com/sachatarba/rsoi_hotels/internal/loyalty/delivery/http/middleware"
	v1 "github.com/sachatarba/rsoi_hotels/internal/reservations/delivery/http/v1"
	"github.com/sachatarba/rsoi_hotels/internal/reservations/domain/services"
	"log/slog"
	"net/http"
)

func NewRouter(app *gin.Engine, cfg *config.Config, bookingService services.IBookingService, logger *slog.Logger) {
	app.Use(middleware.Logger(logger))
	app.Use(gin.Recovery())

	// Health Check
	app.GET("/manage/health", func(c *gin.Context) {
		c.Status(http.StatusOK)
	})

	apiV1Group := app.Group("/api/v1")
	{
		v1.NewReservationRoutes(apiV1Group, bookingService, logger)
	}
}
