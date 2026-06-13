package config

import (
	"log/slog"
	"time"
)

type Config struct {
	LogLevel        string        `mapstructure:"log_level"`
	Workers         int           `mapstructure:"workers"`
	Timeout         time.Duration `mapstructure:"timeout"`
	OutputFormat    string        `mapstructure:"output_format"`
	StringMinLength int           `mapstructure:"string_min_length"`
}

func Default() Config {
	return Config{
		LogLevel:        slog.LevelInfo.String(),
		Workers:         4,
		Timeout:         10 * time.Minute,
		OutputFormat:    "markdown",
		StringMinLength: 5,
	}
}
