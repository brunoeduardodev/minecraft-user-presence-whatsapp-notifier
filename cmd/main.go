package main

import (
	"context"

	"github.com/brunoeduardodev/minecraft-user-presence-whatsapp-notifier/internal/config"
	"github.com/brunoeduardodev/minecraft-user-presence-whatsapp-notifier/internal/gameevents"
	"github.com/brunoeduardodev/minecraft-user-presence-whatsapp-notifier/internal/logger"
	"github.com/brunoeduardodev/minecraft-user-presence-whatsapp-notifier/internal/sftputil"
)

func main() {
	logger.Init()

	ctx := context.Background()
	cfg, err := config.Load()
	if err != nil {
		logger.Error(ctx, err, "failed to load config")
		return
	}

	sc, disconnect, err := sftputil.GetConnection(ctx, cfg.SftpUrl)
	if err != nil {
		logger.Error(ctx, err, "failed to get connection")
		return
	}

	logger.Info(ctx, "connected to ftp server")

	defer func() {
		logger.Info(ctx, "closing sftp client")
		err := disconnect()
		if err != nil {
			logger.Error(ctx, err, "failed to close sftp client")
		}
	}()

	logger.Info(ctx, "opening file")
	sftputil.TailFile(ctx, &sftputil.TailFileOptions{
		Path: "/logs/latest.log",
		Callback: func(line string) {
			event, err := gameevents.ParseGameEvent(line)
			if err != nil {
				return
			}

			if event.Action == "joined" {
				logger.Info(ctx, "joined", "username", event.Username)
			}
			if event.Action == "left" {
				logger.Info(ctx, "left", "username", event.Username)
			}

		},
		IgnoreExisting: true,
		Client:         sc,
	})

	// for {
	// 	hasEntry := scanner.Scan()
	// 	if !hasEntry {
	// 		logger.Info(ctx, "no more entries, sleeping")
	// 		time.Sleep(time.Second)
	// 		continue
	// 	}

	// 	line := scanner.Text()
	// 	logger.Info(ctx, "log line", "line", line)
	// }

	// logs, err := io.ReadAll(file)
	// if err != nil {
	// 	logger.Error(ctx, err, "failed to read file")
	// 	return
	// }

	// logLines := strings.Split(string(logs), "\n")
	// for _, logLine := range logLines {
	// 	logger.Info(ctx, "log line", "line", logLine)
	// }
}
