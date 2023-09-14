package config

import (
	"flag"
	"os"
)

const (
	connStr     = "postgres://postgres:180612@localhost:5432/urldb?sslmode=disable"
	baseURLAddr = "http://localhost:8080"
	serverAddr  = ":8080"
)

// Структура Config для работы с флагами и переменными окружения.
type Config struct {
	ServerAddress  string `json:"server_address"`
	BaseURLAddress string `json:"base_url"`
	DBPath         string `json:"database_path"`
	WorkerCount    int
}

// Функция для создания нового объекта конфигурации.
func NewConfig() (*Config, error) {
	cfg := &Config{
		BaseURLAddress: "",
		ServerAddress:  "",
		DBPath:         "",
		WorkerCount:    15,
	}
	flag.StringVar(&cfg.ServerAddress, "a", "", "host to listen on")
	flag.StringVar(&cfg.BaseURLAddress, "b", "", "base url")
	flag.StringVar(&cfg.DBPath, "d", "", "database path")
	flag.Parse()
	cfg.BaseURLAddress = pickFirstNonEmpty(cfg.BaseURLAddress, os.Getenv("BASE_URL"), baseURLAddr)
	cfg.ServerAddress = pickFirstNonEmpty(cfg.ServerAddress, os.Getenv("SERVER_ADDRESS"), serverAddr)
	cfg.DBPath = pickFirstNonEmpty(cfg.DBPath, os.Getenv("DATABASE_DSN"), connStr)

	return cfg, nil
}

// Объект конфига
var Cfg Config

func pickFirstNonEmpty(strings ...string) string {
	for _, str := range strings {
		if str != "" {
			return str
		}
	}
	return ""
}
