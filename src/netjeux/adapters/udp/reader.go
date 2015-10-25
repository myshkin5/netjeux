package udp

import (
	"fmt"
	"net"
)

type Reader struct {
	connection *net.UDPConn
}

func NewReader(port uint16) (*Reader, error) {
	laddr, err := net.ResolveUDPAddr("udp4", fmt.Sprintf(":%d", port))
	if err != nil {
		return nil, err
	}

	connection, err := net.ListenUDP("udp4", laddr)
	if err != nil {
		return nil, err
	}

	return &Reader{
		connection: connection,
	}, nil
}

func (r *Reader) Read(message []byte) (int, error) {
	return r.connection.Read(message)
}
