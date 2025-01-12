package main

import (
	"context"
	"os"
	"os/signal"
	"syscall"

	"github.com/joho/godotenv"
	"go.uber.org/zap"

	"auth/internal/app"
	"auth/internal/config"
)

func main() {
	prdLogger, _ := zap.NewProduction()
	defer prdLogger.Sync()
	logger := prdLogger.Sugar()

	mainCtx := context.Background()
	ctx, cancel := context.WithCancel(mainCtx)
	defer cancel()

	err := godotenv.Load()
	if err != nil {
		logger.Warnf(".env file not found: %v", err)
	}

	config, err := config.NewConfig()
	if err != nil {
		logger.Fatalf("failed to read config: %s", err.Error())
	}

	newApp := app.NewApp(ctx, logger, config)

	if err := newApp.InitDatabase(); err != nil {
		logger.Fatalf("failed to initialize db: %s", err.Error())
	}

	if err = newApp.RunMigrations(); err != nil {
		logger.Fatalf("failed to run migrations: %s", err.Error())
	}

	newApp.InitService()

	if err = newApp.Run(); err != nil {
		logger.Errorf(err.Error())
		return
	}

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGTERM, syscall.SIGINT)
	<-sigChan

	if err = newApp.Shutdown(); err != nil {
		logger.Errorf(err.Error())
		return
	}
}
