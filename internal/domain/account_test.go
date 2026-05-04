package domain_test

import (
	"testing"

	"github.com/rphmauriciodev/gopher-ledger/internal/domain"
	"github.com/stretchr/testify/assert"
)

func TestAccount_ApplyTransaction(t *testing.T) {
	tests := []struct {
		name            string
		initialBalance  int64
		amount          int64
		txType          domain.TransactionType
		wantErr         error
		expectedBalance int64
	}{
		{
			name:            "Deve realizar crédito com sucesso",
			initialBalance:  100,
			amount:          50,
			txType:          domain.Credit,
			wantErr:         nil,
			expectedBalance: 150,
		},
		{
			name:            "Deve realizar débito com sucesso",
			initialBalance:  100,
			amount:          40,
			txType:          domain.Debit,
			wantErr:         nil,
			expectedBalance: 60,
		},
		{
			name:            "Deve falhar se o saldo for insuficiente no débito",
			initialBalance:  10,
			amount:          50,
			txType:          domain.Debit,
			wantErr:         domain.ErrInsufficientFunds,
			expectedBalance: 10,
		},
		{
			name:            "Deve falhar se o valor for negativo",
			initialBalance:  100,
			amount:          -10,
			txType:          domain.Credit,
			wantErr:         domain.ErrInvalidAmount,
			expectedBalance: 100,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			acc := &domain.Account{Balance: tt.initialBalance}
			err := acc.ApplyTransaction(tt.amount, tt.txType)

			if tt.wantErr != nil {
				assert.ErrorIs(t, err, tt.wantErr)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.expectedBalance, acc.Balance)
			}
		})
	}
}
