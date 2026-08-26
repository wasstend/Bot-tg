package events_telegram

import (
	"errors"
	"fmt"
	"log"
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

func (p *Processor) doCmd(text string, chatID int, username string) error {
	text = strings.TrimSpace(text)

	log.Printf("Got message '%s' from '%s'", text, username)

	// save page: http://...
	// random page: /rnd
	// help: /help
	// start: /start: hi + help

	if isAddCmd(text) {
		return p.savePage(chatID, text, username)
	}

	switch text {
	case StartCmd:
		return p.sendStart(chatID)
	case RndCmd:
		return p.sendRandom(chatID, username)
	case HelpCmd:
		return p.sendHelp(chatID)
	default:
		return p.tg.SendMessage(chatID, msgUnknownCommand)
	}
}

func (p *Processor) savePage(chatID int, pageURL string, username string) (err error) {
	const op = "events_telegram.Processor.savePage"
	defer func() {
		if err != nil {
			err = fmt.Errorf("%s: %w", op, err)
		}
	}()

	sendMsg := NewMessageSender(chatID, p.tg)

	page := &storage.Page{
		URL:      pageURL,
		Username: username,
	}

	isExists, err := p.storage.IsExists(page)
	if err != nil {
		return err
	}

	if isExists {
		return sendMsg(msgAlreadyExists)
	}

	if err := p.storage.Save(page); err != nil {
		return err
	}

	if err := sendMsg(msgSaved); err != nil {
		return err
	}

	return nil
}

func (p *Processor) sendRandom(chatID int, username string) (err error) {
	const op = "events_telegram.Processor.sendRandom"
	defer func() {
		if err != nil {
			err = fmt.Errorf("%s: %w", op, err)
		}
	}()

	sendMsg := NewMessageSender(chatID, p.tg)

	page, err := p.storage.PickRandom(username)
	if err != nil {
		if errors.Is(err, errs.ErrNoSavedPages) {
			return sendMsg(msgNoSavedPages)
		}
		return err
	}

	if err := sendMsg(page.URL); err != nil {
		return err
	}

	return p.storage.Remove(page)

}

func (p *Processor) sendHelp(chatID int) error {
	return NewMessageSender(chatID, p.tg)(msgHelp)
}

func (p *Processor) sendStart(chatID int) error {
	return NewMessageSender(chatID, p.tg)(msgHello)
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
