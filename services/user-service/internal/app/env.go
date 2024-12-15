package app

import (
	"fmt"
	"github.com/joho/godotenv"
	"os"
)

func LoadEnv() error {
	err := godotenv.Load("/app/.env")
	if err != nil {
		return err
	}
	return nil
}

func GetEnvVariable(key string) string {
	value := os.Getenv(key)
	if value == "" {
		fmt.Printf("Environment variable %s is not set\n", key)
	}
	return value
}
