package main

import (
	filespb "github.com/rakshithrajs/cloud/services/files/gen/files/v1"
	"github.com/rakshithrajs/cloud/services/files/internal/config"
	"github.com/rakshithrajs/cloud/services/files/internal/handlers"
	"github.com/rakshithrajs/cloud/services/files/internal/interceptors"
	"github.com/rakshithrajs/cloud/services/files/internal/storage"
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

	store := storage.NewFileStore(db)
	fileHandler := handlers.NewFileHandler(store)
	filespb.RegisterFilesServer(server, fileHandler)

	slog.Info("[main]: starting files gRPC server", slog.String("address", cfg.GRPCAddress))
	if err := server.Serve(listen); err != nil {
		slog.Error("[main]: failed to serve", slog.Any("error", err))
	}
}
