package config

import (
	"errors"
	"fmt"
	"log/slog"
	"os/exec"
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

func (g *gRPCConfig) BuildgRPCAddress() string {
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

func (d *DbConfig) GetDSN() string {
	return fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=%s", d.Host, d.Port, d.User, d.Password, d.DbName, d.SSLMode)
}

type Config struct {
	GRPCServerAddress string
	DSN               string
	JWTSecret         string
	UserStoragePath   string
}

func ProjectRoot() string {
	out, err := exec.Command("go", "env", "GOMOD").Output()
	if err != nil {
		panic(err)
	}
	root := strings.TrimSpace(string(out))
	return root[:len(root)-len("go.mod")]
}

var cfg *Config

func Load() (*Config, error) {
	configPath := ProjectRoot() + ".env"

	env, err := godotenv.Read(configPath)
	if err != nil {
		slog.Error("[Load]: Error loading .env file:", slog.Any("error", err))
	}

	GrpcConf := &gRPCConfig{
		Host: env["GRPC_HOST"],
		Port: env["GRPC_PORT"],
	}

	if GrpcConf.Host == "" || GrpcConf.Port == "" {
		slog.Error("[Load]: Missing environment variable:", slog.Any("error", ErrMissingEnvVariable))
		return nil, ErrMissingEnvVariable
	}

	DbConf := &DbConfig{
		Host:     env["DB_HOST"],
		Port:     env["DB_PORT"],
		DbName:   env["DB_NAME"],
		User:     env["DB_USER"],
		Password: env["DB_PASSWORD"],
		SSLMode:  env["DB_SSLMODE"],
	}

	if DbConf.Host == "" || DbConf.Port == "" || DbConf.DbName == "" || DbConf.User == "" || DbConf.Password == "" || DbConf.SSLMode == "" {
		slog.Error("[Load]: Missing environment variable:", slog.Any("error", ErrMissingEnvVariable))
		return nil, ErrMissingEnvVariable
	}

	JWTSecret := env["JWT_SECRET"]
	if JWTSecret == "" {
		slog.Error("[Load]: Missing environment variable:", slog.Any("error", ErrMissingEnvVariable))
		return nil, ErrMissingEnvVariable
	}

	UserStoragePath := env["USER_STORAGE_PATH"]
	if UserStoragePath == "" {
		slog.Error("[Load]: Missing environment variable:", slog.Any("error", ErrMissingEnvVariable))
		return nil, ErrMissingEnvVariable
	}

	cfg = &Config{
		GRPCServerAddress: GrpcConf.BuildgRPCAddress(),
		DSN:               DbConf.GetDSN(),
		JWTSecret:         JWTSecret,
		UserStoragePath:   UserStoragePath,
	}

	return cfg, nil
}

func GetConfig() (*Config, error) {
	if cfg == nil {
		_, err := Load()
		if err != nil {
			return nil, err
		}
		return cfg, nil
	}
	return cfg, nil
}
