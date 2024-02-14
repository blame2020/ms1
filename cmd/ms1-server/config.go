package main

import (
	"io"
	"log/slog"

	"github.com/kelseyhightower/envconfig"
)

// Config is configurations for this program.
// Configurations are mapped from envronment variables based on "github.com/kelseyhightower/envconfig".
type Config struct {
	GrpcPort int        `default:"50051" desc:"A port number for the gRPC server"`
	LogLevel slog.Level `default:"info"`
}

// NewConfigFromEnv loads configurations from environment, returns it.
func NewConfigFromEnv() (*Config, error) {
	c := &Config{}

	if err := envconfig.Process("test", c); err != nil {
		return nil, err
	}

	return c, nil
}

// Usage prints usage messages to `out`.
func (c *Config) Usage(out io.Writer) error {
	return envconfig.Usagef("test", c, out, envconfig.DefaultTableFormat)
}
