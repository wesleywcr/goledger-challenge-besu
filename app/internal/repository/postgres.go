package repository

import (
	"context"
	"database/sql"
	"fmt"

	"app/internal/model"

	_ "github.com/jackc/pgx/v5/stdlib"
)

type PostgresRepository struct {
	db *sql.DB
}

func NewPostgresRepository(postgresDSN string) (*PostgresRepository, error) {
	db, err := sql.Open("pgx", postgresDSN)
	if err != nil {
		return nil, fmt.Errorf("open database: %w", err)
	}
	repository := &PostgresRepository{db: db}
	if err := repository.initializeTable(context.Background()); err != nil {
		return nil, err
	}
	return repository, nil
}

func (repository *PostgresRepository) SaveState(ctx context.Context, value uint64) error {
	query := `
		INSERT INTO contract_state (id, value)
		VALUES (1, $1)
		ON CONFLICT (id) DO UPDATE
		SET value = EXCLUDED.value, updated_at = NOW()
	`
	if _, err := repository.db.ExecContext(ctx, query, value); err != nil {
		return fmt.Errorf("save state: %w", err)
	}
	return nil
}

func (repository *PostgresRepository) GetState(ctx context.Context) (model.ContractState, error) {
	query := `
		SELECT id, value, updated_at
		FROM contract_state
		WHERE id = 1
	`
	var state model.ContractState
	if err := repository.db.QueryRowContext(ctx, query).Scan(&state.ID, &state.Value, &state.UpdatedAt); err != nil {
		return model.ContractState{}, fmt.Errorf("get state: %w", err)
	}
	return state, nil
}

func (repository *PostgresRepository) Close() error {
	return repository.db.Close()
}

func (repository *PostgresRepository) initializeTable(ctx context.Context) error {
	query := `
		CREATE TABLE IF NOT EXISTS contract_state (
			id BIGINT PRIMARY KEY,
			value BIGINT NOT NULL,
			updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
		)
	`
	if _, err := repository.db.ExecContext(ctx, query); err != nil {
		return fmt.Errorf("create table contract_state: %w", err)
	}
	return nil
}
