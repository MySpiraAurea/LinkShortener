package config

import (
	"flag"
	"fmt"
	"os"
)

type Config struct {
	Filepath     string
	DatabaseDSN  string
	Port         string
	BaseURL      string
}

func Load() *Config {
	filepath := flag.String("f", "", "Путь к файлу для хранения данных")
	databaseDSN := flag.String("d", "", "Строка подключения к БД")
	port := flag.String("port", "8080", "Порт, на котором запускается сервер")
	baseURL := flag.String("b", "", "Базовый URL (например, https://short.example.com)")

	flag.Parse()

	dsn := *databaseDSN
	if dsn == "" {
		dsn = os.Getenv("DATABASE_DSN")
	}

	// Автоопределение базового URL
	bURL := *baseURL
	if bURL == "" {
		host := "localhost"
		if inDocker := os.Getenv("IN_DOCKER"); inDocker != "" {
			host = "0.0.0.0"
		}
		if customHost := os.Getenv("HOST"); customHost != "" {
			host = customHost
		}
		protocol := "http"
		if os.Getenv("HTTPS") != "" {
			protocol = "https"
		}
		bURL = fmt.Sprintf("%s://%s:%s", protocol, host, *port)
	}

	return &Config{
		Filepath:    *filepath,
		DatabaseDSN: dsn,
		Port:        *port,
		BaseURL:     bURL,
	}
}