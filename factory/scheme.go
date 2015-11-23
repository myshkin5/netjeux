package factory

import "github.com/myshkin5/netspel/jsonstruct"

type Scheme interface {
	Init(config jsonstruct.JSONStruct) error

	RunWriter(writer Writer)
	RunReader(reader Reader)
}
