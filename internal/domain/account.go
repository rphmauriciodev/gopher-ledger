package domain

import (
	"context"
	"errors"
	"time"
)

var (
	ErrAccountNotFound   = errors.New("conta não encontrada")
	ErrInsufficientFunds = errors.New("saldo insuficiente para a operação")
	ErrInvalidAmount     = errors.New("o valor da operação deve ser maior que zero")
)

type Account struct {
	ID        string
	OwnerID   string
	Balance   int64
	UpdatedAt time.Time
	Version   int
}

type AccountRepository interface {
	GetByID(ctx context.Context, id string) (*Account, error)
	Save(ctx context.Context, account *Account) error
}

func (a *Account) ApplyTransaction(amount int64, txType TransactionType) error {
	if amount <= 0 {
		return ErrInvalidAmount
	}

	switch txType {
	case Debit:
		if a.Balance < amount {
			return ErrInsufficientFunds
		}
		a.Balance -= amount
	case Credit:
		a.Balance += amount
	default:
		return errors.New("tipo de transação inválido")
	}

	a.UpdatedAt = time.Now()
	return nil
}
