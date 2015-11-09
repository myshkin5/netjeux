package factory

import (
	"errors"
	"io"

	"netspel/jsonstruct"
)

var (
	ErrReaderClosed = errors.New("connection closed while waiting for read")
)

type Writer interface {
	Init(config jsonstruct.JSONStruct) error
	io.Writer
}

type Reader interface {
	Init(config jsonstruct.JSONStruct) error
	Stop()
	io.Reader
}
