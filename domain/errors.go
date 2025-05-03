package domain

import "errors"

var (
	ErrEncoding      = errors.New("encoding error")
	ErrRequest       = errors.New("request error")
	ErrApi           = errors.New("api error")
	ErrCreateRequest = errors.New("create request error")
	ErrEditor        = errors.New("editor error")
)
