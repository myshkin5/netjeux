package mocks

import "github.com/myshkin5/netspel/jsonstruct"

type MockWriter struct {
	Messages chan []byte
}

func NewMockWriter() *MockWriter {
	return &MockWriter{
		Messages: make(chan []byte, 10000),
	}
}

func (m *MockWriter) Init(config jsonstruct.JSONStruct) error {
	return nil
}

func (m *MockWriter) Write(message []byte) (int, error) {
	m.Messages <- message
	return len(message), nil
}

func (m *MockWriter) Close() error {
	return nil
}
