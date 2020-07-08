package awaitws

import (
	"github.com/gorilla/websocket"
	"sync"
)

type Client struct {
	conn                 *websocket.Conn
	nextID               int
	waitingResponses     map[int]waitingResponse
	waitingResponseMutex sync.Mutex
}

func NewClient(conn *websocket.Conn) Client {
	return Client{
		conn:                 conn,
		nextID:               0,
		waitingResponses:     make(map[int]waitingResponse),
		waitingResponseMutex: sync.Mutex{},
	}
}
