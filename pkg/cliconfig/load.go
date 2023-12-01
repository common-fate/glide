package cliconfig

import (
	"os"
	"path/filepath"

	"github.com/BurntSushi/toml"
	"github.com/common-fate/clio"
)

// CurrentContext is a shorthand function which
// calls Load() to load the configuration,
// and then calls cfg.Current()
// to get the current API context.
// It returns an error if either method fails.
func CurrentContext() (*Context, error) {
	cfg, err := Load()
	if err != nil {
		return nil, err
	}

	return cfg.Current()
}

func Load() (*Config, error) {
	// if COMMONFATE_GDEPLOY_CONFIG_FILE is set, use a custom file path
	// for the config file location.
	// the file specified must exist.
	customPath := os.Getenv("COMMONFATE_GDEPLOY_CONFIG_FILE")
	if customPath != "" {
		return openConfigFile(customPath)
	}

	home, err := os.UserHomeDir()
	if err != nil {
		return nil, err
	}

	fp := filepath.Join(home, ".commonfate", "gdeploy")
	cfg, err := openConfigFile(fp)
	if os.IsNotExist(err) {
		// return an empty config if the file doesn't exist
		return Default(), nil
	}
	if err != nil {
		// otherwise if we get an error, return it
		return nil, err
	}

	return cfg, nil
}

func openConfigFile(filepath string) (*Config, error) {
	clio.Debugw("loading config", "path", filepath)
	file, err := os.Open(filepath)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var cfg Config

	_, err = toml.NewDecoder(file).Decode(&cfg)
	if err != nil {
		return nil, err
	}

	clio.Debugw("loaded config", "cfg", cfg)
	return &cfg, nil
}
