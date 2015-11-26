package sse

import (
	"fmt"
	"net/http"

	"github.com/myshkin5/jsonstruct"
	vitosse "github.com/vito/go-sse/sse"
)

const (
	prefix = "sse."

	Port             = prefix + "port"
	RemoteWriterAddr = prefix + "remote-writer-addr"

	DefaultPort             = 38208
	DefaultRemoteWriterAddr = "localhost"
)

type Reader struct {
	sseReader *vitosse.ReadCloser
}

func (r *Reader) Init(config jsonstruct.JSONStruct) error {
	port := config.IntWithDefault(Port, DefaultPort)
	remoteAddr := config.StringWithDefault(RemoteWriterAddr, DefaultRemoteWriterAddr)
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
