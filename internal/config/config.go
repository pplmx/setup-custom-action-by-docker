// internal/config/config.go

package config

import (
	"fmt"

	"github.com/BurntSushi/toml"
)

type Config struct {
	InputText   string    `toml:"input_text"`
	FindWord    string    `toml:"find_word"`
	ReplaceWord string    `toml:"replace_word"`
	NumberList  []float64 `toml:"number_list"`
	InputFile   string    `toml:"input_file"`
	OutputFile  string    `toml:"output_file"`
	AppendText  string    `toml:"append_text"`
	ApiURL      string    `toml:"api_url"`
}

const defaultConfigPath = ".github/configs/setup-custom-action-by-docker.toml"

func LoadConfig(path string) (*Config, error) {
	if path == "" {
		path = defaultConfigPath
	}

	var config Config
	if _, err := toml.DecodeFile(path, &config); err != nil {
		return nil, fmt.Errorf("failed to decode TOML file: %w", err)
	}

	if err := validateConfig(&config); err != nil {
		return nil, fmt.Errorf("invalid configuration: %w", err)
	}

	return &config, nil
}

func validateConfig(cfg *Config) error {
	if cfg.InputText == "" {
		return fmt.Errorf("input_text is required")
	}
	if cfg.InputFile == "" || cfg.OutputFile == "" {
		return fmt.Errorf("input_file and output_file are required")
	}
	if cfg.ApiURL == "" {
		return fmt.Errorf("api_url is required")
	}
	return nil
}
