package postgres

import (
	"context"
	"database/sql"
)

type txKey struct{}

type PostgresTransactionManager struct {
	db *sql.DB
}

type executor interface {
	QueryRowContext(ctx context.Context, query string, args ...interface{}) *sql.Row
	QueryContext(ctx context.Context, query string, args ...interface{}) (*sql.Rows, error)
	ExecContext(ctx context.Context, query string, args ...interface{}) (sql.Result, error)
}

func NewPostgresTransactionManager(db *sql.DB) *PostgresTransactionManager {
	return &PostgresTransactionManager{db: db}
}

func (m *PostgresTransactionManager) Execute(ctx context.Context, fn func(ctx context.Context) error) error {
	tx, err := m.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}

	defer tx.Rollback()

	txCtx := context.WithValue(ctx, txKey{}, tx)

	if err := fn(txCtx); err != nil {
		return err
	}

	if err := tx.Commit(); err != nil {
		return err
	}

	return nil
}
