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

const (
	functionName = "Load"
	logPrefix    = "[" + functionName + "]: "
	NullString   = ""
	ErrorKey     = "error"
)

type ServerConfig struct {
	Host string
	Port string
}

func (g *ServerConfig) Address() string {
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

type Config struct {
	ServerAddress  string
	MMSGRPCAddress string
	DSN            string
	JWTSecret      string
}

var cfg *Config

func moduleRoot() string {
	out, err := exec.Command("go", "env", "GOMOD").Output()
	if err != nil {
		panic(err)
	}
	root := strings.TrimSpace(string(out))
	fmt.Println("Module root:", root)
	return root[:len(root)-len("go.mod")]
}

func Load() (*Config, error) {
	env, err := godotenv.Read(filepath.Join(moduleRoot(), "..", ".env"))
	if err != nil {
		slog.Error(logPrefix+"failed to read .env file", slog.Any("error", err))
	}

	ServerConf := &ServerConfig{
		Host: env["UMS_HOST"],
		Port: env["UMS_PORT"],
	}
	if ServerConf.Host == NullString || ServerConf.Port == NullString {
		slog.Error(logPrefix+"missing UMS gRPC environment variables", slog.Any("error", ErrMissingEnvVariable))
		return nil, ErrMissingEnvVariable
	}

	MMSGRPCConf := &ServerConfig{
		Host: env["MMS_GRPC_HOST"],
		Port: env["MMS_GRPC_PORT"],
	}
	if MMSGRPCConf.Host == NullString || MMSGRPCConf.Port == NullString {
		slog.Error(logPrefix+"missing MMS gRPC environment variables", slog.Any("error", ErrMissingEnvVariable))
		return nil, ErrMissingEnvVariable
	}

	dbConf := &DbConfig{
		Host:     env["UMS_DB_HOST"],
		Port:     env["UMS_DB_PORT"],
		DbName:   env["UMS_DB_NAME"],
		User:     env["UMS_DB_USER"],
		Password: env["UMS_DB_PASSWORD"],
		SSLMode:  env["UMS_DB_SSLMODE"],
	}
	if dbConf.Host == NullString || dbConf.Port == NullString || dbConf.DbName == NullString || dbConf.User == NullString || dbConf.Password == NullString || dbConf.SSLMode == NullString {
		slog.Error(logPrefix+"missing database environment variables", slog.Any("error", ErrMissingEnvVariable))
		return nil, ErrMissingEnvVariable
	}

	jwtSecret := env["JWT_SECRET"]
	if jwtSecret == NullString {
		slog.Error(logPrefix+"missing JWT environment variable", slog.Any("error", ErrMissingEnvVariable))
		return nil, ErrMissingEnvVariable
	}

	cfg = &Config{
		ServerAddress:  ServerConf.Address(),
		MMSGRPCAddress: MMSGRPCConf.Address(),
		DSN:            dbConf.DSN(),
		JWTSecret:      jwtSecret,
	}

	return cfg, nil
}

func GetConfig() (*Config, error) {
	if cfg == nil {
		return Load()
	}
	return cfg, nil
}
