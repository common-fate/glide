package cliconfig

import (
	"os"
	"path/filepath"

	"github.com/BurntSushi/toml"
)

func Save(cfg *Config) error {
	// if COMMONFATE_CONFIG_FILE is set, use a custom file path
	// for the config file location.
	// the file specified must exist.
	customPath := os.Getenv("COMMONFATE_CONFIG_FILE")
	if customPath != "" {
		return saveConfigFile(cfg, customPath)
	}
	home, err := os.UserHomeDir()
	if err != nil {
		return err
	}

	configFolder := filepath.Join(home, ".commonfate")

	_, err = os.Stat(configFolder)
	if os.IsNotExist(err) {
		err := os.MkdirAll(configFolder, os.ModePerm)
		if err != nil {
			return err
		}
	}
	fp := filepath.Join(configFolder, "config")
	return saveConfigFile(cfg, fp)
}

func saveConfigFile(cfg *Config, path string) error {
	file, err := os.OpenFile(path, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0600)
	if err != nil {
		return err
	}
	defer file.Close()
	return toml.NewEncoder(file).Encode(cfg)
}
