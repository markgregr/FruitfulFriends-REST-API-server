package main

import (
	"github.com/Fruitfulfriends-REST-API-server/internal/app"
	"github.com/Fruitfulfriends-REST-API-server/internal/config"
	"github.com/Fruitfulfriends-REST-API-server/internal/lib/logger/handlers/logruspretty"
	"github.com/sirupsen/logrus"
	"os"
)

const (
	envLocal = "local"
	envDev   = "dev"
	envProd  = "prod"
)

func main() {
	cfg := config.MustLoad()

	log := setupLogger(cfg.Env, cfg.LogsPath)

	log.WithField("config", cfg).Info("Application start!")

	application, err := app.New(cfg, log)
	if err != nil {
		panic(err)
	}

	application.Run()

	<-application.Done
	log.Info("Application stopped")
}

func setupLogger(env string, logFilePath string) *logrus.Entry {
	var log = logrus.New()

	logFile, err := os.OpenFile(logFilePath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		panic(err)
	}

	switch env {
	case envLocal:
		return setupPrettySlog(log)
	case envDev:
		log.SetOutput(logFile)
		log.SetFormatter(&logrus.TextFormatter{
			DisableColors: true,
			FullTimestamp: true, // Добавляем временные метки к сообщениям
		})
	case envProd:
		log.SetOutput(logFile)
		log.SetFormatter(&logrus.TextFormatter{
			DisableColors: true,
			FullTimestamp: true, // Добавляем временные метки к сообщениям
		})
		log.SetLevel(logrus.WarnLevel)
	default:
		log.SetOutput(logFile)
		log.SetFormatter(&logrus.TextFormatter{
			DisableColors: true,
			FullTimestamp: true, // Добавляем временные метки к сообщениям
		})
		log.SetLevel(logrus.DebugLevel)
	}

	return logrus.NewEntry(log)
}

func setupPrettySlog(log *logrus.Logger) *logrus.Entry {
	prettyHandler := logruspretty.NewPrettyHandler(os.Stdout)
	log.SetFormatter(prettyHandler)
	return logrus.NewEntry(log)
}
