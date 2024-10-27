package config

import (
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	SftpUrl string
}

func Load() (*Config, error) {
	err := godotenv.Load()
	if err != nil {
		return nil, err
	}
	sftpUrl := os.Getenv("SFTP_URL")

	return &Config{
		SftpUrl: sftpUrl,
	}, nil
}
