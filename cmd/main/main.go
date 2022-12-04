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

	// in_launch_bot: 5419202121:AAFotPHNAuL9B12NHziyFsWEhIDEfNGd3NU
	// volt : 210985494:AAG-GE6m_JwsU31ZDHti91SNmSbePnTSJLk
	botAPI, _, _ := InitializeApplication("5419202121:AAFotPHNAuL9B12NHziyFsWEhIDEfNGd3NU", logger)
	botAPI.StartTelegramServer(true, 60)

	logger.Info("StartTelegramServer")
}
