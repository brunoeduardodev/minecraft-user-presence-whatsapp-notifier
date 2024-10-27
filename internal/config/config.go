package config

import (
	"context"
	"os"
	"strconv"

	"github.com/brunoeduardodev/minecraft-user-presence-whatsapp-notifier/internal/logger"
	"github.com/joho/godotenv"
)

type Config struct {
	SftpUrl         string
	CheckKnownHosts bool
	GroupName       string
	JoinMessage     string
	LeaveMessage    string
}

func Load() (*Config, error) {
	err := godotenv.Load()
	if err != nil {
		logger.Warn(context.Background(), ".env file not found...")
	}
	sftpUrl := os.Getenv("SFTP_URL")

	checkKnowHosts, err := strconv.ParseBool(os.Getenv("CHECK_KNOWN_HOSTS"))
	if err != nil {
		logger.Warn(context.Background(), "failed to parse CHECK_KNOWN_HOSTS")
		checkKnowHosts = false
	}

	return &Config{
		SftpUrl:         sftpUrl,
		JoinMessage:     os.Getenv("JOIN_MESSAGE"),
		LeaveMessage:    os.Getenv("LEAVE_MESSAGE"),
		GroupName:       os.Getenv("GROUP_NAME"),
		CheckKnownHosts: checkKnowHosts,
	}, nil
}
