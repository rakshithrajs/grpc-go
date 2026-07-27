package main

import (
	"log/slog"

	"github.com/gin-gonic/gin"

	MMSpb "github.com/rakshithrajs/cloud/UMS/gen/MMS/v1"
	"github.com/rakshithrajs/cloud/UMS/internal/config"
	mmsGrpc "github.com/rakshithrajs/cloud/UMS/internal/grpcClient"
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
	defer db.Close()

	MMSConn, err := grpc.NewClient(cfg.MMSGRPCAddress, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		slog.Error(logPrefix+"failed to connect to MMS service", slog.Any(config.ErrorKey, err))
		return
	}
	defer MMSConn.Close()

	MMSClient := MMSpb.NewFilesClient(MMSConn)

	gin.SetMode(gin.DebugMode)
	router := gin.New()
	router.SetTrustedProxies(nil)

	router.Use(gin.Logger())

	userStore := storage.NewUserStore(db)
	UserHandler := user.NewUserHandler(userStore)

	userFilesStore := storage.NewUserFilesStore(db)
	mmsClient := mmsGrpc.NewClient(MMSClient, userFilesStore)
	UserFilesHandler := userFiles.NewUserFilesHandler(mmsClient, userFilesStore)

	usersRouterGroup := router.Group(apiPrefix + "/users")
	user.RegisterRoutes(usersRouterGroup, UserHandler)
	userFiles.RegisterRoutes(usersRouterGroup, UserFilesHandler)

	if err := router.Run(cfg.ServerAddress); err != nil {
		slog.Error(logPrefix+"failed to run server", slog.Any(config.ErrorKey, err))
	}
}
