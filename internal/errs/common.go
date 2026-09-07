package errs

import "errors"

var (
	ErrNoSavedPages  = errors.New("no saved pages")
	ErrFileNotExists = errors.New("file does not exist")
	ErrUserNotFound  = errors.New("user not found")
	ErrPageNotFound  = errors.New("page not found")
	// ErrPageAlreadyExists = errors.New("page already exists")
)
