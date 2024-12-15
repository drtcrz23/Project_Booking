package app

import (
	"fmt"
	"os"

	"github.com/joho/godotenv"
)

func LoadEnv() error {
	err := godotenv.Load("../.env")
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
