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
	GRPCAddress string
	DSN         string
	JWTSecret   string
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
		Host: env["ACCOUNT_GRPC_HOST"],
		Port: env["ACCOUNT_GRPC_PORT"],
	}
	if grpcConf.Host == nullString || grpcConf.Port == nullString {
		slog.Error(logPrefix+"missing gRPC environment variables", slog.Any("error", ErrMissingEnvVariable))
		return nil, ErrMissingEnvVariable
	}

	dbConf := &DbConfig{
		Host:     env["ACCOUNT_DB_HOST"],
		Port:     env["ACCOUNT_DB_PORT"],
		DbName:   env["ACCOUNT_DB_NAME"],
		User:     env["ACCOUNT_DB_USER"],
		Password: env["ACCOUNT_DB_PASSWORD"],
		SSLMode:  env["ACCOUNT_DB_SSLMODE"],
	}
	if dbConf.Host == nullString || dbConf.Port == nullString || dbConf.DbName == nullString || dbConf.User == nullString || dbConf.Password == nullString || dbConf.SSLMode == nullString {
		slog.Error(logPrefix+"missing database environment variables", slog.Any("error", ErrMissingEnvVariable))
		return nil, ErrMissingEnvVariable
	}

	jwtSecret := env["JWT_SECRET"]
	if jwtSecret == nullString {
		slog.Error(logPrefix+"missing JWT environment variable", slog.Any("error", ErrMissingEnvVariable))
		return nil, ErrMissingEnvVariable
	}

	cfg = &Config{
		GRPCAddress: grpcConf.Address(),
		DSN:         dbConf.DSN(),
		JWTSecret:   jwtSecret,
	}

	return cfg, nil
}

func GetConfig() (*Config, error) {
	if cfg == nil {
		return Load()
	}
	return cfg, nil
}
