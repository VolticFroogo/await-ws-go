package awaitws

import (
	"github.com/gorilla/websocket"
	"sync"
)

type Client struct {
	conn                 *websocket.Conn
	nextID               int
	waitingResponses     map[int]chan AwaitedResponse
	waitingResponseMutex sync.Mutex
}

// NewClient creates a ws-client using the default options.
func NewClient(conn *websocket.Conn) Client {
	// Create the default client.
	return Client{
		conn:                 conn,
		nextID:               0,
		waitingResponses:     make(map[int]chan AwaitedResponse),
		waitingResponseMutex: sync.Mutex{},
	}
}
