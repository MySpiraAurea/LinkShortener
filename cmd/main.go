package main

import (
    "log/slog"
    "os"

    "link-shortener/internal/application/server"
    "link-shortener/internal/controller"
    "link-shortener/internal/handlers"
    "link-shortener/internal/storage"
)

func main() {
    logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))
    slog.SetDefault(logger)

    store, err := storage.NewStorage("links.json")
    if err != nil {
        logger.Error("Не удалось загрузить хранилище", "error", err)
        os.Exit(1)
    }

    linkController := controller.NewLinkController(store)

    // Создаём новый обработчик, используя LinkController
    handler := handlers.NewHandler(linkController)

    appServer := server.New(handler, "8080")
    if err := appServer.Run(); err != nil {
        logger.Error("Ошибка при запуске сервера", "error", err)
        os.Exit(1)
    }
}