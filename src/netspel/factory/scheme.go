package factory

type Scheme interface {
	Init(Config) error
	RunWriter(writer Writer)
	RunReader(reader Reader)
}
