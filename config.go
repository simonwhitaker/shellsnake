package main

import (
	"encoding/json"
	"os"
	"path/filepath"
)

type Config struct {
	HighScore int
}

func configFilePath() (string, error) {
	userConfigDir, err := os.UserConfigDir()
	if err != nil {
		return "", err
	}
	configFile := filepath.Join(userConfigDir, "bubblesnake", "config.json")
	return configFile, nil
}

func (c *Config) Save() error {
	data, err := json.Marshal(c)
	if err != nil {
		return err
	}

	configFile, err := configFilePath()
	if err != nil {
		return err
	}
	err = os.Mkdir(filepath.Dir(configFile), 0755)
	if err != nil && !os.IsExist(err) {
		return err
	}

	return os.WriteFile(configFile, data, 0644)
}

func LoadConfig() (Config, error) {
	configFile, err := configFilePath()
	if err != nil {
		return Config{}, err
	}

	data, err := os.ReadFile(configFile)
	if err != nil {
		return Config{}, err
	}

	var config Config
	err = json.Unmarshal(data, &config)
	if err != nil {
		return Config{}, err
	}
	return config, nil
}
