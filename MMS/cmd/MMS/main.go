package main

import (
	"log/slog"
	"net"

	MMSpb "github.com/rakshithrajs/cloud/MMS/gen/MMS/v1"
	"github.com/rakshithrajs/cloud/MMS/internal/config"
	"github.com/rakshithrajs/cloud/MMS/internal/handlers"
	"github.com/rakshithrajs/cloud/MMS/internal/interceptors"
	"github.com/rakshithrajs/cloud/MMS/internal/storage"
	UMSpb "github.com/rakshithrajs/cloud/UMS/gen/UMS/v1"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
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

	UMSConn, err := grpc.NewClient(cfg.UMSGRPCAddress, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		slog.Error(logPrefix+"failed to connect to UMS service", slog.Any("error", err))
		return
	}
	defer UMSConn.Close()

	UMSClient := UMSpb.NewUMSClient(UMSConn)

	listen, err := net.Listen("tcp", cfg.GRPCAddress)
	if err != nil {
		slog.Error(logPrefix+"failed to listen", slog.Any("error", err))
		return
	}

	server := grpc.NewServer(grpc.UnaryInterceptor(interceptors.NewAuthInterceptor(UMSClient)))

	store := storage.NewMMStore(db)
	fileHandler := handlers.NewFileHandler(store)
	MMSpb.RegisterFilesServer(server, fileHandler)

	slog.Info(logPrefix+"starting MMS gRPC server", slog.String("address", cfg.GRPCAddress))
	if err := server.Serve(listen); err != nil {
		slog.Error(logPrefix+"failed to serve", slog.Any("error", err))
	}
}
