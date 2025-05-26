package app

import (
	"log/slog"
	grpcapp "sso/iternal/app/grpc"
	"sso/iternal/services/auth"
	"sso/iternal/storage/sqlite"
	"time"
)

type App struct {
	GRPCSrv *grpcapp.App
}

func New(log *slog.Logger, grpcPort int, storagePath string, tokenTTL time.Duration) *App {
	storage, err := sqlite.New(storagePath)
	if err != nil {
		log.Error("failed to create storage", slog.String("path", storagePath), slog.Any("error", err))
		panic(err)
	}

	authService := auth.New(log, storage, storage, storage, tokenTTL)

	grpcApp := grpcapp.NewApp(log, authService, grpcPort)

	return &App{
		GRPCSrv: grpcApp,
	}
}
