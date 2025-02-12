package config

import (
	"fmt"
	"log"
	"os"
)

type Config struct {
	PostgreString string
	JwtKey        []byte
}

func NewConfig() *Config {
	dbPort, exists := os.LookupEnv("DB_PORT")
	if !exists {
		log.Fatal("DB_PORT not found in .env file")
	}
	dbUser, exists := os.LookupEnv("DB_USER")
	if !exists {
		log.Fatal("DB_USER not found in .env file")
	}
	dbPassword, exists := os.LookupEnv("DB_PASSWORD")
	if !exists {
		log.Fatal("DB_PASSWORD not found in .env file")
	}
	dbName, exists := os.LookupEnv("DB_NAME")
	if !exists {
		log.Fatal("DB_NAME not found in .env file")
	}
	dbHost, exists := os.LookupEnv("DB_HOST")
	if !exists {
		log.Fatal("DB_HOST not found in .env file")
	}

	jwtKey, exists := os.LookupEnv("JWT_KEY")
	if !exists {
		log.Fatal("JWT_KEY not found in .env file")
	}
	return &Config{
		PostgreString: fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=disable", dbHost, dbUser, dbPassword, dbName, dbPort),
		JwtKey:        []byte(jwtKey),
	}
}
