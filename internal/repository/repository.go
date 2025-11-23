package repository

import (
	"context"
	"github.com/jackc/pgx/v5/pgxpool"
	"time"
)

type Repository struct {
	Pool *pgxpool.Pool
}

func NewRepository() (Repository, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	poolConfig, err := pgxpool.ParseConfig("postgres://postgres:0000@db:5432/auth_user_db?sslmode=disable")
	if err != nil {
		return Repository{}, err
	}

	poolConfig.MaxConns = 10
	poolConfig.MinConns = 2
	poolConfig.MaxConnLifetime = 10 * time.Minute
	poolConfig.MaxConnIdleTime = 5 * time.Minute

	dbPool, err := pgxpool.NewWithConfig(ctx, poolConfig)
	if err != nil {
		return Repository{}, err
	}

	if err = dbPool.Ping(ctx); err != nil {
		dbPool.Close()
		return Repository{}, err
	}

	return Repository{Pool: dbPool}, err
}
