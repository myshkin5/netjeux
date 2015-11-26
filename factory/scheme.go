package factory

import "github.com/myshkin5/jsonstruct"

type Scheme interface {
	Init(config jsonstruct.JSONStruct) error

	RunWriter(writer Writer)
	RunReader(reader Reader)
}
