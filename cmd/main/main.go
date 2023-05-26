package main

import (
	"context"
	"fmt"
	"github.com/joho/godotenv"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/jaeger"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/sdk/resource"
	"go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.4.0"
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

	logger.Info(fmt.Sprintf("DB host: %s", cfg.Host))
	botAPI, _, _ := InitializeApplication(os.Getenv("KEY"), cfg, logger)
	botAPI.StartAsyncScheduler()
	botAPI.StartTelegramServer(true, 60)

	flush, err := initTracer(logger)
	if err != nil {
		logger.Error("can't init jaeger tracer", zap.Error(err))
	}
	defer flush()

	logger.Info("StartTelegramServer")
}

// http://localhost:14268/api/traces

func initTracer(logger *zap.Logger) (func(), error) {
	sampler := trace.TraceIDRatioBased(1)

	jaegerHost := os.Getenv("JAEGER_AGENT_HOST")
	if jaegerHost == "" {
		jaegerHost = "localhost"
	}
	jaegerPort := os.Getenv("JAEGER_AGENT_PORT")
	if jaegerPort == "" {
		jaegerPort = "14268"
	}

	exporter, err := jaeger.New(
		jaeger.WithAgentEndpoint(
			jaeger.WithAgentHost(jaegerHost),
			jaeger.WithAgentPort(jaegerPort),
		),
	)
	if err != nil {
		return nil, fmt.Errorf("error occurred while trying to setup tracing exporter: %w", err)
	}

	tp := trace.NewTracerProvider(
		trace.WithBatcher(exporter, trace.WithBatchTimeout(batchTimeout*time.Second)),
		trace.WithSampler(sampler),
		trace.WithResource(
			resource.NewWithAttributes(
				semconv.SchemaURL,
				semconv.ServiceNameKey.String(appName),
			),
		),
	)
	otel.SetTextMapPropagator(propagation.TraceContext{})
	otel.SetTracerProvider(tp)

	closeFn := func() {
		if err := tp.Shutdown(context.Background()); err != nil {
			logger.Error("can't init jaeger tracer", zap.Error(err))
		}
	}

	return closeFn, nil
}
