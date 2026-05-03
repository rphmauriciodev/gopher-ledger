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
