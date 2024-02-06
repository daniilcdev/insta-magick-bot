package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

type IMConfig interface {
	InDir() string
	OutDir() string
	TempDir() string
}

type imConfig struct {
	in, out, tmp string
}

func (c *imConfig) InDir() string {
	return c.in
}

func (c *imConfig) OutDir() string {
	return c.out
}

func (c *imConfig) TempDir() string {
	return c.out
}

func Load() IMConfig {
	err := godotenv.Load("./config/env/imagemagick.env")
	if err != nil {
		log.Printf("env loading failed: '%v'\n", err)
	}

	return &imConfig{
		in:  os.Getenv("IM_IN_DIR"),
		out: os.Getenv("IM_OUT_DIR"),
		tmp: os.Getenv("IM_TEMP_DIR"),
	}
}
