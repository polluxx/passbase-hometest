package main

import (
	"passbase-hometest/config"
	"passbase-hometest/domain/app"
	"passbase-hometest/domain/database"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"

	_ "passbase-hometest/cmd/server/docs"
)

var (
	log          = buildZapLog()
	configFolder = "../../config"
)

func buildZapLog() *zap.SugaredLogger {
	cfg := zap.NewDevelopmentConfig()
	cfg.EncoderConfig.EncodeLevel = zapcore.LowercaseColorLevelEncoder
	// Always use the console encoder (JSON is default for zap.NewProductionConfig())
	cfg.Encoding = "console"
	logger, _ := cfg.Build()
	zap.ReplaceGlobals(logger)
	return zap.S().Named("converter app")
}

func main() {

	configuration, err := config.LoadFromPath(configFolder)
	if err != nil {
		log.Fatalf("can't load config: %s", err)
	}

	if configuration == nil {
		log.Fatal("config is nil")
	}

	repository := database.Connect(configuration.Database, database.SQLite)

	application := app.NewApp(repository)
	err = application.Router.Engine.Run()
	if err != nil {
		log.Fatalf("can't run server: %s", err)
	}
}
