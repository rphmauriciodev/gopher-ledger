package domain

import (
	"context"
	"time"
)

type TransactionType string

const (
	Credit TransactionType = "CREDIT"
	Debit  TransactionType = "DEBIT"
)

type Transaction struct {
	ID            string
	AccountID     string
	Amount        int64
	Type          TransactionType
	CorrelationID string
	CreatedAt     time.Time
}

type TransactionRepository interface {
	Create(ctx context.Context, tx *Transaction) error
	GetByAccountID(ctx context.Context, accountID string) ([]*Transaction, error)
}

type Processor interface {
	Process(ctx context.Context, tx Transaction) error
}
