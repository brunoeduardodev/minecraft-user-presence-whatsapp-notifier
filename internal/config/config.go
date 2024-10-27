package config

import (
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	SftpUrl      string
	GroupName    string
	JoinMessage  string
	LeaveMessage string
}

func Load() (*Config, error) {
	err := godotenv.Load()
	if err != nil {
		return nil, err
	}
	sftpUrl := os.Getenv("SFTP_URL")

	return &Config{
		SftpUrl:      sftpUrl,
		JoinMessage:  os.Getenv("JOIN_MESSAGE"),
		LeaveMessage: os.Getenv("LEAVE_MESSAGE"),
		GroupName:    os.Getenv("GROUP_NAME"),
	}, nil
}
