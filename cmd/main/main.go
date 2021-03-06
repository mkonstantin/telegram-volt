package main

import (
	"fmt"
	"os"
	"telegram-api/internal/infrastructure_layer/telegram"
	"telegram-api/pkg/log"
)

func main() {
	logger, err := log.NewLogger(true, "info", "truck-api")
	if err != nil {
		fmt.Fprintf(os.Stderr, "cannot init log %s", err)
		return
	}
	logger.Info("app starting")

	botAPI, _, _ := InitializeApplication("210985494:AAG-GE6m_JwsU31ZDHti91SNmSbePnTSJLk", logger)

	telegram.StartTelegramServer(botAPI.BotAPI, true, 60)
	logger.Info("StartTelegramServer")
}
