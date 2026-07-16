package main

import (
	authpb "cloud/gen/auth/v1"
	filespb "cloud/gen/files/v1"
	userspb "cloud/gen/users/v1"
	"cloud/internal/config"
	auth "cloud/internal/handlers/auth"
	files "cloud/internal/handlers/files"
	users "cloud/internal/handlers/users"
	"cloud/internal/interceptors"
	"cloud/internal/storage"
	"log/slog"
	"net"

	"google.golang.org/grpc"
)

func main() {
	config, err := config.GetConfig()
	if err != nil {
		slog.Error("[main]: Failed to get config", slog.Any("error", err))
		return
	}

	conn, err := storage.Connect(config)
	if err != nil {
		slog.Error("[main]: Failed to connect to database", slog.Any("error", err))
		return
	}
	defer conn.Db.Close()

	listen, err := net.Listen("tcp", config.GRPCServerAddress)
	if err != nil {
		slog.Error("[main]: Failed to listen", slog.Any("error", err))
		return
	}

	gRPCServer := grpc.NewServer(grpc.UnaryInterceptor(interceptors.AuthInterceptor))

	UserService := storage.NewUserService(conn.Db)
	FileService := storage.NewFileService(conn.Db)

	authHandler := auth.NewAuthHandler(UserService)
	authpb.RegisterAuthServer(gRPCServer, authHandler)

	usersHandler := users.NewUsersHandler(UserService)
	userspb.RegisterUsersServer(gRPCServer, usersHandler)

	filesHandler := files.NewFileHandler(FileService)
	filespb.RegisterFilesServer(gRPCServer, filesHandler)

	slog.Info("[main]: Starting gRPC server", slog.String("address", config.GRPCServerAddress))
	if err := gRPCServer.Serve(listen); err != nil {
		slog.Error("[main]: Failed to serve", slog.Any("error", err))
	}
}
