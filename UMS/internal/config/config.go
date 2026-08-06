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
	// ErrMissingEnvVariable is returned when a required environment variable is missing.
	ErrMissingEnvVariable = errors.New("missing environment variable")
)

const (
	// functionName is Load
	functionName = "Load"

	// logPrefix is the prefix for log messages in this package
	logPrefix = "[" + functionName + "]: "

	// NullString is the representation of a ""
	NullString = ""

	// ErrorKey is the key for error messages in logs
	ErrorKey = "error"

	// UserIDMetadataKey is the key for userID
	UserIDMetadataKey = "userID"
)

// ServerConfig contains the host and port used to build server addresses.
type ServerConfig struct {
	Host string
	Port string
}

// Address returns the host:port server address.
func (g *ServerConfig) Address() string {
	return g.Host + ":" + g.Port
}

// DbConfig contains the database connection parameters.
type DbConfig struct {
	Host     string
	Port     string
	DbName   string
	User     string
	Password string
	SSLMode  string
}

// DSN returns the PostgreSQL connection string.
func (d *DbConfig) DSN() string {
	return fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=%s", d.Host, d.Port, d.User, d.Password, d.DbName, d.SSLMode)
}

// Config holds all application-wide configuration values loaded from environment variables.
type Config struct {
	ServerAddress  string
	MMSGRPCAddress string
	TMSGRPCAddress string
	DSN            string
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

// Load loads the configuration from environment variables.
func Load() (*Config, error) {
	env, err := godotenv.Read(filepath.Join(moduleRoot(), "..", ".env"))
	if err != nil {
		slog.Error(logPrefix+"failed to read .env file", slog.Any(ErrorKey, err))
	}

	ServerConf := &ServerConfig{
		Host: env["UMS_HOST"],
		Port: env["UMS_PORT"],
	}
	if ServerConf.Host == NullString || ServerConf.Port == NullString {
		slog.Error(logPrefix+"missing UMS gRPC environment variables", slog.Any(ErrorKey, ErrMissingEnvVariable))
		return nil, ErrMissingEnvVariable
	}

	MMSGRPCConf := &ServerConfig{
		Host: env["MMS_GRPC_HOST"],
		Port: env["MMS_GRPC_PORT"],
	}
	if MMSGRPCConf.Host == NullString || MMSGRPCConf.Port == NullString {
		slog.Error(logPrefix+"missing MMS gRPC environment variables", slog.Any(ErrorKey, ErrMissingEnvVariable))
		return nil, ErrMissingEnvVariable
	}

	TMSGRPCConf := &ServerConfig{
		Host: env["TMS_GRPC_HOST"],
		Port: env["TMS_GRPC_PORT"],
	}
	if TMSGRPCConf.Host == NullString || TMSGRPCConf.Port == NullString {
		slog.Error(logPrefix+"missing TMS gRPC environment variables", slog.Any(ErrorKey, ErrMissingEnvVariable))
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
		slog.Error(logPrefix+"missing database environment variables", slog.Any(ErrorKey, ErrMissingEnvVariable))
		return nil, ErrMissingEnvVariable
	}

	cfg = &Config{
		ServerAddress:  ServerConf.Address(),
		MMSGRPCAddress: MMSGRPCConf.Address(),
		TMSGRPCAddress: TMSGRPCConf.Address(),
		DSN:            dbConf.DSN(),
	}

	return cfg, nil
}

// GetConfig returns the current configuration, loading it if necessary.
func GetConfig() (*Config, error) {
	if cfg == nil {
		return Load()
	}
	return cfg, nil
}

// SetConfig sets the current configuration.
func SetConfig(c *Config) {
	cfg = c
}
