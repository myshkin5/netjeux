package udp

import (
	"fmt"
	"net"

	"netspel/factory"
)

type Reader struct {
	connection *net.UDPConn
}

func (r *Reader) Init(config factory.Config) error {
	port, ok := config.AdditionalInt(Port)
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
	return r.connection.Read(message)
}
