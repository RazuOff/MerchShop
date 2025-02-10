package config

import (
	"log"
	"os"
)

type Config struct {
	PostgreString string
	JwtKey        []byte
}

func NewConfig() *Config {
	postgreString, exists := os.LookupEnv("POSTGRE_STRING")
	if !exists {
		log.Fatal("POSTGRE_STRING not found in .env file")
	}

	jwtKey, exists := os.LookupEnv("JWT_KEY")
	if !exists {
		log.Fatal("JWT_KEY not found in .env file")
	}
	return &Config{PostgreString: postgreString, JwtKey: []byte(jwtKey)}
}
