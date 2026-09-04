package postgres

import (
	"context"
	"fmt"
	"tgbot/internal/storage"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)

type Postgres struct {
	*pgxpool.Pool
	timeout time.Duration
}

func New(ctx context.Context) (*Postgres, error) {
	config := NewMust()

	connString := fmt.Sprintf(
		"postgres://%s:%s@%s:%s/%s?sslmode=disable",
		config.Username,
		config.Password,
		config.Host,
		config.Port,
		config.Database,
	)

	pgxconfig, err := pgxpool.ParseConfig(connString)
	if err != nil {
		return nil, fmt.Errorf("failed to parse pgx config: %w", err)
	}

	pool, err := pgxpool.NewWithConfig(ctx, pgxconfig)
	if err != nil {
		return nil, fmt.Errorf("failed to create postgres pool: %w", err)
	}

	if err := pool.Ping(ctx); err != nil {
		return nil, fmt.Errorf("failed to ping postgres pool: %w", err)
	}

	return &Postgres{
		Pool:    pool,
		timeout: config.Timeout,
	}, nil
}

func (p *Postgres) Save(page *storage.Page) error {
	//ctx, cancel := context.WithTimeout(context.Background(), p.timeout)
	//defer cancel()
	//
	//query := ``

	return nil
}

func (p *Postgres) PickRandom(userName string) (*storage.Page, error) {
	return nil, nil
}

func (p *Postgres) Remove(page *storage.Page) error {
	return nil
}

func (p *Postgres) IsExists(page *storage.Page) (bool, error) {
	return false, nil
}
