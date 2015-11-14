package sse

import (
	"fmt"
	"net/http"

	"github.com/myshkin5/netspel/jsonstruct"
	vitosse "github.com/vito/go-sse/sse"
)

const (
	prefix = "sse."

	RemoteAddr = prefix + "remote-addr"
	Port       = prefix + "port"
)

type Reader struct {
	sseReader *vitosse.ReadCloser
}

func (r *Reader) Init(config jsonstruct.JSONStruct) error {
	port, ok := config.Int(Port)
	if !ok {
		return fmt.Errorf("%s must be specified in the config additional section", Port)
	}

	remoteAddr, ok := config.String(RemoteAddr)
	if !ok {
		return fmt.Errorf("%s must be specified in the config additional section", RemoteAddr)
	}

	resp, err := http.DefaultClient.Get(fmt.Sprintf("http://%s:%d/", remoteAddr, port))
	if err != nil {
		return err
	}

	r.sseReader = vitosse.NewReadCloser(resp.Body)

	return nil
}

func (r *Reader) Read(message []byte) (int, error) {
	event, err := r.sseReader.Next()
	if err != nil {
		return 0, err
	}
	return copy(message, event.Data), nil
}

func (r *Reader) Close() error {
	return r.sseReader.Close()
}
