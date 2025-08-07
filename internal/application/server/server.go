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
    "link-shortener/internal/middleware"
)

type Server struct {
    server *http.Server
}

func New(handler *handlers.Handler, port string) *Server {
    mux := http.NewServeMux()

    handlers.RegisterRoutes(mux, handler)

    stack := middleware.Logging(mux)

    return &Server{
        server: &http.Server{
            Addr:    ":" + port,
            Handler: stack,
        },
    }
}

func (s *Server) Run() error {
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