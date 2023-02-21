package main

import (
	"fmt"
	"github.com/joho/godotenv"
	"go.uber.org/zap"
	"os"
	"strconv"
	"telegram-api/config"
	"telegram-api/pkg/log"
)

func main() {
	logger, err := log.NewLogger(true, "info", "telegram-api")
	if err != nil {
		fmt.Fprintf(os.Stderr, "cannot init log %s", err)
		return
	}
	logger.Info("app starting")

	err = godotenv.Load()
	if err != nil {
		logger.Error("Error load env file", zap.Error(err))
		return
	}

	port, err := strconv.Atoi(os.Getenv("PORT"))
	mOpenConn, err := strconv.Atoi(os.Getenv("MAXOPENCONN"))
	mIdleConn, err := strconv.Atoi(os.Getenv("MAXIDLECONN"))
	lifeTime, err := strconv.Atoi(os.Getenv("CONNLIFETIME"))
	if err != nil {
		logger.Error("Error Getenv", zap.Error(err))
		return
	}

	cfg := config.AppConfig{
		Username:              os.Getenv("USERNAME"),
		Password:              os.Getenv("PASSWORD"),
		Host:                  os.Getenv("HOST"),
		Port:                  port,
		Database:              os.Getenv("DATABASE"),
		MaxOpenConnections:    mOpenConn,
		MaxIdleConnections:    mIdleConn,
		ConnectionMaxLifeTime: lifeTime,
	}

	logger.Info(fmt.Sprintf("DB host: %s", cfg.Host))
	botAPI, _, _ := InitializeApplication("5566428356:AAH6_BR_A8O_33VEZTw2PNtHHTtaEwB9Rrk", cfg, logger)
	botAPI.StartAsyncScheduler()
	botAPI.StartTelegramServer(true, 60)

	logger.Info("StartTelegramServer")
}
