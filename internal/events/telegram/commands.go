package events_telegram

import (
	"errors"
	"fmt"
	"net/url"
	"strings"
	"tgbot/internal/clients/telegram"
	"tgbot/internal/errs"
	"tgbot/internal/storage"
)

const (
	StartCmd = "/start"
	RndCmd   = "/rnd"
	HelpCmd  = "/help"
)

func (b *Bot) doCmd(text string, chatID int, username string) error {
	const op = "events_telegram.doCmd"

	text = strings.TrimSpace(text)

	b.log.Info("Got new message", "text", text, "username", username)

	// save page: http://...
	// random page: /rnd
	// help: /help
	// start: /start: hi + help

	if isAddCmd(text) {
		return b.savePage(chatID, text, username)

	}

	switch text {
	case StartCmd:
		return b.sendStart(chatID)
	case RndCmd:
		return b.sendRandom(chatID, username)
	case HelpCmd:
		return b.sendHelp(chatID)
	default:
		return b.tg.SendMessage(chatID, msgUnknownCommand)
	}
}

func (b *Bot) savePage(chatID int, pageURL string, username string) (err error) {
	const op = "events_telegram.Processor.savePage"
	defer func() {
		if err != nil {
			err = fmt.Errorf("%s: %w", op, err)
		}
	}()

	sendMsg := NewMessageSender(chatID, b.tg)

	page := &storage.Page{
		URL:      pageURL,
		Username: username,
	}

	isExists, err := b.storage.IsExists(page)
	if err != nil {
		return err
	}

	if isExists {
		return sendMsg(msgAlreadyExists)
	}

	if err := b.storage.Save(page); err != nil {
		return err
	}

	if err := sendMsg(msgSaved); err != nil {
		return err
	}

	return nil
}

func (b *Bot) sendRandom(chatID int, username string) (err error) {
	const op = "events_telegram.Processor.sendRandom"
	defer func() {
		if err != nil {
			err = fmt.Errorf("%s: %w", op, err)
		}
	}()

	sendMsg := NewMessageSender(chatID, b.tg)

	page, err := b.storage.PickRandom(username)
	if err != nil {
		if errors.Is(err, errs.ErrNoSavedPages) {
			return sendMsg(msgNoSavedPages)
		}
		return err
	}

	if err := sendMsg(page.URL); err != nil {
		return err
	}

	return b.storage.Remove(page)

}

func (b *Bot) sendHelp(chatID int) error {
	return NewMessageSender(chatID, b.tg)(msgHelp)
}

func (b *Bot) sendStart(chatID int) error {
	return NewMessageSender(chatID, b.tg)(msgHello)
}

func NewMessageSender(chatID int, tg *telegram.Client) func(string) error {
	return func(msg string) error {
		return tg.SendMessage(chatID, msg)
	}
}

func isAddCmd(text string) bool {
	return isURL(text)
}

func isURL(text string) bool {
	u, err := url.Parse(text)

	return err == nil && u.Host != ""
}
