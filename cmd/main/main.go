package main

import (
	"fmt"
	"os"
	"telegram-api/pkg/log"
)

func main() {
	logger, err := log.NewLogger(true, "info", "truck-api")
	if err != nil {
		fmt.Fprintf(os.Stderr, "cannot init log %s", err)
		return
	}
	logger.Info("app starting")

	botAPI, _, _ := InitializeApplication("5419202121:AAFotPHNAuL9B12NHziyFsWEhIDEfNGd3NU", logger)
	botAPI.StartTelegramServer(true, 60)

	logger.Info("StartTelegramServer")
}
