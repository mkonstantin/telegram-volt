package main

import (
	"context"
	"fmt"
	"github.com/joho/godotenv"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/exporters/jaeger"
	"go.opentelemetry.io/otel/sdk/resource"
	tracesdk "go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.17.0"
	"go.uber.org/zap"
	"os"
	"strconv"
	"strings"
	"telegram-api/config"
	"telegram-api/pkg/log"
	"time"
)

const (
	appName      = "telegram-volt"
	environment  = "production"
	batchTimeout = 5
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

	adminsStr := os.Getenv("ADMINS")
	admins := strings.Split(adminsStr, ",")

	cfg := config.AppConfig{
		Username:              os.Getenv("USERNAME"),
		Password:              os.Getenv("PASSWORD"),
		Host:                  os.Getenv("HOST"),
		Port:                  port,
		Database:              os.Getenv("DATABASE"),
		MaxOpenConnections:    mOpenConn,
		MaxIdleConnections:    mIdleConn,
		ConnectionMaxLifeTime: lifeTime,
		Version:               os.Getenv("VERSION"),
		Admins:                admins,
	}

	flush, err := createTracing(logger)
	if err != nil {
		logger.Error("can't init jaeger tracer", zap.Error(err))
		return
	}
	logger.Info("Start Jaeger tracer")
	defer flush()

	logger.Info(fmt.Sprintf("DB host: %s", cfg.Host))
	botAPI, _, _ := InitializeApplication(os.Getenv("KEY"), cfg, logger)
	botAPI.StartAsyncScheduler()
	botAPI.StartTelegramServer(true, 60)

	logger.Info("StartTelegramServer")
}

func createTracing(logger *zap.Logger) (func(), error) {
	tp, err := makeTracerProvider()
	if err != nil {
		logger.Error("TraceProvider creation error", zap.Error(err))
		return nil, err
	}

	// Register our TracerProvider as the global so any imported
	// instrumentation in the future will default to using it.
	otel.SetTracerProvider(tp)

	// Cleanly shutdown and flush telemetry when the application exits.
	closeFn := func() {
		if err = tp.Shutdown(context.Background()); err != nil {
			logger.Error("TraceProvider CloseFn error", zap.Error(err))
		}
	}

	return closeFn, nil
}

func makeTracerProvider() (*tracesdk.TracerProvider, error) {

	jaegerHost := os.Getenv("JAEGER_HOST")
	jaegerPort := os.Getenv("JAEGER_PORT")

	// http://localhost:14268/api/traces
	url := fmt.Sprintf("http://%s:%s/api/traces", jaegerHost, jaegerPort)

	// Create the Jaeger exporter
	exp, err := jaeger.New(jaeger.WithCollectorEndpoint(jaeger.WithEndpoint(url)))
	if err != nil {
		return nil, err
	}

	tp := tracesdk.NewTracerProvider(
		// Always be sure to batch in production.
		tracesdk.WithBatcher(exp, tracesdk.WithBatchTimeout(batchTimeout*time.Second)),
		// Record information about this application in a Resource.
		tracesdk.WithResource(resource.NewWithAttributes(
			semconv.SchemaURL,
			semconv.ServiceName(appName),
			attribute.String("environment", environment),
		)),
	)

	return tp, nil
}
