package cmd

import (
    "context"
    "log"
    "fmt"
    "net/http"
    "os"
    "os/signal"
    "syscall"
    "time"

    "link-shortener/internal/handlers"
    "link-shortener/internal/storage"
)

func RunServer() {
    // Инициализируем хранилище
    store, err := storage.NewStorage("links.json")
    if err != nil {
        log.Fatalf("Не удалось загрузить хранилище: %v", err)
    }

    // Создаем обработчик
    handler := handlers.NewHandler(store)

    // Настройка маршрутов
    mux := http.NewServeMux()
    log.Println("Регистрация маршрута /api/shorten")
    mux.HandleFunc("/api/shorten", handler.Shorten)
    mux.HandleFunc("/api/shorten/", handler.Shorten) // Поддержка слэша
    log.Println("Регистрация маршрута /{id}")
    mux.HandleFunc("/", handler.Redirect)

    // Тестовый маршрут
    mux.HandleFunc("/ping", func(w http.ResponseWriter, r *http.Request) {
        fmt.Fprintf(w, "pong")
    })

    // Настройка сервера
    server := &http.Server{
        Addr:    ":8080",
        Handler: mux,
    }

    // Graceful Shutdown
    go func() {
        sigChan := make(chan os.Signal, 1)
        signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
        <-sigChan

        log.Println("Получен сигнал завершения работы...")

        ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
        defer cancel()

        if err := server.Shutdown(ctx); err != nil {
            log.Printf("Ошибка при завершении работы: %v\n", err)
        }
    }()

    log.Println("Сервер запущен на :8080")
    if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
        log.Fatalf("Ошибка при запуске сервера: %v\n", err)
    }
}