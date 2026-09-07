package events_telegram

import (
	"context"
	"errors"
	"fmt"
	"tgbot/internal/clients/telegram"
	"tgbot/internal/events"
	"tgbot/internal/logger"
	"tgbot/internal/storage"
)

type Bot struct {
	tg      *telegram.Client
	offset  int
	storage Storage
	log     *logger.Logger
	ctx     context.Context
}

// Storage interface represents storage of the saved pages(urls) and users
type Storage interface {
	// TODO: context
	SavePage(ctx context.Context, page storage.Page) error
	GetPage(ctx context.Context, page storage.Page) (storage.Page, error)
	CheckPage(ctx context.Context, page storage.Page) (bool, error)
	PickRandomPage(ctx context.Context, user storage.User) (storage.Page, error)

	SaveUser(ctx context.Context, userName string) (storage.User, error)
	GetUser(ctx context.Context, userName string) (storage.User, error)
	CheckUser(ctx context.Context, userName string) (bool, error)

	Remove(ctx context.Context, p storage.Page) error
}

type Meta struct {
	ChatID   int
	Username string
}

var (
	ErrUnknownEventType = errors.New("unknows event type")
	ErrUnknownMetaType  = errors.New("unknown meta type")
)

func New(
	client *telegram.Client,
	storage Storage,
	log *logger.Logger,
	ctx context.Context,
) *Bot {
	return &Bot{
		tg:      client,
		storage: storage,
		log:     log,
		ctx:     ctx,
	}
}

func (b *Bot) Fetch(limit int) ([]events.Event, error) {
	const op = "events_telegram.Processor.Fetch"

	updates, err := b.tg.Updates(b.offset, limit)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	if len(updates) == 0 {
		return nil, nil
	}

	res := make([]events.Event, 0, len(updates))

	for _, update := range updates {
		res = append(res, toEvent(update))
	}

	b.offset = updates[len(updates)-1].ID + 1

	return res, nil

}

func (b *Bot) Process(event events.Event) error {
	switch event.Type {
	case events.Message:
		return b.processMessage(event)
	default:
		return fmt.Errorf("process message: %w", ErrUnknownEventType)
	}
}

func (b *Bot) processMessage(event events.Event) error {
	meta, err := meta(event)
	if err != nil {
		return fmt.Errorf("process message: %w", err)
	}

	if err := b.doCmd(event.Text, meta.ChatID, meta.Username); err != nil {
		return fmt.Errorf("process message: %w", err)
	}

	return nil
}

func meta(event events.Event) (Meta, error) {
	res, ok := event.Meta.(Meta)
	if !ok {
		return Meta{}, fmt.Errorf("get meta: %w", ErrUnknownMetaType)
	}

	return res, nil
}

func toEvent(upd telegram.Update) events.Event {
	updType := fetchType(upd)

	result := events.Event{
		Type: updType,
		Text: fetchText(upd),
	}

	if updType == events.Message {
		result.Meta = Meta{
			ChatID:   upd.Message.Chat.ID,
			Username: upd.Message.From.Username,
		}
	}

	return result
}

func fetchType(upd telegram.Update) events.Type {
	if upd.Message != nil {
		return events.Message
	}
	return events.Unknown
}

func fetchText(upd telegram.Update) string {
	if upd.Message != nil {
		return upd.Message.Text
	}
	return ""
}
