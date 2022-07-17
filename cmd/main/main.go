package main

import (
	"fmt"
	"telegram-api/internal/infrastructure_layer/telegram"
)

func main() {
	fmt.Println("sdfs")

	botAPI, _, _ := InitializeApplication("210985494:AAG-GE6m_JwsU31ZDHti91SNmSbePnTSJLk")
	telegram.StartTelegramServer(botAPI.BotAPI, true, 60)

}
