package config

import (
	"fmt"
	"time"

	"github.com/spf13/viper"
)

type Config struct {
	DatabaseURL     string        `mapstructure:"database_url"`
	MaxOpenConns    int           `mapstructure:"max_open_conns"`
	MaxIdleConns    int           `mapstructure:"max_idle_conns"`
	ConnMaxLifetime time.Duration `mapstructure:"conn_max_lifetime"`
	GRPCPort        string        `mapstructure:"grpc_port"`
}

func LoadConfig() (*Config, error) {
	v := viper.New()
	v.AddConfigPath(".")
	v.SetConfigName(".env")
	v.SetConfigType("env")

	v.AutomaticEnv()

	v.SetDefault("GRPC_PORT", "50051")
	v.SetDefault("MAX_OPEN_CONNS", 25)
	v.SetDefault("MAX_IDLE_CONNS", 10)
	v.SetDefault("CONN_MAX_LIFETIME", "5m")

	if err := v.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); ok {
		} else {
			return nil, fmt.Errorf("erro lendo arquivo de configuração: %w", err)
		}
	}

	var cfg Config
	if err := v.Unmarshal(&cfg); err != nil {
		return nil, fmt.Errorf("erro ao converter configurações: %w", err)
	}

	if cfg.DatabaseURL == "" {
		return nil, fmt.Errorf("a variável DATABASE_URL é obrigatória")
	}

	return &cfg, nil
}
