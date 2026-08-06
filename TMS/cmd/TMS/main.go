package main

import (
	"log/slog"
	"net"
	"os"
	"os/signal"
	"syscall"

	TMSpb "github.com/rakshithrajs/cloud/TMS/gen/TMS/v1"
	"github.com/rakshithrajs/cloud/TMS/internal/config"
	"github.com/rakshithrajs/cloud/TMS/internal/handlers"

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

	listen, err := net.Listen("tcp", cfg.GRPCAddress)
	if err != nil {
		slog.Error(logPrefix+"failed to listen", slog.Any(config.ErrorKey, err))
		return
	}

	server := grpc.NewServer()

	tokenHandler := handlers.NewTokenHandler()
	TMSpb.RegisterTokensServer(server, tokenHandler)

	slog.Info(logPrefix+"starting TMS gRPC server", slog.String("address", cfg.GRPCAddress))

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
