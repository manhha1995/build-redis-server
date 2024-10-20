package server

import (
	"net"
	"os"
	"sync"
	"sync/atomic"
	"time"

	log "github.com/sirupsen/logrus"
)

// Config stores tcp server properties
type SignalHandler struct {
	status         bool
	ipAddress      string
	timeout        time.Duration
	maxConnections uint32
}

const MAX_THREAD_NUMBERS = 10
const WATING_STATUS = 1
const RUNNING_STATUS = 2
const TERMINATED_STATUS = 3
const TRANSACTION_STATUS = 4

var connectedClients map[int]*SignalHandler
var eStatus int32 = WATING_STATUS

func init() {
	connectedClients = make(map[int]*SignalHandler)
}

func waitForSignal(signalHandler *SignalHandler, signal chan os.Signal, s *sync.WaitGroup) {
	closeChan := make(chan struct{})
	defer s.Wait()
	go func() {
		<-signal
		log.Printf("waiting for signal from server %v")
		close(signal)
	}()

	listener, err := net.Listen("tcp", signalHandler.ipAddress)
	if err != nil {
		panic("stop the server")
	}
	log.Println("bind: %s, start listening address", signalHandler.ipAddress)
	runAsync(listener, s, closeChan)

}

func runAsync(listener net.Listener, s *sync.WaitGroup, closeChan chan struct{}) {
	err := make(chan error, 1)
	defer s.Done()
	go func() {
		select {
		case <-closeChan:
			log.Info("close channel")
		case <-err:
			log.Error("accept the error %s", err)
			break
		}
	}()
	defer func() {
		atomic.StoreInt32(&eStatus, TERMINATED_STATUS)
	}()

	log.Println("starting invoke a TCP server on")
	for {
		_, errChan := listener.Accept()
		if errChan != nil {
			if ne, ok := errChan.(net.Error); ok && ne.Timeout() {
				log.Println("accept occurs tmporary error: %v", errChan)
				time.Sleep(5 * time.Millisecond)
				continue
			}
			err <- errChan
		}
	}
}
