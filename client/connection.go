package client

import "sync"

type Connection struct {
}

var (
	connPool = sync.Pool{
		New: func() interface{} {
			return &Connection{}
		},
	}
)

const (
	MAX_CONNECTION_POOL = 10
)
