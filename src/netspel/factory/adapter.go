package factory

import (
	"errors"
	"io"
)

var (
	ErrReaderClosed = errors.New("connection closed while waiting for read")
)

type Writer interface {
	Init(Config) error
	io.Writer
}

type Reader interface {
	Init(Config) error
	Stop()
	io.Reader
}
