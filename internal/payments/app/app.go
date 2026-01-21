package app

import (
	"github.com/gin-gonic/gin"
	"github.com/sachatarba/rsoi_hotels/config"
	"github.com/sachatarba/rsoi_hotels/internal/payments/database/repository"
	"github.com/sachatarba/rsoi_hotels/internal/payments/delivery/http"
	"github.com/sachatarba/rsoi_hotels/internal/payments/services"
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

	paymentCfg := cfg.Payment

	timeout := time.Duration(paymentCfg.ConnTimeoutSeconds) * time.Second
	pg, err := postgres.New(paymentCfg.Url,
		postgres.WithConnAttempts(paymentCfg.ConnAttempts),
		postgres.WithConnTimeout(timeout),
		postgres.WithMaxPoolSize(paymentCfg.MaxPoolSize),
	)

	if err != nil {
		l.Error("Error creating postgres instance", "error", err)
		return
	}
	defer pg.Pool.Close()

	// Init Repo & Service
	repo := repository.NewPaymentRepository(pg.Pool, l)
	service := services.NewPaymentService(repo, l)

	// Init HTTP Server
	app := gin.New()
	http.NewRouter(app, cfg, service, l)

	server := httpserver.New(app, httpserver.WithPort(paymentCfg.Port))
	server.Start()

	l.Info("Payment Service started", "port", paymentCfg.Port)

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
