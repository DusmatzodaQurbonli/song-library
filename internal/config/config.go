package config

import (
	"encoding/json"
	"fmt"
	"os"
)

type DB struct {
	Host string `json:"host"`
	Port string `json:"port"`
	User string `json:"user"`
	Pass string `json:"pass"`
	Name string `json:"name"`
}

type LogLevel struct {
	INFO  string `json:"info"`
	DEBUG string `json:"debug"`
	WARN  string `json:"warn"`
	ERROR string `json:"error"`
}

type Server struct {
	Host string `json:"host"`
	Port string `json:"port"`
}

type Config struct {
	DB           DB       `json:"db"`
	LogLevel     LogLevel `json:"log_level"`
	Server       Server   `json:"server"`
	MusicInfoAPI string   `json:"music_info_api"`
}

func New() (*Config, error) {
	filePath := "configs/config.json"
	file, err := os.Open(filePath)
	if err != nil {
		return nil, fmt.Errorf("Failed to open configuration file: %w", err)
	}
	defer file.Close()

	var cfg Config

	decoder := json.NewDecoder(file)
	if err := decoder.Decode(&cfg); err != nil {
		return nil, fmt.Errorf("JSON decoding error: %w", err)
	}

	return &cfg, nil
}
