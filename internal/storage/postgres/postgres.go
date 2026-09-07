package postgres

import (
	"context"
	"errors"
	"fmt"
	"tgbot/internal/errs"
	"tgbot/internal/storage"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
)

const ErrViolatesForeignKey = "23503"

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

func (p *Postgres) SavePage(ctx context.Context, page storage.Page) error {
	ctx, cancel := context.WithTimeout(ctx, p.timeout)
	defer cancel()

	query := `
	INSERT INTO tgbot.pages (url, user_id)
	VALUES ($1, $2)
	`

	_, err := p.Exec(ctx, query, page.URL, page.UserID)
	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgErr.Code == ErrViolatesForeignKey {
			return fmt.Errorf("user %d does not exist: %w", page.UserID, errs.ErrUserNotFound)
		}
		return fmt.Errorf("save page: %w", err)
	}

	return nil
}

func (p *Postgres) GetPage(ctx context.Context, page storage.Page) (storage.Page, error) {
	ctx, cancel := context.WithTimeout(ctx, p.timeout)
	defer cancel()

	query := `
	SELECT id, url, user_id FROM tgbot.pages
	WHERE url = $1 AND user_id = $2
	`

	row := p.QueryRow(ctx, query, page.URL, page.UserID)
	var pageModel storage.Page
	if err := row.Scan(&pageModel.ID, &pageModel.URL, &pageModel.UserID); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return storage.Page{}, errs.ErrPageNotFound
		}
		return storage.Page{}, fmt.Errorf("get page: %w", err)
	}
	return pageModel, nil

}

func (p *Postgres) CheckPage(ctx context.Context, page storage.Page) (bool, error) {
	ctx, cancel := context.WithTimeout(ctx, p.timeout)
	defer cancel()

	query := `
	SELECT EXISTS(SELECT id, url, user_id FROM tgbot.pages
	WHERE url = $1 AND user_id = $2)
	`

	var exists bool
	if err := p.QueryRow(ctx, query, page.URL, page.UserID).Scan(&exists); err != nil {
		return false, fmt.Errorf("check page: %w", err)
	}

	return exists, nil
}

func (p *Postgres) SaveUser(ctx context.Context, userName string) (storage.User, error) {
	ctx, cancel := context.WithTimeout(ctx, p.timeout)
	defer cancel()

	query := `
	INSERT INTO tgbot.users (username) VALUES ($1)
	RETURNING id, username
	`

	row := p.QueryRow(ctx, query, userName)
	var user storage.User
	if err := row.Scan(&user.ID, &user.Username); err != nil {
		return storage.User{}, fmt.Errorf("save user: %w", err)
	}

	return user, nil
}

func (p *Postgres) GetUser(ctx context.Context, userName string) (storage.User, error) {
	ctx, cancel := context.WithTimeout(ctx, p.timeout)
	defer cancel()

	query := `
	SELECT id, username FROM tgbot.users
	WHERE username = $1
	`

	row := p.QueryRow(ctx, query, userName)
	var user storage.User
	err := row.Scan(&user.ID, &user.Username)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return storage.User{}, errs.ErrUserNotFound
		}
		return storage.User{}, fmt.Errorf("get user: %w", err)
	}

	return user, nil
}

func (p *Postgres) CheckUser(ctx context.Context, userName string) (bool, error) {
	ctx, cancel := context.WithTimeout(ctx, p.timeout)
	defer cancel()

	query := `
	SELECT EXISTS(SELECT id, username FROM tgbot.users
	WHERE username = $1)
	`

	var exists bool
	if err := p.QueryRow(ctx, query, userName).Scan(&exists); err != nil {
		return false, fmt.Errorf("check user: %w", err)
	}

	return exists, nil
}

func (p *Postgres) PickRandomPage(ctx context.Context, user storage.User) (storage.Page, error) {
	ctx, cancel := context.WithTimeout(ctx, p.timeout)
	defer cancel()

	query := `
	SELECT pages.id, pages.url, pages.user_id FROM tgbot.pages
	JOIN tgbot.users users on users.id = pages.user_id
	WHERE users.username = $1
	ORDER BY RANDOM() LIMIT 1
	`

	row := p.QueryRow(ctx, query, user.Username)
	var page storage.Page
	if err := row.Scan(&page.ID, &page.URL, &page.UserID); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return storage.Page{}, errs.ErrPageNotFound
		}
		return storage.Page{}, fmt.Errorf("pick random page: %w", err)
	}

	return page, nil
}

func (p *Postgres) Remove(ctx context.Context, page storage.Page) error {
	ctx, cancel := context.WithTimeout(ctx, p.timeout)
	defer cancel()

	query := `
	DELETE FROM tgbot.pages
	WHERE id = $1
	`

	_, err := p.Exec(ctx, query, page.ID)
	if err != nil {
		return fmt.Errorf("remove page: %w", err)
	}

	return nil
}
