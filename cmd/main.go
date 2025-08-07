package main

import (
    "log/slog"
    "os"
    
    "link-shortener/internal/application/server"
    "link-shortener/internal/config"
    "link-shortener/internal/controller"
    "link-shortener/internal/handlers"
    "link-shortener/internal/storage"
    "link-shortener/internal/repository/postgres"
)

func main() {
    logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))
    slog.SetDefault(logger)

    cfg := config.Load()

    var repo storage.LinkRepository
    var err error

    if cfg.DatabaseDSN != "" {
        logger.Info("Используется PostgreSQL", "dsn", maskDSN(cfg.DatabaseDSN))
        repo, err = postgres.New(cfg.DatabaseDSN)
        if err != nil {
            logger.Error("Ошибка подключения к PostgreSQL", "error", err)
            os.Exit(1)
        }
    } else if cfg.Filepath != "" {
        logger.Info("Используется файловое хранилище", "file", cfg.Filepath)
        repo, err = storage.NewStorage(cfg.Filepath)
        if err != nil {
            logger.Error("Ошибка инициализации файлового хранилища", "error", err)
            os.Exit(1)
        }
    } else {
        logger.Info("Используется хранилище в памяти")
        repo = storage.NewInMemory()
    }

    linkController := controller.NewLinkController(repo)
    handler := handlers.NewHandler(linkController)

    appServer := server.New(handler, cfg.Port)
    if err := appServer.Run(); err != nil {
        logger.Error("Ошибка при запуске сервера", "error", err)
        os.Exit(1)
    }
}

func maskDSN(dsn string) string {
    if len(dsn) == 0 {
        return ""
    }
    return "******"
}