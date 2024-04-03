package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

func Config(key string) string {
	err := godotenv.Load(".env")
	if err != nil {
		log.Println("error loading .env trying to load .env.local")
		err = godotenv.Load(".env.local")
	}
	if err != nil {
		log.Fatal("error loading all env")
	}
	return os.Getenv(key)
}
