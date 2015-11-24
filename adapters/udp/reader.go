package udp

import (
	"fmt"
	"io"
	"net"

	"github.com/myshkin5/netspel/jsonstruct"
)

type Reader struct {
	connection *net.UDPConn
}

func (r *Reader) Init(config jsonstruct.JSONStruct) error {
	port := config.IntWithDefault(Port, DefaultPort)
	laddr, err := net.ResolveUDPAddr("udp4", fmt.Sprintf(":%d", port))
	if err != nil {
		return err
	}

	r.connection, err = net.ListenUDP("udp4", laddr)
	if err != nil {
		return err
	}

	return nil
}

func (r *Reader) Read(message []byte) (int, error) {
	count, err := r.connection.Read(message)
	opErr, ok := err.(*net.OpError)
	if err != nil && ok && opErr.Err.Error() == "use of closed network connection" {
		return 0, io.EOF
	}

	return count, err
}

func (r *Reader) Close() error {
	return r.connection.Close()
}
