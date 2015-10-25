package udp

import "net"

type Writer struct {
	connection *net.UDPConn
}

func NewWriter(raddr *net.UDPAddr) (*Writer, error) {
	connection, err := net.DialUDP("udp4", nil, raddr)
	if err != nil {
		return nil, err
	}

	return &Writer{
		connection: connection,
	}, nil
}

func (w *Writer) Write(message []byte) (int, error) {
	return w.connection.Write(message)
}
