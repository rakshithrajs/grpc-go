package main

import (
	"log/slog"

	"github.com/gin-gonic/gin"

	MMSpb "github.com/rakshithrajs/cloud/UMS/gen/MMS/v1"
	TMSpb "github.com/rakshithrajs/cloud/UMS/gen/TMS/v1"
	"github.com/rakshithrajs/cloud/UMS/internal/config"
	grpcClient "github.com/rakshithrajs/cloud/UMS/internal/grpcClient"
	user "github.com/rakshithrajs/cloud/UMS/internal/handlers/user"
	userFiles "github.com/rakshithrajs/cloud/UMS/internal/handlers/userFiles"
	"github.com/rakshithrajs/cloud/UMS/internal/storage"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

const (
	functionName = "main"
	logPrefix    = "[" + functionName + "]: "
	apiPrefix    = "/api"
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
	defer func() {
		slog.Info(logPrefix + "Closing Database connection")
		if err := db.Close(); err != nil {
			slog.Error(logPrefix+"failed to close database connection", slog.Any(config.ErrorKey, err))
		}
	}()

	MMSConn, err := grpc.NewClient(cfg.MMSGRPCAddress, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		slog.Error(logPrefix+"failed to connect to MMS service", slog.Any(config.ErrorKey, err))
		return
	}
	defer func() {
		slog.Info(logPrefix + "Closing MMS gRPC connection")
		if err := MMSConn.Close(); err != nil {
			slog.Error(logPrefix+"failed to close MMS gRPC connection", slog.Any(config.ErrorKey, err))
		}
	}()

	TMSConn, err := grpc.NewClient(cfg.TMSGRPCAddress, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		slog.Error(logPrefix+"failed to connect to TMS service", slog.Any(config.ErrorKey, err))
		return
	}
	defer func() {
		slog.Info(logPrefix + "Closing TMS gRPC connection")
		if err := TMSConn.Close(); err != nil {
			slog.Error(logPrefix+"failed to close TMS gRPC connection", slog.Any(config.ErrorKey, err))
		}
	}()

	MMSClient := MMSpb.NewFilesClient(MMSConn)
	TMSClient := TMSpb.NewTokensClient(TMSConn)

	gin.SetMode(gin.DebugMode)
	router := gin.New()
	router.SetTrustedProxies(nil)

	router.Use(gin.Logger())

	userStore := storage.NewUserStore(db)
	tmsClient := grpcClient.NewTMSClient(TMSClient)
	UserHandler := user.NewUserHandler(userStore, tmsClient)

	userFilesStore := storage.NewUserFilesStore(db)
	mmsClient := grpcClient.NewMMSClient(MMSClient, userFilesStore)
	UserFilesHandler := userFiles.NewUserFilesHandler(mmsClient, userFilesStore)

	usersRouterGroup := router.Group(apiPrefix + "/users")
	user.RegisterRoutes(usersRouterGroup, UserHandler, tmsClient)
	userFiles.RegisterRoutes(usersRouterGroup, UserFilesHandler, tmsClient)

	if err := router.Run(cfg.ServerAddress); err != nil {
		slog.Error(logPrefix+"failed to run server", slog.Any(config.ErrorKey, err))
	}
}
