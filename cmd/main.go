package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/mdp/qrterminal/v3"

	"github.com/brunoeduardodev/minecraft-user-presence-whatsapp-notifier/internal/config"
	"github.com/brunoeduardodev/minecraft-user-presence-whatsapp-notifier/internal/gameevents"
	"github.com/brunoeduardodev/minecraft-user-presence-whatsapp-notifier/internal/logger"
	"github.com/brunoeduardodev/minecraft-user-presence-whatsapp-notifier/internal/sftputil"
	_ "github.com/mattn/go-sqlite3"

	"go.mau.fi/whatsmeow"
	"go.mau.fi/whatsmeow/proto/waE2E"
	"go.mau.fi/whatsmeow/store/sqlstore"
	"go.mau.fi/whatsmeow/types"
	waLog "go.mau.fi/whatsmeow/util/log"
)

// func eventHandler(evt interface{}) {
// 	switch v := evt.(type) {
// 	case *events.Message:
// 		logger.Info(context.Background(), "Received a message!", "message", v.Message.GetConversation())
// 	}
// }

func main() {
	logger.Init()

	ctx := context.Background()
	cfg, err := config.Load()
	if err != nil {
		logger.Error(ctx, err, "failed to load config")
		return
	}

	botPrefix := "🤖 "

	dbLog := waLog.Stdout("Database", "DEBUG", true)
	container, err := sqlstore.New("sqlite3", "file:wastore.db?_foreign_keys=on", dbLog)
	if err != nil {
		panic(err)
	}
	// If you want multiple sessions, remember their JIDs and use .GetDevice(jid) or .GetAllDevices() instead.
	deviceStore, err := container.GetFirstDevice()
	if err != nil {
		panic(err)
	}

	clientLog := waLog.Stdout("Client", "INFO", true)
	logger.Info(ctx, "connecting to whatsapp")
	waClient := whatsmeow.NewClient(deviceStore, clientLog)
	// waClient.AddEventHandler(eventHandler)
	if waClient.Store.ID == nil {
		// No ID stored, new login
		qrChan, _ := waClient.GetQRChannel(context.Background())
		err = waClient.Connect()
		if err != nil {
			panic(err)
		}
		for evt := range qrChan {
			if evt.Event == "code" {
				// Render the QR code here
				// e.g. qrterminal.GenerateHalfBlock(evt.Code, qrterminal.L, os.Stdout)
				// or just manually `echo 2@... | qrencode -t ansiutf8` in a terminal

				qrterminal.GenerateHalfBlock(evt.Code, qrterminal.L, os.Stdout)

			} else {
				fmt.Println("Login event:", evt.Event)
			}
		}
	} else {
		// Already logged in, just connect
		err = waClient.Connect()
		if err != nil {
			panic(err)
		}
	}

	logger.Info(ctx, "connected to whatsapp")

	groups, err := waClient.GetJoinedGroups()
	if err != nil {
		panic(err)
	}

	var groupJid *types.JID
	for _, group := range groups {
		if group.Name == cfg.GroupName {
			groupJid = &group.JID
		}
	}

	if groupJid == nil {
		logger.Info(ctx, "group not found")
		return
	}

	logToGroup := func(message string) {
		formattedMessage := botPrefix + message
		waClient.SendMessage(ctx, *groupJid, &waE2E.Message{
			Conversation: &formattedMessage,
		})
	}

	logToGroup("hello")

	sc, disconnect, err := sftputil.GetConnection(ctx, cfg.SftpUrl, cfg.CheckKnownHosts)
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
	go sftputil.TailFile(ctx, &sftputil.TailFileOptions{
		Path: "/logs/latest.log",
		Callback: func(line string) {
			event, err := gameevents.ParseGameEvent(line)
			if err != nil {
				return
			}

			if event.Action == "joined" {
				logger.Info(ctx, "joined", "username", event.Username)
				logToGroup(fmt.Sprintf("%s %s", event.Username, cfg.JoinMessage))
			}
			if event.Action == "left" {
				logger.Info(ctx, "left", "username", event.Username)
				logToGroup(fmt.Sprintf("%s %s", event.Username, cfg.LeaveMessage))
			}

		},
		IgnoreExisting: true,
		Client:         sc,
	})

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	<-c

	waClient.Disconnect()
}
