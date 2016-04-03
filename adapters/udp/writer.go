package udp

import (
	"net"
	"strconv"

	"github.com/myshkin5/jsonstruct"
)

const (
	prefix = ".udp."

	Port             = prefix + "port"
	RemoteReaderAddr = prefix + "remote-reader-addr"

	DefaultPort             = 57955
	DefaultRemoteReaderAddr = "localhost"
)

type Writer struct {
	connection *net.UDPConn
}

func (w *Writer) Init(config jsonstruct.JSONStruct) error {
	port := config.IntWithDefault(Port, DefaultPort)
	remoteAddr := config.StringWithDefault(RemoteReaderAddr, DefaultRemoteReaderAddr)
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
