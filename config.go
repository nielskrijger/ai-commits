package main

import (
	"errors"
	"os"

	"gopkg.in/yaml.v3"
)

type FileConfigReader struct {
	Filename string
}

type Config struct {
	APIKey string `yaml:"api_key"`
}

func (r *FileConfigReader) openConfigFile() (*os.File, error) {
	file, err := os.Open(r.Filename)
	if err != nil {
		return nil, err
	}
	return file, nil
}

func (r *FileConfigReader) decodeConfigFile(file *os.File) (Config, error) {
	var cfg Config
	decoder := yaml.NewDecoder(file)
	if err := decoder.Decode(&cfg); err != nil {
		return cfg, err
	}
	return cfg, nil
}

func (r *FileConfigReader) getAPIKey(cfg Config) (string, error) {
	if cfg.APIKey == "" {
		return "", os.ErrNotExist
	}

	return cfg.APIKey, nil
}

func (r *FileConfigReader) ReadAPIKey() (string, error) {
	file, err := r.openConfigFile()
	if err != nil {
		apiKey, ok := os.LookupEnv("API_KEY")

		if !ok {
			return "", errors.New("environment variable API_KEY is not set")
		}
		return apiKey, nil
	}

	defer file.Close()

	cfg, err := r.decodeConfigFile(file)
	if err != nil {
		return "", err
	}

	return r.getAPIKey(cfg)
}
