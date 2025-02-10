package config

import (
	"log"
	"os"
)

type Config struct {
	PostgreString string
}

func NewConfig() *Config {
	postgreString, exists := os.LookupEnv("PostgreString")
	if !exists {
		log.Fatal("PostgreString not found in .env file")
	}
	return &Config{PostgreString: postgreString}
}
