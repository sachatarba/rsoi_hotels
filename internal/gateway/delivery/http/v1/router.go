package v1

import (
	"github.com/gin-gonic/gin"
	"github.com/sachatarba/rsoi_hotels/internal/gateway/delivery/http/v1/handlers"
	"github.com/sachatarba/rsoi_hotels/internal/gateway/domain/services"
	"log/slog"
)

func NewGatewayRoutes(routes *gin.RouterGroup, service services.IGatewayService, logger *slog.Logger) {
	h := handlers.NewGatewayHandler(service, logger)

	routes.GET("/hotels", h.GetHotels)
	routes.GET("/me", h.GetUserInfo)
	routes.GET("/reservations/:reservationUid", h.GetReservation)
	routes.POST("/reservations", h.BookHotel)
	routes.DELETE("/reservations/:reservationUid", h.CancelReservation)
	routes.GET("/loyalty", h.GetLoyalty)
	routes.GET("/reservations", h.GetUserReservations)
}
