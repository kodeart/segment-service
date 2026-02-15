package main

import (
    "context"
    "errors"
    "log/slog"
    "net/http"
    "os"
    "os/signal"
    "syscall"
    "time"

    "segment-service/infra/config"
    "segment-service/infra/database"
    "segment-service/internal/api"
    "segment-service/internal/repository"
    "segment-service/internal/service"
)

func main() {
    ctx, cancel := signal.NotifyContext(
        context.Background(),
        syscall.SIGINT,
        syscall.SIGTERM,
    )
    defer cancel()

    if err := run(ctx); err != nil {
        exit(err)
    }
}

func run(ctx context.Context) error {
    // Load service configuration
    cfg, err := config.Load()
    if err != nil {
        return err
    }
    // Create a database instance
    db, err := database.NewPostgresDatabase(ctx, cfg.GetPostgresDsn())
    if err != nil {
        return err
    }
    defer db.Close()

    // Create all dependencies
    segmentRepo := repository.NewSegmentsRepository(db)
    segmentService := service.NewSegmentService(segmentRepo)
    segmentHandler := api.NewSegmentHandler(segmentService)

    // Create the HTTP server and run it
    srv := &http.Server{
        Addr:    cfg.ServerAddr(),
        Handler: api.NewRouter(segmentHandler),
    }

    go func() {
        err := srv.ListenAndServe()
        if err != nil && !errors.Is(err, http.ErrServerClosed) {
            exit(err)
        }
    }()

    slog.Info("server listening on", "addr", srv.Addr, "env", cfg.Env)

    // Wait for the shutdown signal
    <-ctx.Done()

    slog.Info("stopping server gracefully...")
    shutdownCtx, shutdownCancel := context.WithTimeout(context.Background(), 10*time.Second)
    defer shutdownCancel()

    return srv.Shutdown(shutdownCtx)
}

func exit(err error) {
    slog.Error("[server]", "message", err.Error())
    os.Exit(1)
}
