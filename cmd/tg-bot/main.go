package main

import (
	"context"
	"tgbot/internal/clients/telegram"
	event_consumer "tgbot/internal/consumer/event-consumer"
	events_telegram "tgbot/internal/events/telegram"
	"tgbot/internal/logger"
	"tgbot/internal/storage/postgres"
	"tgbot/internal/token"
)

const (
	tgBotHost       = "api.telegram.org"
	storageBasePath = "files_storage"
	batchSize       = 100
)

func main() {

	ctx := context.TODO()

	log := logger.New()
	defer log.Close()

	log.Debug("initialized logger")

	tgClient := telegram.New(tgBotHost, token.MustToken())

	//storage := files.New(storageBasePath)

	storage, err := postgres.New(ctx)
	if err != nil {
		log.Fatal("failed to connect to postgres", "error", err)
	} else {
		log.Debug("connected to postgres")
	}

	tgBot := events_telegram.New(tgClient, storage, log)

	consumer := event_consumer.New(tgBot, tgBot, batchSize, log)

	if err := consumer.Start(); err != nil {
		log.Fatal("service is stopped", "error", err)
	}
}
