package client

import (
	"log"
	"net"
	"net/http"
	"sync"

	"github.com/golang/protobuf/ptypes/duration"
	"github.com/hdt3213/godis/interface/redis"
)

type Client struct {
	request   *http.Request
	ipAddress string
	conn      net.Conn
	status    int32
	sync      *sync.WaitGroup
	heartbeat bool // check the liveness
	reply     redis.Reply
	duration  duration.Duration
}

func NewClient(addr string) (*Client, error) {
	log.Print("create the new client")
	conn, err := net.Dial("tcp", addr)
	if err != nil {
		panic(err)
	}
	return &Client{
		ipAddress: addr,
		conn:      conn,
		sync:      &sync.WaitGroup{},
	}, nil
}

//start client request
func (c *Client) Start() {
	log.Println("start the client session")
	// handle concurrent
	go func() {
		c.handleRequest()
		c.Heartbeat()
	}()
}

func (c *Client) handleRequest() {

}

func (c *Client) Heartbeat() {
	if c.duration.AsDuration().Milliseconds() >= 3000 {
		c.doHeartbeat()
	}
}

func (c *Client) doHeartbeat() {

}
