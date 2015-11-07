package factory

import (
	"errors"
	"io"
)

var (
	ErrReaderClosed = errors.New("connection closed while waiting for read")
)

type Writer interface {
	Init(config map[string]interface{}) error
	io.Writer
}

type Reader interface {
	Init(config map[string]interface{}) error
	Stop()
	io.Reader
}
