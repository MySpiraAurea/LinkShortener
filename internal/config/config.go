package config

import (
    "flag"
    "os"
)

type Config struct {
    Filepath     string
    DatabaseDSN  string
    Port         string
}

func Load() *Config {
    filepath := flag.String("f", "", "Путь к файлу для хранения данных")
    databaseDSN := flag.String("d", "", "Строка подключения к БД")
    port := flag.String("port", "8080", "Порт, на котором запускается сервер")

    flag.Parse()

    dsn := *databaseDSN
    if dsn == "" {
        dsn = os.Getenv("DATABASE_DSN")
    }

    return &Config{
        Filepath:    *filepath,
        DatabaseDSN: dsn,
        Port:        *port,
    }
}