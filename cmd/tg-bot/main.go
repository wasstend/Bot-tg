package main

import (
	"log"
	"tgbot/internal/clients/telegram"
	event_consumer "tgbot/internal/consumer/event-consumer"
	events_telegram "tgbot/internal/events/telegram"
	"tgbot/internal/storage/files"
	"tgbot/internal/token"
)

const (
	tgBotHost       = "api.telegram.org"
	storageBasePath = "files_storage"
	batchSize       = 100
)

func main() {

	tgClient := telegram.New(tgBotHost, token.TokenMust())

	storage := files.New(storageBasePath)

	tgBot := events_telegram.New(tgClient, storage)

	consumer := event_consumer.New(tgBot, tgBot, batchSize)

	if err := consumer.Start(); err != nil {
		log.Fatal("service is stopped", err)
	}
}
