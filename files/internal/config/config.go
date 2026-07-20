package config

import (
	"errors"
	"fmt"
	"log/slog"
	"os"
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
	functionName = "Load"
	logPrefix    = "[" + functionName + "]: "
	nullString   = ""
)

type Config struct {
	GRPCAddress        string
	AccountGRPCAddress string
	DSN                string
	UserStoragePath    string
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

func envRoot() string {
	dir := moduleRoot()
	for {
		if _, err := os.Stat(filepath.Join(dir, ".env")); err == nil {
			return dir
		}
		parent := filepath.Dir(dir)
		if parent == dir {
			return moduleRoot()
		}
		dir = parent
	}
}

func Load() (*Config, error) {
	env, err := godotenv.Read(filepath.Join(envRoot(), ".env"))
	if err != nil {
		slog.Error(logPrefix+"failed to read .env file", slog.Any("error", err))
	}

	grpcConf := &gRPCConfig{
		Host: env["FILES_GRPC_HOST"],
		Port: env["FILES_GRPC_PORT"],
	}
	if grpcConf.Host == nullString || grpcConf.Port == nullString {
		slog.Error(logPrefix+"missing gRPC environment variables", slog.Any("error", ErrMissingEnvVariable))
		return nil, ErrMissingEnvVariable
	}

	accountGRPCConf := &gRPCConfig{
		Host: env["ACCOUNT_GRPC_HOST"],
		Port: env["ACCOUNT_GRPC_PORT"],
	}
	if accountGRPCConf.Host == nullString || accountGRPCConf.Port == nullString {
		slog.Error(logPrefix, slog.Any("error", ErrMissingEnvVariable))
		return nil, ErrMissingEnvVariable
	}

	dbConf := &DbConfig{
		Host:     env["FILES_DB_HOST"],
		Port:     env["FILES_DB_PORT"],
		DbName:   env["FILES_DB_NAME"],
		User:     env["FILES_DB_USER"],
		Password: env["FILES_DB_PASSWORD"],
		SSLMode:  env["FILES_DB_SSLMODE"],
	}
	if dbConf.Host == nullString || dbConf.Port == nullString || dbConf.DbName == nullString || dbConf.User == nullString || dbConf.Password == nullString || dbConf.SSLMode == nullString {
		slog.Error(logPrefix+"missing database environment variables", slog.Any("error", ErrMissingEnvVariable))
		return nil, ErrMissingEnvVariable
	}

	userStoragePath := env["USER_STORAGE_PATH"]
	if userStoragePath == nullString {
		slog.Error(logPrefix+"missing user storage path environment variable", slog.Any("error", ErrMissingEnvVariable))
		return nil, ErrMissingEnvVariable
	}

	cfg = &Config{
		GRPCAddress:        grpcConf.Address(),
		AccountGRPCAddress: accountGRPCConf.Address(),
		DSN:                dbConf.DSN(),
		UserStoragePath:    userStoragePath,
	}

	return cfg, nil
}

func GetConfig() (*Config, error) {
	if cfg == nil {
		return Load()
	}
	return cfg, nil
}
