package v1

import (
	"github.com/gin-gonic/gin"
	"github.com/sachatarba/rsoi_hotels/internal/reservations/delivery/http/v1/handlers"
	"github.com/sachatarba/rsoi_hotels/internal/reservations/domain/services"
	"log/slog"
)

func NewReservationRoutes(routes *gin.RouterGroup, service services.IBookingService, logger *slog.Logger) {
	handler := handlers.BookingHandler{
		Logger:         logger,
		BookingService: service,
	}

	// Группа для отелей
	hotels := routes.Group("/hotels")
	{
		hotels.GET("", handler.GetHotels)
		hotels.GET("/:hotelUid", handler.GetHotelByUid)
	}

	// Группа для бронирований
	reservations := routes.Group("/reservations")
	{
		reservations.POST("", handler.CreateReservation)
		reservations.GET("", handler.GetUserReservations)
		reservations.GET("/:reservationUid", handler.GetReservation)
		reservations.DELETE("/:reservationUid", handler.CancelReservation)
	}
}
