package factory

type Scheme interface {
	Init(config map[string]interface{}) error
	RunWriter(writer Writer)
	RunReader(reader Reader)
}
