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
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	port := os.Getenv("PORT")
	token := os.Getenv("TOKEN")
	return &Env{
		Port:  port,
		Token: token,
	}
}
