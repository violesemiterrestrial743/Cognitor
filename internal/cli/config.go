package cli

import (
	"github.com/kernelstub/cognitor/internal/config"
	"github.com/spf13/viper"
)

func loadConfig(path string) (config.Config, error) {
	cfg := config.Default()
	v := viper.New()
	v.SetConfigType("yaml")
	v.SetConfigName("default")
	v.AddConfigPath("configs")
	if path != "" {
		v.SetConfigFile(path)
	}
	v.SetDefault("log_level", cfg.LogLevel)
	v.SetDefault("workers", cfg.Workers)
	v.SetDefault("timeout", cfg.Timeout)
	v.SetDefault("output_format", cfg.OutputFormat)
	v.SetDefault("string_min_length", cfg.StringMinLength)
	if err := v.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); !ok || path != "" {
			return cfg, err
		}
	}
	if err := v.Unmarshal(&cfg); err != nil {
		return cfg, err
	}
	return cfg, config.Validate(cfg)
}
