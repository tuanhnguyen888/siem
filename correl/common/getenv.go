package common

import (
	"errors"
	"github.com/joho/godotenv"
	"os"
)

func GoDotEnvVariable(key string) (string,error) {
	// load .env file
	err := godotenv.Load("D:\\doan_siem\\correl\\.env")
	if err != nil {
		return "",errors.New("error loading .env file")
	}
	str, ok := os.LookupEnv(key)
	if !ok {
		return "",errors.New("no find .env key")
	}
	return str, nil
}
