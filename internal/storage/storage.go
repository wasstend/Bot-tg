package storage

import (
	"crypto/sha1"
	"encoding/hex"
	"fmt"
)

type Storage interface {
	Save(p *Page) error
	PickRandom(userName string) (*Page, error)
	Remove(p *Page) error
	IsExists(p *Page) (bool, error)
}

type Page struct {
	URL      string
	Username string
}

func (p Page) Hash() (string, error) {
	const op = "storage.Page.Hash"

	hash := sha1.New()

	if _, err := hash.Write([]byte(p.URL)); err != nil {
		return "", fmt.Errorf("%s: %w", op, err)
	}

	if _, err := hash.Write([]byte(p.Username)); err != nil {
		return "", fmt.Errorf("%s: %w", op, err)
	}

	return hex.EncodeToString(hash.Sum(nil)), nil
}
