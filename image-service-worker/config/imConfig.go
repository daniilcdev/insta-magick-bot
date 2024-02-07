package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

type WorkerConfig struct {
	In, Out string
}

func Load() (*WorkerConfig, error) {
	err := godotenv.Load("./config/env/imagemagick.env")
	if err != nil {
		log.Printf("env loading failed: '%v'\n", err)
		return nil, err
	}

	cfg := WorkerConfig{
		In:  os.Getenv("IM_IN_DIR"),
		Out: os.Getenv("IM_OUT_DIR"),
	}

	if err := validatePath(cfg.In, cfg.Out); err != nil {
		return nil, err
	}

	return &cfg, nil
}

func validatePath(paths ...string) error {
	for _, path := range paths {
		if _, err := os.Stat(path); err != nil {
			return err
		}
	}

	return nil
}
