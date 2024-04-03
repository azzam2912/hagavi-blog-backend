package config

import (
	"log"
	"os"
	"path/filepath"
	"github.com/joho/godotenv"
)

func Config(key string) string {
	dir, err := filepath.Abs(filepath.Dir(os.Args[0]))
        if err != nil {
            log.Fatal(err)
        }
        environmentPath := filepath.Join(dir, ".env")
        err = godotenv.Load(environmentPath)
	if err != nil {
		log.Println("error loading .env trying to load .env.local")
		err = godotenv.Load(".env.local")
	}
	if err != nil {
		log.Fatal("error loading all env")
	}
	return os.Getenv(key)
}
