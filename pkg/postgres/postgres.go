package postgres

import (
	"context"
	"fmt"
	"github.com/jackc/pgx/v5/pgxpool"
	"log"
	"time"
)

const (
	_defaultMaxPoolSize  = 1
	_defaultConnAttempts = 10
	_defaultConnTimeout  = time.Second
)

type Postgres struct {
	maxPoolSize  int
	connAttempts int
	connTimeout  time.Duration

	Pool *pgxpool.Pool
}

func New(url string, opts ...Option) (*Postgres, error) {
	pg := &Postgres{
		maxPoolSize:  _defaultMaxPoolSize,
		connAttempts: _defaultConnAttempts,
		connTimeout:  _defaultConnTimeout,
	}

	for _, opt := range opts {
		opt(pg)
	}

	if pg.connAttempts == 0 {
		pg.connAttempts = _defaultConnAttempts
	}

	poolConfig, err := pgxpool.ParseConfig(url)
	if err != nil {
		return nil, fmt.Errorf("postgres: error parsing connection url %s: %w", url, err)
	}
	poolConfig.MaxConns = int32(pg.maxPoolSize)

	for pg.connAttempts > 0 {
		pg.Pool, err = pgxpool.NewWithConfig(context.Background(), poolConfig)
		if err == nil {
			err = pg.Pool.Ping(context.Background())
			if err == nil {
				return pg, nil
			}
			pg.Pool.Close()
		}

		log.Printf("Ping postgres error: trying to connect to postgres at %s, attempts left %d. Error: %v", url, pg.connAttempts, err)

		time.Sleep(pg.connTimeout)

		pg.connAttempts--
	}

	return nil, fmt.Errorf("postgres: error creating connection pool: %w", err)
}
