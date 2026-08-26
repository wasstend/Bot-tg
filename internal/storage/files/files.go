package files

import (
	"encoding/gob"
	"fmt"
	"math/rand"
	"os"
	"path/filepath"
	"tgbot/internal/errs"
	"tgbot/internal/storage"
)

const (
	defaultPerm = 0774
)

type Storage struct {
	basePath string
}

func New(basePath string) *Storage {
	return &Storage{basePath: basePath}
}

func (s *Storage) Save(page *storage.Page) error {
	const op = "files.Save"

	filePath := filepath.Join(s.basePath, page.Username)

	if err := os.MkdirAll(filePath, defaultPerm); err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	fileName, err := fileName(page)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	filePath = filepath.Join(filePath, fileName)

	file, err := os.Create(filePath)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	defer file.Close()

	if err := gob.NewEncoder(file).Encode(page); err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	return nil

}

func (s *Storage) PickRandom(userName string) (*storage.Page, error) {
	const op = "files.PickRandom"

	path := filepath.Join(s.basePath, userName)

	files, err := os.ReadDir(path)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	if len(files) == 0 {
		return nil, errs.ErrNoSavedPages
	}

	randomFileNumber := rand.Intn(len(files))

	file := files[randomFileNumber]

	return s.decodePage(filepath.Join(path, file.Name()))

}

func (s *Storage) Remove(p *storage.Page) error {
	const op = "files.Remove"

	fileName, err := fileName(p)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	path := filepath.Join(s.basePath, p.Username, fileName)

	if err := os.Remove(path); err != nil {
		return fmt.Errorf("%s: %s: %w", op, path, err)
	}

	return nil
}

func (s *Storage) IsExists(p *storage.Page) (bool, error) {
	const op = "files.IsExists"

	fileName, err := fileName(p)
	if err != nil {
		return false, fmt.Errorf("%s: %w", op, err)
	}

	path := filepath.Join(s.basePath, p.Username, fileName)

	switch _, err := os.Stat(path); {
	case os.IsNotExist(err):
		return false, nil
	case err != nil:
		return false, fmt.Errorf("%s: %w", op, err)
	}

	return true, nil
}

func (s *Storage) decodePage(filePath string) (*storage.Page, error) {
	const op = "files.decodePage"

	file, err := os.Open(filePath)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}
	defer file.Close()

	var p storage.Page

	if err := gob.NewDecoder(file).Decode(&p); err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return &p, nil
}

func fileName(p *storage.Page) (string, error) {
	return p.Hash()
}
