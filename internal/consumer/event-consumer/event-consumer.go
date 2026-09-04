package event_consumer

import (
	"log"
	"log/slog"
	"tgbot/internal/events"
	"time"
)

type Consumer struct {
	fetcher   events.Fetcher
	processor events.Processor
	batchSize int
	log       *slog.Logger
}

func New(
	fetcher events.Fetcher,
	processor events.Processor,
	batchSize int,
	log *slog.Logger,
) Consumer {
	return Consumer{
		fetcher:   fetcher,
		processor: processor,
		batchSize: batchSize,
		log:       log,
	}
}

func (c Consumer) Start() error {

	c.log.Info("starting bot...")
	for {
		gotEvents, err := c.fetcher.Fetch(c.batchSize)
		if err != nil {
			c.log.Error("consumer error", "error", err)
			time.Sleep(1 * time.Second)
			continue
			// попробовать сделать retry логику
		}

		if len(gotEvents) == 0 {
			time.Sleep(1 * time.Second)
			continue
		}

		if err := c.handleEvents(gotEvents); err != nil {
			log.Print(err)
			continue
		}
	}
}

/*
Возможные проблемы:
	1. Потеря событий:
		ретраи, возвращение в хранилище, fallback, подтверждение для фетчера (чтобы не делал сдвиг всегда)
	2. Обработка всей пачки событий:
		останавливаться после первой ошибки, вести счетчик ошибок
	3. Параллельная обработка (строка 56):
		воркер пул
*/

func (c Consumer) handleEvents(events []events.Event) error {
	const op = "event_consumer.handleEvents"

	for _, event := range events {
		c.log.Info("Got new event", "text", event.Text)

		if err := c.processor.Process(event); err != nil {
			c.log.Error("Can't handle event", "error", err)

			continue
		}
	}

	return nil
}
