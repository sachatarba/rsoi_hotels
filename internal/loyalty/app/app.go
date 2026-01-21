package app

import (
	"github.com/gin-gonic/gin"
	"github.com/sachatarba/rsoi_hotels/config"
	"github.com/sachatarba/rsoi_hotels/internal/loyalty/database/repository"
	"github.com/sachatarba/rsoi_hotels/internal/loyalty/delivery/http"
	"github.com/sachatarba/rsoi_hotels/internal/loyalty/services"
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

	loyaltyConfig := cfg.Loyalty

	timeout := time.Duration(loyaltyConfig.ConnTimeoutSeconds) * time.Second
	pg, err := postgres.New(loyaltyConfig.Url,
		postgres.WithConnAttempts(loyaltyConfig.ConnAttempts),
		postgres.WithConnTimeout(timeout),
		postgres.WithMaxPoolSize(loyaltyConfig.MaxPoolSize),
	)

	if err != nil {
		l.Error("Error creating postgres instance", "error", err)
		return
	}

	loyaltyRepo := repository.NewLoyaltyRepository(pg.Pool, l)

	loyaltyService := services.NewLoyaltyService(loyaltyRepo, l)

	app := gin.New()

	http.NewRouter(app, cfg, loyaltyService, l)

	server := httpserver.New(app, httpserver.WithPort(loyaltyConfig.Port))

	server.Start()

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
