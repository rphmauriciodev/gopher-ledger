package grpc

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/rphmauriciodev/gopher-ledger/api/ledgerproto"
	"github.com/rphmauriciodev/gopher-ledger/internal/domain"
)

type LedgerHandler struct {
	ledgerproto.UnimplementedLedgerServiceServer
	txQueue chan<- domain.Transaction
}

func NewLedgerHandler(queue chan<- domain.Transaction) *LedgerHandler {
	return &LedgerHandler{txQueue: queue}
}

func (h *LedgerHandler) ProcessTransaction(ctx context.Context, req *ledgerproto.TransactionRequest) (*ledgerproto.TransactionResponse, error) {
	transactionID := uuid.New().String()
	tx := domain.Transaction{
		ID:            transactionID,
		AccountID:     req.AccountId,
		Amount:        req.Amount,
		Type:          domain.TransactionType(req.Type),
		CorrelationID: req.CorrelationId,
		CreatedAt:     time.Now(),
	}

	select {
	case h.txQueue <- tx:
		return &ledgerproto.TransactionResponse{
			Message:       "Transação aceita para processamento",
			Success:       true,
			TransactionId: transactionID,
		}, nil
	default:
		return &ledgerproto.TransactionResponse{
			Message: "Servidor sobrecarregado, tente novamente",
			Success: false,
		}, nil
	}
}
