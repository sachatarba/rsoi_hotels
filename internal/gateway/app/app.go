package app

import (
	"github.com/gin-gonic/gin"
	"github.com/sachatarba/rsoi_hotels/config"
	"github.com/sachatarba/rsoi_hotels/internal/gateway/delivery/http"
	"github.com/sachatarba/rsoi_hotels/internal/gateway/repository"
	"github.com/sachatarba/rsoi_hotels/internal/gateway/services"
	"github.com/sachatarba/rsoi_hotels/pkg/httpserver"
	"github.com/sachatarba/rsoi_hotels/pkg/logger"
	"os"
	"os/signal"
	"syscall"
)

func Run(cfg *config.Config) {
	l := logger.New(cfg.Env)

	// Init Repos (Adapters)
	loyaltyRepo := repository.NewLoyaltyRepository(cfg.Gateway.LoyaltyUrl)
	paymentRepo := repository.NewPaymentRepository(cfg.Gateway.PaymentUrl)
	reservationRepo := repository.NewReservationRepository(cfg.Gateway.ReservationUrl)

	// Init Service
	gatewayService := services.NewGatewayService(reservationRepo, paymentRepo, loyaltyRepo, l)

	// Init HTTP
	app := gin.New()
	http.NewRouter(app, cfg, gatewayService, l)

	// Start Server
	server := httpserver.New(app, httpserver.WithPort(cfg.Gateway.Port))
	server.Start()

	l.Info("Gateway Service started", "port", cfg.Gateway.Port)

	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt, syscall.SIGTERM)

	select {
	case s := <-interrupt:
		l.Info("app - Run - signal", "signal", s.String())
	case err := <-server.Notify():
		l.Error("app - Run - httpServer.Notify:", "error", err)
	}

	err := server.Shutdown()
	if err != nil {
		l.Error("app - Run - httpServer.Shutdown:", "error", err)
	}
}
