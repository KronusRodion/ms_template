package app

import (
	"log/slog"
	"ms_template/internal/api/notes"
	"ms_template/internal/config"
	grpcserver "ms_template/internal/grpc"
)

type App struct {
    GRPCServer *grpcserver.App
}

func New(
    log *slog.Logger,
    cfg  *config.Config,
) *App {
	noteServer := notes.NewServer()
    grpcApp := grpcserver.New(log, noteServer, *cfg.GRPC.Port, *cfg.Prometheus.Port)

    return &App{
        GRPCServer: grpcApp,
    }
}