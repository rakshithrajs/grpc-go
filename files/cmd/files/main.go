package main

import (
	"log/slog"
	"net"

	accountpb "github.com/rakshithrajs/cloud/services/account/gen/account/v1"
	filespb "github.com/rakshithrajs/cloud/services/files/gen/files/v1"
	"github.com/rakshithrajs/cloud/services/files/internal/config"
	"github.com/rakshithrajs/cloud/services/files/internal/handlers"
	"github.com/rakshithrajs/cloud/services/files/internal/interceptors"
	"github.com/rakshithrajs/cloud/services/files/internal/storage"

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

	accountConn, err := grpc.NewClient(cfg.AccountGRPCAddress, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		slog.Error(logPrefix+"failed to connect to account service", slog.Any("error", err))
		return
	}
	defer accountConn.Close()

	accountClient := accountpb.NewAccountClient(accountConn)

	listen, err := net.Listen("tcp", cfg.GRPCAddress)
	if err != nil {
		slog.Error(logPrefix+"failed to listen", slog.Any("error", err))
		return
	}

	server := grpc.NewServer(grpc.UnaryInterceptor(interceptors.NewAuthInterceptor(accountClient)))

	store := storage.NewFileStore(db)
	fileHandler := handlers.NewFileHandler(store)
	filespb.RegisterFilesServer(server, fileHandler)

	slog.Info(logPrefix+"starting files gRPC server", slog.String("address", cfg.GRPCAddress))
	if err := server.Serve(listen); err != nil {
		slog.Error(logPrefix+"failed to serve", slog.Any("error", err))
	}
}
