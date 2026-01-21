package postgres

import "time"

type Option func(*Postgres)

func WithConnTimeout(connTimeout time.Duration) func(postgres *Postgres) {
	return func(postgres *Postgres) {
		postgres.connTimeout = connTimeout
	}
}

func WithConnAttempts(maxConnections int) func(postgres *Postgres) {
	return func(postgres *Postgres) {
		postgres.connAttempts = maxConnections
	}
}

func WithMaxPoolSize(maxPoolSize int) func(postgres *Postgres) {
	return func(postgres *Postgres) {
		postgres.maxPoolSize = maxPoolSize
	}
}
