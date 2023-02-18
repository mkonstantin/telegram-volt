package main

import (
	"fmt"
	"os"
	"telegram-api/pkg/log"
)

func main() {
	logger, err := log.NewLogger(true, "info", "telegram-api")
	if err != nil {
		fmt.Fprintf(os.Stderr, "cannot init log %s", err)
		return
	}
	logger.Info("app starting")

	botAPI, _, _ := InitializeApplication("5566428356:AAH6_BR_A8O_33VEZTw2PNtHHTtaEwB9Rrk", logger)
	botAPI.StartAsyncScheduler()
	botAPI.StartTelegramServer(true, 60)

	logger.Info("StartTelegramServer")
}
