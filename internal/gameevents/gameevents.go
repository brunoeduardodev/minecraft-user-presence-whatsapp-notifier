package gameevents

import (
	"context"
	"fmt"
	"strings"

	"github.com/brunoeduardodev/minecraft-user-presence-whatsapp-notifier/internal/logger"
)

type GameEvent struct {
	Username string
	Action   string
}

func getUsernameFromJoinEvent(line string) string {
	parts := strings.Split(line, "[Server thread/INFO]: ")
	if len(parts) != 2 {
		logger.Warn(context.Background(), "failed to parse username from join event")
		return ""
	}

	message := parts[1]

	username := strings.Split(message, " ")[0]
	return username
}

func ParseGameEvent(line string) (*GameEvent, error) {
	if strings.Contains(line, "joined the game") {
		return &GameEvent{
			Action:   "joined",
			Username: getUsernameFromJoinEvent(line),
		}, nil
	}

	if strings.Contains(line, "left the game") {
		return &GameEvent{
			Action:   "left",
			Username: getUsernameFromJoinEvent(line),
		}, nil
	}

	return nil, fmt.Errorf("unknown game event")
}
