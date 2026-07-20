package main

import (
	accountpb "github.com/rakshithrajs/cloud/services/account/gen/account/v1"
	"github.com/rakshithrajs/cloud/services/account/internal/interceptors"
	"github.com/rakshithrajs/cloud/services/account/internal/config"
	"github.com/rakshithrajs/cloud/services/account/internal/handlers"
	"github.com/rakshithrajs/cloud/services/account/internal/storage"
	"log/slog"
	"net"

	"google.golang.org/grpc"
)

func main() {
	cfg, err := config.GetConfig()
	if err != nil {
		slog.Error("[main]: failed to get config", slog.Any("error", err))
		return
	}

	db, err := storage.Connect(cfg.DSN)
	if err != nil {
		slog.Error("[main]: failed to connect to database", slog.Any("error", err))
		return
	}
	defer db.Close()

	listen, err := net.Listen("tcp", cfg.GRPCAddress)
	if err != nil {
		slog.Error("[main]: failed to listen", slog.Any("error", err))
		return
	}

	server := grpc.NewServer(grpc.UnaryInterceptor(interceptors.AuthInterceptor))

	store := storage.NewUserStore(db)
	accountHandler := handlers.NewAccountHandler(store)
	accountpb.RegisterAccountServer(server, accountHandler)

	slog.Info("[main]: starting account gRPC server", slog.String("address", cfg.GRPCAddress))
	if err := server.Serve(listen); err != nil {
		slog.Error("[main]: failed to serve", slog.Any("error", err))
	}
}
