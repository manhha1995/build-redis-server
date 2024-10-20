package client

import (
	"net"
)

type Client struct {
	conn net.Conn
}

// remote address
func (c *Client) remoteAddr() string {
	return c.conn.RemoteAddr().String()
}
