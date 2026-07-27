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

// Address returns the full address for the gRPC server in the format "host:port".
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

// DSN returns the Data Source Name for connecting to the database, formatted as a connection string.
func (d *DbConfig) DSN() string {
	return fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=%s", d.Host, d.Port, d.User, d.Password, d.DbName, d.SSLMode)
}

const (
	// functionName is Load
	functionName = "Load"

	// logPrefix is the prefix for log messages in this package
	logPrefix = "[" + functionName + "]: "

	// NullString is the representation of a ""
	NullString = ""

	// UserIDMetadataKey is the key to be used in metadata for userID
	UserIDMetadataKey = "userID"

	// ErrorKey is the key for error messages in logs
	ErrorKey = "error"
)

type Config struct {
	GRPCAddress     string
	DSN             string
	UserStoragePath string
}

var cfg *Config

// moduleRoot returns the root directory of the Go module by executing "go env GOMOD" and trimming the output.
func moduleRoot() string {
	out, err := exec.Command("go", "env", "GOMOD").Output()
	if err != nil {
		panic(err)
	}
	root := strings.TrimSpace(string(out))
	return root[:len(root)-len("go.mod")]
}

// Load reads the configuration from the .env file and returns a Config struct. It also validates that all required environment variables are present.
func Load() (*Config, error) {
	env, err := godotenv.Read(filepath.Join(moduleRoot(), "..", ".env"))
	if err != nil {
		slog.Error(logPrefix+"failed to read .env file", slog.Any(ErrorKey, err))
	}

	grpcConf := &gRPCConfig{
		Host: env["MMS_GRPC_HOST"],
		Port: env["MMS_GRPC_PORT"],
	}
	if grpcConf.Host == NullString || grpcConf.Port == NullString {
		slog.Error(logPrefix+"missing gRPC environment variables", slog.Any(ErrorKey, ErrMissingEnvVariable))
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
		slog.Error(logPrefix+"missing database environment variables", slog.Any(ErrorKey, ErrMissingEnvVariable))
		return nil, ErrMissingEnvVariable
	}

	userStoragePath := env["USER_STORAGE_PATH"]
	if userStoragePath == NullString {
		slog.Error(logPrefix+"missing user storage path environment variable", slog.Any(ErrorKey, ErrMissingEnvVariable))
		return nil, ErrMissingEnvVariable
	}

	cfg = &Config{
		GRPCAddress:     grpcConf.Address(),
		DSN:             dbConf.DSN(),
		UserStoragePath: userStoragePath,
	}

	return cfg, nil
}

// GetConfig returns the current configuration. If the configuration has not been loaded yet, it calls Load to load it.
func GetConfig() (*Config, error) {
	if cfg == nil {
		return Load()
	}
	return cfg, nil
}

// SetConfig sets the current configuration. This is useful for testing purposes.
func SetConfig(c *Config) {
	cfg = c
}
