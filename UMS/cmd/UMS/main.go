package main

import (
	"log/slog"
	"net"

	"github.com/gin-gonic/gin"
	UMSpb "github.com/rakshithrajs/cloud/UMS/gen/UMS/v1"

	"github.com/rakshithrajs/cloud/UMS/internal/config"
	"github.com/rakshithrajs/cloud/UMS/internal/handlers"
	"github.com/rakshithrajs/cloud/UMS/internal/interceptors"
	"github.com/rakshithrajs/cloud/UMS/internal/storage"

	"google.golang.org/grpc"
)

const (
	functionName = "main"
	logPrefix    = "[" + functionName + "]: "
	apiPrefix    = "/api"
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

	gin.SetMode(gin.DebugMode)
	router := gin.New()
	router.SetTrustedProxies(nil)

	router.Use(gin.Logger())

	store := storage.NewUserStore(db)
	UMSHandler := handlers.NewUMSHandler(store)

	UMSRouterGroup := router.Group(apiPrefix + "/user")
	handlers.RegisterRoutes(UMSRouterGroup, UMSHandler)

	listen, err := net.Listen("tcp", cfg.GRPCAddress)
	if err != nil {
		slog.Error(logPrefix+"failed to listen", slog.Any("error", err))
		return
	}

	server := grpc.NewServer(grpc.UnaryInterceptor(interceptors.AuthInterceptor))
	UMSpb.RegisterUMSServer(server, UMSHandler)

	slog.Info(logPrefix+"starting UMS gRPC server", slog.String("address", cfg.GRPCAddress))
	if err := server.Serve(listen); err != nil {
		slog.Error(logPrefix+"failed to serve", slog.Any("error", err))
	}
}
