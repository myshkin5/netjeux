package factory

import (
	"fmt"
	"reflect"
)

var (
	WriterManager *InstanceManager
	ReaderManager *InstanceManager
	SchemeManager *InstanceManager
)

func init() {
	WriterManager = NewInstanceManager()
	ReaderManager = NewInstanceManager()
	SchemeManager = NewInstanceManager()
}

type InstanceManager struct {
	types map[string]reflect.Type
}

func NewInstanceManager() *InstanceManager {
	return &InstanceManager{
		types: make(map[string]reflect.Type),
	}
}

func (m *InstanceManager) RegisterType(name string, instanceType reflect.Type) {
	m.types[name] = instanceType
}

func (m *InstanceManager) CreateInstance(name string) (reflect.Value, error) {
	instanceType, ok := m.types[name]
	if !ok {
		return reflect.Value{}, fmt.Errorf("Type not found, %s", name)
	}

	return reflect.New(instanceType), nil
}

func CreateWriter(name string) (Writer, error) {
	value, err := WriterManager.CreateInstance(name)
	if err != nil {
		return nil, err
	}

	writerValue, ok := value.Interface().(Writer)
	if !ok {
		return nil, fmt.Errorf("Type does not implement writer interface, %s", name)
	}

	return writerValue, nil
}

func CreateReader(name string) (Reader, error) {
	value, err := ReaderManager.CreateInstance(name)
	if err != nil {
		return nil, err
	}

	readerValue, ok := value.Interface().(Reader)
	if !ok {
		return nil, fmt.Errorf("Type does not implement Reader interface, %s", name)
	}

	return readerValue, nil
}

func CreateScheme(name string) (Scheme, error) {
	value, err := SchemeManager.CreateInstance(name)
	if err != nil {
		return nil, err
	}

	schemeValue, ok := value.Interface().(Scheme)
	if !ok {
		return nil, fmt.Errorf("Type does not implement Scheme interface, %s", name)
	}

	return schemeValue, nil
}
