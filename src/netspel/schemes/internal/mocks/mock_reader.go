package mocks

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
