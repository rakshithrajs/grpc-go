package main

import (
	"log/slog"
	"net"
	"os"
	"os/signal"
	"syscall"

	MMSpb "github.com/rakshithrajs/cloud/MMS/gen/MMS/v1"
	"github.com/rakshithrajs/cloud/MMS/internal/config"
	"github.com/rakshithrajs/cloud/MMS/internal/handlers"
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
		slog.Error(logPrefix+"failed to get config", slog.Any(config.ErrorKey, err))
		return
	}

	db, err := storage.Connect(cfg.DSN)
	if err != nil {
		slog.Error(logPrefix+"failed to connect to database", slog.Any(config.ErrorKey, err))
		return
	}
	slog.Info("Db Connected")
	defer func() {
		slog.Info(logPrefix + "Closing Database connection")
		if err := db.Close(); err != nil {
			slog.Error(logPrefix+"failed to close database connection", slog.Any(config.ErrorKey, err))
		}
	}()

	listen, err := net.Listen("tcp", cfg.GRPCAddress)
	if err != nil {
		slog.Error(logPrefix+"failed to listen", slog.Any(config.ErrorKey, err))
		return
	}

	server := grpc.NewServer()

	store := storage.NewFileStore(db)
	fileHandler := handlers.NewFileHandler(store)
	MMSpb.RegisterFilesServer(server, fileHandler)

	slog.Info(logPrefix+"starting MMS gRPC server", slog.String("address", cfg.GRPCAddress))

	done := make(chan os.Signal, 1)
	signal.Notify(done, os.Interrupt, syscall.SIGTERM, syscall.SIGINT)

	go func() {
		if err := server.Serve(listen); err != nil {
			slog.Error(logPrefix+"failed to serve", slog.Any(config.ErrorKey, err))
		}
	}()

	<-done

	slog.Info(logPrefix + "Shutting down server")
	server.GracefulStop()
}
