package main

import (
	"log/slog"
	"net"

	MMSpb "github.com/rakshithrajs/cloud/MMS/gen/MMS/v1"
	"github.com/rakshithrajs/cloud/MMS/internal/config"
	"github.com/rakshithrajs/cloud/MMS/internal/handlers"
	"github.com/rakshithrajs/cloud/MMS/internal/interceptors"
	"github.com/rakshithrajs/cloud/MMS/internal/storage"

	"google.golang.org/grpc"
)

const (
	functionName = "main"
	logPrefix    = "[" + functionName + "]: "
)

func main() {
	cfg, err := config.GetConfig()
	if err != nil {
		slog.Error(logPrefix+"failed to get config", slog.Any("error", err))
		return
	}

	db, err := storage.Connect(cfg.DSN)
	if err != nil {
		slog.Error(logPrefix+"failed to connect to database", slog.Any("error", err))
		return
	}
	defer db.Close()

	listen, err := net.Listen("tcp", cfg.GRPCAddress)
	if err != nil {
		slog.Error(logPrefix+"failed to listen", slog.Any("error", err))
		return
	}

	server := grpc.NewServer(grpc.UnaryInterceptor(interceptors.AuthInterceptor()))

	store := storage.NewFileStore(db)
	fileHandler := handlers.NewFileHandler(store)
	MMSpb.RegisterFilesServer(server, fileHandler)

	slog.Info(logPrefix+"starting MMS gRPC server", slog.String("address", cfg.GRPCAddress))
	if err := server.Serve(listen); err != nil {
		slog.Error(logPrefix+"failed to serve", slog.Any("error", err))
	}
}
