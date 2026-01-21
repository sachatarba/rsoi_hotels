package v1

import (
	"github.com/gin-gonic/gin"
	"github.com/sachatarba/rsoi_hotels/internal/loyalty/delivery/http/v1/handlers"
	"github.com/sachatarba/rsoi_hotels/internal/loyalty/domain/services"
	"log/slog"
)

func NewLoyaltyRoutes(routes *gin.RouterGroup,
	service services.ILoyaltyService, logger *slog.Logger) {
	handler := handlers.LoyaltyHandler{
		Service: service,
		Logger:  logger,
	}

	loyaltyGroup := routes.Group("/loyalties")

	{
		loyaltyGroup.GET("", handler.GetLoyalty)
		loyaltyGroup.POST("/reservations", handler.AddReservations)
	}
}
