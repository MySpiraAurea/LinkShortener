package main

import (
    "log/slog"
    "os"

    "link-shortener/cmd"
)

func main() {
    // Настройка логгера
    logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))
    slog.SetDefault(logger)

    // Запуск сервера
    cmd.RunServer()
}