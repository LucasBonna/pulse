package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

type Env struct {
	Port  string
	Token string
}

func InitEnvs() *Env {
	_ = godotenv.Load()

	port := os.Getenv("PORT")
	token := os.Getenv("TOKEN")

	if port == "" {
		port = ":8080" // default
	}
	if token == "" {
		log.Fatal("TOKEN environment variable is required")
	}

	return &Env{
		Port:  port,
		Token: token,
	}
}
