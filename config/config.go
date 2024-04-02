package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

func Config(key string) string {
	err := godotenv.Load(".env")
	if err != nil {
		err = godotenv.Load(".env.local")
	}
	if err != nil {
		log.Printf("error loading env")
	}
	return os.Getenv(key)
}
