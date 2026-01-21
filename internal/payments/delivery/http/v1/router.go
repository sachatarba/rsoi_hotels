package v1

import (
	"github.com/gin-gonic/gin"
	"github.com/sachatarba/rsoi_hotels/internal/payments/delivery/http/v1/handlers"
	"github.com/sachatarba/rsoi_hotels/internal/payments/domain/services"
	"log/slog"
)

func NewPaymentRoutes(routes *gin.RouterGroup, service services.IPaymentService, logger *slog.Logger) {
	handler := handlers.PaymentHandler{
		Service: service,
		Logger:  logger,
	}

	g := routes.Group("/payments")
	{
		g.POST("", handler.Create)
		g.GET("/:paymentUid", handler.Get)
		g.DELETE("/:paymentUid", handler.Cancel)
	}
}
