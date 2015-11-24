package sse

import (
	"errors"
	"fmt"
	"net/http"
	"sync"

	"github.com/myshkin5/netspel/jsonstruct"
	"github.com/myshkin5/netspel/logs"
	vitosse "github.com/vito/go-sse/sse"
)

type Writer struct {
	server    *http.Server
	messages  chan []byte
	responses chan response
	readers   sync.WaitGroup
}

type response struct {
	count int
	err   error
}

func (w *Writer) Init(config jsonstruct.JSONStruct) error {
	port := config.IntWithDefault(Port, DefaultPort)
	server := &http.Server{
		Addr:    fmt.Sprintf("localhost:%d", port),
		Handler: http.HandlerFunc(w.handle),
	}

	go func() {
		err := server.ListenAndServe()
		if err != nil {
			logs.Logger.Warning("Error when starting server, %e", err.Error())
		}
	}()

	w.messages = make(chan []byte)
	w.responses = make(chan response, 1)

	return nil
}

func (w *Writer) Write(message []byte) (int, error) {
	w.messages <- message
	resp := <-w.responses
	return resp.count, resp.err
}

func (w *Writer) Close() error {
	select {
	case <-w.messages:
	default:
	}
	close(w.messages)
	w.responses <- response{
		count: 0,
		err:   errors.New("Writer closing"),
	}
	w.readers.Wait()
	return nil
}

func (w *Writer) handle(rw http.ResponseWriter, r *http.Request) {
	w.readers.Add(1)
	defer w.readers.Done()

	if r.RequestURI == "/ready" {
		rw.WriteHeader(http.StatusOK)
		return
	}

	rw.Header().Set("Content-Type", "text/event-stream; charset=utf-8")
	rw.Header().Set("Cache-Control", "no-cache")
	rw.Header().Set("Connection", "keep-alive")
	rw.WriteHeader(http.StatusOK)
	rw.Write([]byte("\n"))
	flusher := rw.(http.Flusher)
	flusher.Flush()

	closeNotifier := rw.(http.CloseNotifier).CloseNotify()

	for {
		select {
		case <-closeNotifier:
			return
		case message, ok := <-w.messages:
			if !ok {
				return
			}

			event := vitosse.Event{
				Data: message,
			}

			err := event.Write(rw)
			flusher.Flush()

			w.responses <- response{
				count: len(message),
				err:   err,
			}
		}
	}
}
