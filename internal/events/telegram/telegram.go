package events_telegram

import (
	"errors"
	"fmt"
	"tgbot/internal/clients/telegram"
	"tgbot/internal/events"
	"tgbot/internal/storage"
)

type Processor struct {
	tg      *telegram.Client
	offset  int
	storage storage.Storage
}

type Meta struct {
	ChatID   int
	Username string
}

var (
	ErrUnknownEventType = errors.New("unknows event type")
	ErrUnknownMetaType  = errors.New("unknown meta type")
)

func New(client *telegram.Client, storage storage.Storage) *Processor {
	return &Processor{
		tg:      client,
		storage: storage,
	}
}

func (p *Processor) Fetch(limit int) ([]events.Event, error) {
	const op = "events_telegram.Processor.Fetch"

	updates, err := p.tg.Updates(p.offset, limit)
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

	p.offset = updates[len(updates)-1].ID + 1

	return res, nil

}

func (p *Processor) Process(event events.Event) error {
	switch event.Type {
	case events.Message:
		return p.processMessage(event)
	default:
		return fmt.Errorf("process message: %w", ErrUnknownEventType)
	}
}

func (p *Processor) processMessage(event events.Event) error {
	meta, err := meta(event)
	if err != nil {
		return fmt.Errorf("process message: %w", err)
	}

	if err := p.doCmd(event.Text, meta.ChatID, meta.Username); err != nil {
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
