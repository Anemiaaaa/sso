package app

import (
	"log/slog"
	grpcapp "sso/iternal/app/grpc"
	"time"
)

type App struct {
	GRPCSrv *grpcapp.App
}

func New(log *slog.Logger, grpcPort int, storagePath string, tokenTTL time.Duration) *App {
	// Создаем новый gRPC сервер
	grpcApp := grpcapp.NewApp(log, grpcPort)

	return &App{
		GRPCSrv: grpcApp,
	}
}
