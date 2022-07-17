package main

import (
	"fmt"
	"telegram-api/internal/service_layer/service"
)

func main() {
	fmt.Println("sdfs")

	botAPI, _, _ := InitializeApplication("210985494:AAG-GE6m_JwsU31ZDHti91SNmSbePnTSJLk")
	service.StartTelegramServer(botAPI.BotAPI, true, 60)

}
