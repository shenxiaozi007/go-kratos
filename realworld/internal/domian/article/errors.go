package article

import "errors"

var (
	ErrTitleNotFound   = errors.New("article title not found")
	ErrContentNotFound = errors.New("article content not found")
	ErrAlreadyPub      = errors.New("article already published")
)
