package server

import (
    "context"
    "log"
    "net/http"
    "os"
    "os/signal"
    "syscall"
    "time"
    "link-shortener/internal/handlers"
)

type Server struct {
    server *http.Server
    handler *handlers.Handler
}

func New(handler *handlers.Handler, port string) *Server {
    mux := http.NewServeMux()

    // Регистрация маршрутов через handlers
    handlers.RegisterRoutes(mux, handler)

    return &Server{
        server: &http.Server{
            Addr:    ":" + port,
            Handler: mux,
        },
        handler: handler,
    }
}

func (s *Server) Run() error {
    // Graceful shutdown
    go func() {
        sigChan := make(chan os.Signal, 1)
        signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
        <-sigChan

        log.Println("Получен сигнал завершения работы...")

        ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
        defer cancel()

        if err := s.server.Shutdown(ctx); err != nil {
            log.Printf("Ошибка при завершении работы: %v\n", err)
        }
    }()

    log.Printf("Сервер запущен на :%s\n", s.server.Addr)
    return s.server.ListenAndServe()
}