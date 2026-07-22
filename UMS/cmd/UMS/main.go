package main

import (
	"log/slog"

	"github.com/gin-gonic/gin"

	MMSpb "github.com/rakshithrajs/cloud/UMS/gen/MMS/v1"
	"github.com/rakshithrajs/cloud/UMS/internal/config"
	user "github.com/rakshithrajs/cloud/UMS/internal/handlers/user"
	userfiles "github.com/rakshithrajs/cloud/UMS/internal/handlers/userFiles"
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
		slog.Error(logPrefix+"failed to get config", slog.Any("error", err))
		return
	}

	db, err := storage.Connect(cfg.DSN)
	if err != nil {
		slog.Error(logPrefix+"failed to connect to database", slog.Any("error", err))
		return
	}
	defer db.Close()

	MMSConn, err := grpc.NewClient(cfg.MMSGRPCAddress, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		slog.Error(logPrefix+"failed to connect to MMS service", slog.Any("error", err))
		return
	}
	defer MMSConn.Close()

	MMSClient := MMSpb.NewFilesClient(MMSConn)

	gin.SetMode(gin.DebugMode)
	router := gin.New()
	router.SetTrustedProxies(nil)

	router.Use(gin.Logger())

	store := storage.NewUserStore(db)
	UMSHandler := user.NewUMSHandler(store)

	userFilesStore := storage.NewUserFilesStore(db)
	UserFilesHandler := userfiles.NewUserFilesHandler(userFilesStore, MMSClient)

	UMSRouterGroup := router.Group(apiPrefix + "/user")
	user.RegisterRoutes(UMSRouterGroup, UMSHandler)

	filesRouterGroup := router.Group(apiPrefix + "/files")
	userfiles.RegisterRoutes(filesRouterGroup, UserFilesHandler)

	if err := router.Run(cfg.ServerAddress); err != nil {
		slog.Error(logPrefix+"failed to run server", slog.Any("error", err))
	}
}
