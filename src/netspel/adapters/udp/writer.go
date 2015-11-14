package udp

import (
	"fmt"
	"net"
	"strconv"

	"netspel/jsonstruct"
)

const (
	prefix = "udp."

	RemoteAddr = prefix + "remote-addr"
	Port       = prefix + "port"
)

type Writer struct {
	connection *net.UDPConn
}

func (w *Writer) Init(config jsonstruct.JSONStruct) error {
	remoteAddr, ok := config.String(RemoteAddr)
	if !ok {
		return fmt.Errorf("%s must be specified in the config additional section", RemoteAddr)
	}

	port, ok := config.Int(Port)
	if !ok {
		return fmt.Errorf("%s must be specified in the config additional section", Port)
	}

	raddr, err := net.ResolveUDPAddr("udp4", net.JoinHostPort(remoteAddr, strconv.Itoa(port)))
	if err != nil {
		return err
	}

	w.connection, err = net.DialUDP("udp4", nil, raddr)
	if err != nil {
		return err
	}

	return nil
}

func (w *Writer) Write(message []byte) (int, error) {
	return w.connection.Write(message)
}

func (w *Writer) Close() error {
	return w.connection.Close()
}
