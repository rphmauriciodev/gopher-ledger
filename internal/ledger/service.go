package ledger

import (
	"context"
	"time"

	"github.com/rphmauriciodev/gopher-ledger/internal/domain"
	"github.com/rs/zerolog"
)

type LedgerService struct {
	accountRepo     domain.AccountRepository
	transactionRepo domain.TransactionRepository
	txManager       domain.TransactionManager
	logger          zerolog.Logger
}

func NewLedgerService(ar domain.AccountRepository, tr domain.TransactionRepository, tx domain.TransactionManager, l zerolog.Logger) *LedgerService {
	return &LedgerService{
		accountRepo:     ar,
		transactionRepo: tr,
		txManager:       tx,
		logger:          l.With().Str("component", "ledger-service").Logger(),
	}
}

func (s *LedgerService) ExecuteTransaction(ctx context.Context, txReq domain.Transaction) error {

	return s.txManager.Execute(ctx, func(txCtx context.Context) error {

		s.logger.Info().
			Str("account_id", txReq.AccountID).
			Int64("amount", txReq.Amount).
			Str("correlation_id", txReq.CorrelationID).
			Msg("Processando nova transação")

		acc, err := s.accountRepo.GetByID(ctx, txReq.AccountID)
		if err != nil {
			s.logger.Error().
				Err(err).
				Str("transaction_id", txReq.ID).
				Msg("Falha ao buscar conta")
			return err
		}

		if txReq.Type == domain.Debit && acc.Balance < txReq.Amount {
			err := domain.ErrInsufficientFunds
			s.logger.Error().
				Err(err).
				Str("transaction_id", txReq.ID).
				Msg("O saldo da conta não é suficiente para completar a transação")
			return err
		}

		if txReq.Type == domain.Debit {
			acc.Balance -= txReq.Amount
		} else {
			acc.Balance += txReq.Amount
		}
		acc.UpdatedAt = time.Now()

		if err := s.accountRepo.Save(ctx, acc); err != nil {
			s.logger.Error().
				Err(err).
				Str("transaction_id", txReq.ID).
				Msg("erro ao atualizar saldo")
			return err
		}

		if err := s.transactionRepo.Create(ctx, &txReq); err != nil {
			s.logger.Error().
				Err(err).
				Str("transaction_id", txReq.ID).
				Msg("erro ao gravar histórico")
			return err
		}
		s.logger.Info().
			Str("correlation_id", txReq.CorrelationID).
			Msg("Transação finalizada")
		return nil
	})
}
