package mocks

import (
	"io"

	"github.com/myshkin5/netspel/jsonstruct"
)

type ReadMessage struct {
	Buffer []byte
	Error  error
}

type MockReader struct {
	ReadMessages chan ReadMessage
}

func NewMockReader() *MockReader {
	return &MockReader{
		ReadMessages: make(chan ReadMessage, 10000),
	}
}

func (m *MockReader) Init(config jsonstruct.JSONStruct) error {
	return nil
}

func (m *MockReader) Read(message []byte) (int, error) {
	readMessage := <-m.ReadMessages
	bytesRead := 0
	if readMessage.Error == nil {
		bytesRead = len(readMessage.Buffer)
		if len(message) < bytesRead {
			bytesRead = len(message)
		}
		copy(message, readMessage.Buffer)
	}
	return bytesRead, readMessage.Error
}

func (m *MockReader) Close() error {
	m.ReadMessages <- ReadMessage{Buffer: []byte{}, Error: io.EOF}
	return nil
}
