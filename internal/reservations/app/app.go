package app

import (
	"github.com/gin-gonic/gin"
	"github.com/sachatarba/rsoi_hotels/config"
	repository "github.com/sachatarba/rsoi_hotels/internal/reservations/database/repository"
	"github.com/sachatarba/rsoi_hotels/internal/reservations/delivery/http"
	"github.com/sachatarba/rsoi_hotels/internal/reservations/services"
	"github.com/sachatarba/rsoi_hotels/pkg/httpserver"
	"github.com/sachatarba/rsoi_hotels/pkg/logger"
	"github.com/sachatarba/rsoi_hotels/pkg/postgres"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func Run(cfg *config.Config) {
	l := logger.New(cfg.Env)

	reservationCfg := cfg.Reservation

	timeout := time.Duration(reservationCfg.ConnTimeoutSeconds) * time.Second
	pg, err := postgres.New(reservationCfg.Url,
		postgres.WithConnAttempts(reservationCfg.ConnAttempts),
		postgres.WithConnTimeout(timeout),
		postgres.WithMaxPoolSize(reservationCfg.MaxPoolSize),
	)

	if err != nil {
		l.Error("Error creating postgres instance", "error", err)
		return
	}
	defer pg.Pool.Close()

	// Init Repos
	hotelRepo := repository.NewHotelRepository(pg.Pool, l)
	reservationRepo := repository.NewReservationRepository(pg.Pool, l)

	// Init Service
	bookingService := services.NewBookingService(reservationRepo, hotelRepo, l)

	// Init HTTP Server
	ginApp := gin.New()
	http.NewRouter(ginApp, cfg, bookingService, l)

	server := httpserver.New(ginApp, httpserver.WithPort(reservationCfg.Port))
	server.Start()

	l.Info("Reservation Service started", "port", reservationCfg.Port)

	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt, syscall.SIGTERM)

	select {
	case s := <-interrupt:
		l.Info("app - Run - signal", "signal", s.String())
	case err = <-server.Notify():
		l.Error("app - Run - httpServer.Notify:", "error", err)
	}

	err = server.Shutdown()
	if err != nil {
		l.Error("app - Run - httpServer.Shutdown:", "error", err)
	}
}
