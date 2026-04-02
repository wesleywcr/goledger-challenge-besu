package config

import (
	"fmt"
	"os"
)

type Config struct {
	HttpPort        string
	RpcURL          string
	ContractAddress string
	PrivateKey      string
	PostgresDSN     string
}

func LoadConfig() (Config, error) {
	cfg := Config{
		HttpPort:        getEnvOrDefault("HTTP_PORT", "8080"),
		RpcURL:          os.Getenv("BESU_RPC_URL"),
		ContractAddress: os.Getenv("CONTRACT_ADDRESS"),
		PrivateKey:      os.Getenv("PRIVATE_KEY"),
		PostgresDSN:     os.Getenv("POSTGRES_DSN"),
	}
	if cfg.RpcURL == "" {
		return Config{}, fmt.Errorf("missing BESU_RPC_URL")
	}
	if cfg.ContractAddress == "" {
		return Config{}, fmt.Errorf("missing CONTRACT_ADDRESS")
	}
	if cfg.PrivateKey == "" {
		return Config{}, fmt.Errorf("missing PRIVATE_KEY")
	}
	if cfg.PostgresDSN == "" {
		return Config{}, fmt.Errorf("missing POSTGRES_DSN")
	}
	return cfg, nil
}

func getEnvOrDefault(key string, fallback string) string {
	value := os.Getenv(key)
	if value == "" {
		return fallback
	}
	return value
}
