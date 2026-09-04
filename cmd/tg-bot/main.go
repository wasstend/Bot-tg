package main

import (
	"os"
	"tgbot/internal/clients/telegram"
	event_consumer "tgbot/internal/consumer/event-consumer"
	events_telegram "tgbot/internal/events/telegram"
	"tgbot/internal/logger"
	"tgbot/internal/storage/files"
	"tgbot/internal/token"
)

const (
	tgBotHost       = "api.telegram.org"
	storageBasePath = "files_storage"
	batchSize       = 100
)

func main() {

	log, mustCloseFile := logger.SetupLogger()
	defer func(file *os.File) {
		err := file.Close()
		if err != nil {
			panic("failed to close logs file: " + err.Error())
		}
	}(mustCloseFile)

	log.Debug("initialized logger")

	tgClient := telegram.New(tgBotHost, token.MustToken())

	storage := files.New(storageBasePath)

	tgBot := events_telegram.New(tgClient, storage, log)

	consumer := event_consumer.New(tgBot, tgBot, batchSize, log)

	if err := consumer.Start(); err != nil {
		log.Error("service is stopped", err)
		os.Exit(1)
	}
}
