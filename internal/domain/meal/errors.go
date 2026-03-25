package meal

import "errors"

var (
	ErrEmptyName = errors.New("meal name cannot be empty")
	ErrNotFound  = errors.New("meal not found")
)
