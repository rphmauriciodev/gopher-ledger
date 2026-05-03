package postgres

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/rphmauriciodev/gopher-ledger/internal/domain"
)

type PostgresAccountRepository struct {
	db *sql.DB
}

func NewPostgresAccountRepository(db *sql.DB) *PostgresAccountRepository {
	return &PostgresAccountRepository{db: db}
}

func (r *PostgresAccountRepository) getExecutor(ctx context.Context) executor {
	if tx, ok := ctx.Value(txKey{}).(*sql.Tx); ok {
		return tx
	}
	return r.db
}

func (r *PostgresAccountRepository) GetByID(ctx context.Context, id string) (*domain.Account, error) {
	const query = `SELECT id, balance, version FROM accounts WHERE id = $1`

	var acc domain.Account
	err := r.getExecutor(ctx).QueryRowContext(ctx, query, id).Scan(&acc.ID, &acc.Balance, &acc.Version)

	if err == sql.ErrNoRows {
		return nil, domain.ErrAccountNotFound
	}
	return &acc, err
}

func (r *PostgresAccountRepository) Save(ctx context.Context, acc *domain.Account) error {
	const query = `UPDATE accounts SET balance = $1, version = version + 1 WHERE id = $2 AND version = $3`

	result, err := r.getExecutor(ctx).ExecContext(ctx, query, acc.Balance, acc.ID, acc.Version)
	if err != nil {
		return err
	}

	rows, _ := result.RowsAffected()
	if rows == 0 {
		return fmt.Errorf("concorrência detectada: a conta foi alterada por outro processo")
	}
	return nil
}
