package main

import (
	"context"
	"fmt"
	"net"
	"os"
	"os/signal"
	"syscall"
	"time"

	pb "github.com/rphmauriciodev/gopher-ledger/api/ledgerproto"
	"github.com/rphmauriciodev/gopher-ledger/internal/config"
	"github.com/rphmauriciodev/gopher-ledger/internal/domain"
	"github.com/rphmauriciodev/gopher-ledger/internal/infra/database"
	pg "github.com/rphmauriciodev/gopher-ledger/internal/infra/database/postgres"
	grpcHandler "github.com/rphmauriciodev/gopher-ledger/internal/infra/grpc"
	"github.com/rphmauriciodev/gopher-ledger/internal/ledger"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"google.golang.org/grpc"
)

func main() {
	zerolog.TimeFieldFormat = zerolog.TimeFormatUnix
	logger := log.With().Str("app", "gopher-ledger").Logger()

	cfg, err := config.LoadConfig()
	if err != nil {
		logger.Fatal().Err(err).Msg("Falha crítica ao carregar configurações")
	}

	db, err := database.NewPostgresConnection(*cfg, logger)
	if err != nil {
		logger.Fatal().Err(err).Msg("Falha crítica no boot do banco de dados")
	}
	defer db.Close()

	accountRepo := pg.NewPostgresAccountRepository(db)
	transactionRepo := pg.NewPostgresTransactionRepository(db)
	txManager := pg.NewPostgresTransactionManager(db)

	txQueue := make(chan domain.Transaction, 100)

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	service := ledger.NewLedgerService(accountRepo, transactionRepo, txManager, logger)

	wp := ledger.NewWorkerPool(service, txQueue, 10)
	go func() {
		logger.Info().Int("workers", 10).Msg("Worker Pool iniciado")
		wp.Start(ctx)
	}()

	server := grpc.NewServer()
	handler := grpcHandler.NewLedgerHandler(txQueue)
	pb.RegisterLedgerServiceServer(server, handler)

	lis, err := net.Listen("tcp", fmt.Sprintf(":%s", cfg.GRPCPort))
	if err != nil {
		logger.Fatal().Err(err).Msg("Falha ao abrir porta gRPC")
	}

	go func() {
		logger.Info().Str("port", cfg.GRPCPort).Msg("Servidor gRPC rodando")
		if err := server.Serve(lis); err != nil && err != grpc.ErrServerStopped {
			logger.Fatal().Err(err).Msg("Erro inesperado no servidor gRPC")
		}
	}()
	<-ctx.Done()
	logger.Info().Msg("Sinal de interrupção recebido. Desligando gopher-ledger graciosamente...")

	server.GracefulStop()
	close(txQueue)

	time.Sleep(time.Second * 2)
	logger.Info().Msg("Gopher-ledger encerrado com sucesso.")
}
