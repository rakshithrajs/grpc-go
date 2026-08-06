package config

import (
	"errors"
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

type gRPCConfig struct {
	Host string
	Port string
}

// Address returns the full address for the gRPC server in the format "host:port".
func (g *gRPCConfig) Address() string {
	return g.Host + ":" + g.Port
}

const (
	// functionName is Load
	functionName = "Load"

	// logPrefix is the prefix for log messages in this package
	logPrefix = "[" + functionName + "]: "

	// NullString is the representation of an empty string
	NullString = ""

	// ErrorKey is the key for error messages in logs
	ErrorKey = "error"
)

// Config holds all application-wide configuration values loaded from environment variables.
type Config struct {
	GRPCAddress string
	JWTSecret   string
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
		Host: env["TMS_GRPC_HOST"],
		Port: env["TMS_GRPC_PORT"],
	}
	if grpcConf.Host == NullString || grpcConf.Port == NullString {
		slog.Error(logPrefix+"missing TMS gRPC environment variables", slog.Any(ErrorKey, ErrMissingEnvVariable))
		return nil, ErrMissingEnvVariable
	}

	jwtSecret := env["JWT_SECRET"]
	if jwtSecret == NullString {
		slog.Error(logPrefix+"missing JWT environment variable", slog.Any(ErrorKey, ErrMissingEnvVariable))
		return nil, ErrMissingEnvVariable
	}

	cfg = &Config{
		GRPCAddress: grpcConf.Address(),
		JWTSecret:   jwtSecret,
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
