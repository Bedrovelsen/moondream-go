package moondream

import "errors"

var (
	ErrInvalidAPIKey   = errors.New("invalid API key")
	ErrImageNotFound   = errors.New("image not found")
	ErrAPIRequestFailed = errors.New("API request failed")
)
