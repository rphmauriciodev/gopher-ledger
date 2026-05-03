package database

import (
	"context"
	"database/sql"
	"time"

	_ "github.com/lib/pq"
	"github.com/rphmauriciodev/gopher-ledger/internal/config"
	"github.com/rs/zerolog"
)

func NewPostgresConnection(cfg config.Config, logger zerolog.Logger) (*sql.DB, error) {
	db, err := sql.Open("postgres", cfg.DatabaseURL)
	if err != nil {
		logger.Error().Err(err).Msg("Erro ao parsear URI do banco")
		return nil, err
	}
	db.SetMaxOpenConns(cfg.MaxOpenConns)

	db.SetMaxIdleConns(cfg.MaxIdleConns)

	db.SetConnMaxLifetime(cfg.ConnMaxLifetime)

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := db.PingContext(ctx); err != nil {
		logger.Error().Err(err).Msg("Não foi possível pingar o banco de dados")
		return nil, err
	}

	logger.Info().
		Int("max_open", cfg.MaxOpenConns).
		Int("max_idle", cfg.MaxIdleConns).
		Msg("Conexão com Postgres estabelecida com sucesso")

	return db, nil
}
