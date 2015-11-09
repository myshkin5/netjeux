package udp

import (
	"fmt"
	"net"

	"netspel/factory"
	"netspel/jsonstruct"
)

type Reader struct {
	connection *net.UDPConn
}

func (r *Reader) Init(config jsonstruct.JSONStruct) error {
	port, ok := config.Int(Port)
	if !ok {
		return fmt.Errorf("%s must be specified in the config additional section", Port)
	}

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
		return 0, factory.ErrReaderClosed
	}

	return count, err
}

func (r *Reader) Stop() {
	r.connection.Close()
}
