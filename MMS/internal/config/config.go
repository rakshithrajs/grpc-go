package config

import (
	"errors"
	"fmt"
	"log/slog"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/joho/godotenv"
)

var (
	ErrMissingEnvVariable = errors.New("missing environment variable")
)

type gRPCConfig struct {
	Host string
	Port string
}

func (g *gRPCConfig) Address() string {
	return g.Host + ":" + g.Port
}

type DbConfig struct {
	Host     string
	Port     string
	DbName   string
	User     string
	Password string
	SSLMode  string
}

func (d *DbConfig) DSN() string {
	return fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=%s", d.Host, d.Port, d.User, d.Password, d.DbName, d.SSLMode)
}

const (
	functionName      = "Load"
	logPrefix         = "[" + functionName + "]: "
	NullString        = ""
	UserIDMetadataKey = "userID"
)

type Config struct {
	GRPCAddress     string
	DSN             string
	UserStoragePath string
}

var cfg *Config

func moduleRoot() string {
	out, err := exec.Command("go", "env", "GOMOD").Output()
	if err != nil {
		panic(err)
	}
	root := strings.TrimSpace(string(out))
	return root[:len(root)-len("go.mod")]
}

func Load() (*Config, error) {
	env, err := godotenv.Read(filepath.Join(moduleRoot(), "..", ".env"))
	if err != nil {
		slog.Error(logPrefix+"failed to read .env file", slog.Any("error", err))
	}

	grpcConf := &gRPCConfig{
		Host: env["MMS_GRPC_HOST"],
		Port: env["MMS_GRPC_PORT"],
	}
	if grpcConf.Host == NullString || grpcConf.Port == NullString {
		slog.Error(logPrefix+"missing gRPC environment variables", slog.Any("error", ErrMissingEnvVariable))
		return nil, ErrMissingEnvVariable
	}

	dbConf := &DbConfig{
		Host:     env["MMS_DB_HOST"],
		Port:     env["MMS_DB_PORT"],
		DbName:   env["MMS_DB_NAME"],
		User:     env["MMS_DB_USER"],
		Password: env["MMS_DB_PASSWORD"],
		SSLMode:  env["MMS_DB_SSLMODE"],
	}
	if dbConf.Host == NullString || dbConf.Port == NullString || dbConf.DbName == NullString || dbConf.User == NullString || dbConf.Password == NullString || dbConf.SSLMode == NullString {
		slog.Error(logPrefix+"missing database environment variables", slog.Any("error", ErrMissingEnvVariable))
		return nil, ErrMissingEnvVariable
	}

	userStoragePath := env["USER_STORAGE_PATH"]
	if userStoragePath == NullString {
		slog.Error(logPrefix+"missing user storage path environment variable", slog.Any("error", ErrMissingEnvVariable))
		return nil, ErrMissingEnvVariable
	}

	cfg = &Config{
		GRPCAddress:     grpcConf.Address(),
		DSN:             dbConf.DSN(),
		UserStoragePath: userStoragePath,
	}

	return cfg, nil
}

func GetConfig() (*Config, error) {
	if cfg == nil {
		return Load()
	}
	return cfg, nil
}

func SetConfig(c *Config) {
	cfg = c
}
