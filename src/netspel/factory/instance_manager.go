package factory

import (
	"fmt"
	"reflect"
)

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
