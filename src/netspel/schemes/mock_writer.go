package schemes

type MockWriter struct {
	Messages chan []byte
}

func NewMockWriter() *MockWriter {
	return &MockWriter{
		Messages: make(chan []byte, 10000),
	}
}

func (m *MockWriter) Write(message []byte) (int, error) {
	m.Messages <- message
	return len(message), nil
}
