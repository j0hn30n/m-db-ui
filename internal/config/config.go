package config

import (
	"os"
)

type Config struct {
	Host     string
	Port     string
	MongoURI string
}

func Load() *Config {
	port := os.Getenv("PORT")
	if port == "" {
		port = "8082"
	}

	host := os.Getenv("HOST")
	if host == "" {
		host = "127.0.0.1"
	}

	mongoURI := os.Getenv("MONGO_URI")
	if mongoURI == "" {
		mongoURI = "mongodb://localhost:27017"
	}

	return &Config{
		Host:     host,
		Port:     port,
		MongoURI: mongoURI,
	}
}