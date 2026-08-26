package errs

import "errors"

var (
	ErrNoSavedPages  = errors.New("no saved pages")
	ErrFileNotExists = errors.New("file does not exist")
	// ErrPageAlreadyExists = errors.New("page already exists")
)
