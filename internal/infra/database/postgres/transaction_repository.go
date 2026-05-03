package postgres

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/rphmauriciodev/gopher-ledger/internal/domain"
)

type PostgresTransactionRepository struct {
	db *sql.DB
}

func NewPostgresTransactionRepository(db *sql.DB) *PostgresTransactionRepository {
	return &PostgresTransactionRepository{db: db}
}

func (r *PostgresTransactionRepository) getExecutor(ctx context.Context) executor {
	if tx, ok := ctx.Value(txKey{}).(*sql.Tx); ok {
		return tx
	}
	return r.db
}

func (r *PostgresTransactionRepository) Create(ctx context.Context, tx *domain.Transaction) error {
	const query = `
		INSERT INTO transactions (id, account_id, amount, type, correlation_id, created_at)
		VALUES ($1, $2, $3, $4, $5, $6)
	`

	_, err := r.getExecutor(ctx).ExecContext(
		ctx,
		query,
		tx.ID,
		tx.AccountID,
		tx.Amount,
		tx.Type,
		tx.CorrelationID,
		tx.CreatedAt,
	)

	if err != nil {
		return fmt.Errorf("falha ao inserir transação no banco: %w", err)
	}

	return nil
}

func (r *PostgresTransactionRepository) GetByAccountID(ctx context.Context, accountID string) ([]*domain.Transaction, error) {
	const query = `
		SELECT id, account_id, amount, type, correlation_id, created_at 
		FROM transactions 
		WHERE account_id = $1 
		ORDER BY created_at DESC
	`

	rows, err := r.getExecutor(ctx).QueryContext(ctx, query, accountID)
	if err != nil {
		return nil, fmt.Errorf("falha ao buscar histórico: %w", err)
	}
	defer rows.Close()

	var transactions []*domain.Transaction
	for rows.Next() {
		var t domain.Transaction
		err := rows.Scan(
			&t.ID,
			&t.AccountID,
			&t.Amount,
			&t.Type,
			&t.CorrelationID,
			&t.CreatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("falha ao ler linha da transação: %w", err)
		}
		transactions = append(transactions, &t)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return transactions, nil
}
