package factory

import (
	"io"

	"github.com/myshkin5/netspel/jsonstruct"
)

type Writer interface {
	Init(config jsonstruct.JSONStruct) error
	io.WriteCloser
}

type Reader interface {
	Init(config jsonstruct.JSONStruct) error
	io.ReadCloser
}
