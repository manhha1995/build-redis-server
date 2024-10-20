package server

import (
	"io"
	"net"
	"sync"
	"time"

	"github.com/hdt3213/godis/lib/sync/atomic"
	log "github.com/sirupsen/logrus"
)

type EchoClient struct {
	con     net.Conn
	waiting sync.WaitGroup
}

type EchoHandler struct {
	activeConnection sync.Map
	closing          atomic.Boolean
}

func MakeEchoHandler() *EchoHandler {
	return &EchoHandler{}
}

// close connection
func (h *EchoClient) Close() error {
	h.waiting.Wait()
	h.con.Close()
	return nil
}

// handle connection
func (h *EchoHandler) HandleConnection(c net.Conn) {
	if h.closing.Get() {
		_ = c.Close()
		return
	}
	client := &EchoClient{
		con: c,
	}
	h.activeConnection.Store(client, struct{}{})
	buffer := make([]byte, 1024)

	for {
		// read connection data
		reader, err := c.Read(buffer)

		if err != nil {
			if err == io.EOF {
				log.Info("connection closed")
				h.activeConnection.Delete(client)
			} else {
				log.Warn(err)
			}
			return
		}

		// sleep echo for a while
		time.Sleep(5 * time.Millisecond)

		// write the reader
		b := buffer[:reader]
		c.Write(b)
		client.waiting.Done()
	}
}

// close echo handler
func (h *EchoHandler) Close() error {
	log.Info("close echo handler")
	h.closing.Set(true)
	h.activeConnection.Range(func(key, value interface{}) bool {
		client := key.(*EchoClient)
		client.con.Close()
		return true
	})
	return nil
}
