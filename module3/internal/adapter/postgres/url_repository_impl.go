package postgres

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"module3/internal/entities"
	"module3/internal/repository"
	"time"

	"github.com/lib/pq"
	_ "github.com/lib/pq"
)

type PostgresRepository struct {
	db *sql.DB
}

func NewPostgresRepository(db *sql.DB) *PostgresRepository {
	return &PostgresRepository{db: db}
}

func (pr *PostgresRepository) Save(url, alias string) (*entities.Url, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	var savedId int64
	query := "INSERT INTO urls (url, alias) VALUES ($1, $2) RETURNING id;"

	err := pr.db.QueryRowContext(ctx, query, url, alias).Scan(&savedId)
	if err != nil {
		if pqErr, ok := err.(*pq.Error); ok && pqErr.Code == "23505" {
			return nil, repository.ErrAliasAlreadyExists
		}
		return nil, fmt.Errorf("db error: %w", err)
	}

	return &entities.Url{
		Id:    savedId,
		Url:   url,
		Alias: alias,
	}, nil
}

func (pr *PostgresRepository) Update(url *entities.Url) (*entities.Url, error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	var updatedId int64
	query := "UPDATE urls SET url =  $1, alias = $2 WHERE id = $3 RETURNING id;"

	err := pr.db.QueryRowContext(ctx, query, url.Url, url.Alias, url.Id).Scan(&updatedId)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, repository.ErrNotFound
		}
		return nil, fmt.Errorf("db error: %w", err)
	}

	return &entities.Url{
		Id:    updatedId,
		Url:   url.Url,
		Alias: url.Alias,
	}, nil
}

func (pr *PostgresRepository) DeleteById(id int64) error {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	query := "DELETE FROM urls WHERE id = $1"

	_, err := pr.db.ExecContext(ctx, query, id)
	if err != nil {
		return fmt.Errorf("db error: %w", err)
	}

	return nil
}

func (pr *PostgresRepository) FindById(id int64) (*entities.Url, error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	var url entities.Url
	query := "SELECT id, url, alias FROM urls where id = $1"

	err := pr.db.QueryRowContext(ctx, query, id).Scan(&url.Id, &url.Url, &url.Alias)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, repository.ErrNotFound
		}
		return nil, fmt.Errorf("db error: %w", err)
	}

	return &url, nil
}

func (pr *PostgresRepository) FindUrlStringByAlias(alias string) (string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	var url string
	query := "SELECT url FROM urls where alias = $1"

	err := pr.db.QueryRowContext(ctx, query, alias).Scan(&url)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return "", repository.ErrNotFound
		}
		return "", fmt.Errorf("db error: %w", err)
	}

	return url, nil
}
