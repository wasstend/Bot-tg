package main

import (
	"tgbot/internal/clients/telegram"
	events_telegram "tgbot/internal/events/telegram"
	"tgbot/internal/storage/files"
	"tgbot/internal/token"
)

const (
	tgBotHost       = "api.telegram.org"
	storageBasePath = "storage"
)

func main() {

	tgClient := telegram.New(tgBotHost, token.TokenMust())

	storage := files.New(storageBasePath)

	eventsHandler := events_telegram.New(tgClient, storage)
	_ = eventsHandler

	// consumer.Start(fetcher, processor)
}
