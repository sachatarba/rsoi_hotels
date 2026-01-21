package config

import (
	"fmt"
	"github.com/caarlos0/env/v11"
)

type (
	Config struct {
		Env         string `env:"ENV" envDefault:"development"`
		Loyalty     LoyaltyService
		Reservation ReservationService
		Payment     PaymentService
		Gateway     GatewayService
	}

	LoyaltyService struct {
		Port               string `env:"LOYALTY_SERVICE_PORT,required"`
		Url                string `env:"PG_LOYALTY_URL,required"`
		MaxPoolSize        int    `env:"PG_LOYALTY_MAX_POOL_SIZE,required"`
		ConnAttempts       int    `env:"PG_LOYALTY_CONN_ATTEMPTS,required"`
		ConnTimeoutSeconds int    `env:"PG_LOYALTY_CONN_TIMEOUT_SECONDS,required"`
	}

	ReservationService struct {
		Port               string `env:"RESERVATION_SERVICE_PORT,required"`
		Url                string `env:"PG_RESERVATION_URL,required"`
		MaxPoolSize        int    `env:"PG_RESERVATION_MAX_POOL_SIZE,required"`
		ConnAttempts       int    `env:"PG_RESERVATION_CONN_ATTEMPTS,required"`
		ConnTimeoutSeconds int    `env:"PG_RESERVATION_CONN_TIMEOUT_SECONDS,required"`
	}

	PaymentService struct {
		Port               string `env:"PAYMENT_SERVICE_PORT,required"`
		Url                string `env:"PG_PAYMENT_URL,required"`
		MaxPoolSize        int    `env:"PG_PAYMENT_MAX_POOL_SIZE,required"`
		ConnAttempts       int    `env:"PG_PAYMENT_CONN_ATTEMPTS,required"`
		ConnTimeoutSeconds int    `env:"PG_PAYMENT_CONN_TIMEOUT_SECONDS,required"`
	}

	GatewayService struct {
		Port string `env:"GATEWAY_SERVICE_PORT" envDefault:"8080"`

		LoyaltyUrl     string `env:"LOYALTY_SERVICE_URL" envDefault:"http://loyalty-service:8050"`
		PaymentUrl     string `env:"PAYMENT_SERVICE_URL" envDefault:"http://payment-service:8060"`
		ReservationUrl string `env:"RESERVATION_SERVICE_URL" envDefault:"http://reservation-service:8070"`
	}
)

func New() (*Config, error) {
	cfg := &Config{}
	if err := env.Parse(cfg); err != nil {
		return nil, fmt.Errorf("error parsing config: %w", err)
	}

	return cfg, nil
}
