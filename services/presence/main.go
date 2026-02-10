package main

import (
	"context"
	"errors"
	"log/slog"
	"os"
	"os/signal"
	"syscall"
	"time"

	"google.golang.org/grpc"

	"github.com/adrputra/face-recognition-svc/presence-svc/internal/config"
	"github.com/adrputra/face-recognition-svc/presence-svc/internal/server"
)

const shutdownTimeout = 10 * time.Second

func main() {
	logger := slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
		Level: slog.LevelInfo,
	}))
	slog.SetDefault(logger)

	configPath := os.Getenv("CONFIG_PATH")
	if configPath == "" {
		configPath = config.DefaultConfigPath
	}

	cfg, err := config.Load(configPath)
	if err != nil {
		slog.Error("Failed to load config", "error", err)
		os.Exit(1)
	}
	slog.Info("Config loaded", "path", configPath, "grpc_port", cfg.GRPC.Port)

	grpcServer, err := server.NewGRPCServer(cfg)
	if err != nil {
		slog.Error("Failed to create gRPC server", "error", err)
		os.Exit(1)
	}

	errCh := make(chan error, 1)
	go func() {
		slog.Info("gRPC server started", "addr", grpcServer.Addr())
		errCh <- grpcServer.Start()
	}()

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	select {
	case <-ctx.Done():
		slog.Info("Shutdown initiated")
	case err := <-errCh:
		if err != nil && !errors.Is(err, grpc.ErrServerStopped) {
			slog.Error("gRPC server stopped unexpectedly", "error", err)
		}
	}

	shutdownCtx, cancel := context.WithTimeout(context.Background(), shutdownTimeout)
	defer cancel()
	if err := grpcServer.Stop(shutdownCtx); err != nil && !errors.Is(err, context.Canceled) {
		slog.Error("gRPC server shutdown error", "error", err)
	}
}
