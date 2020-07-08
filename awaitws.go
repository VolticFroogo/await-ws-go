package awaitws

import (
	"github.com/gorilla/websocket"
	"sync"
)

type Server struct {
	conn                 *websocket.Conn
	nextID               int
	waitingResponses     map[int]waitingResponse
	waitingResponseMutex sync.Mutex
}

func NewServer(conn *websocket.Conn) Server {
	return Server{
		conn:                 conn,
		nextID:               0,
		waitingResponses:     make(map[int]waitingResponse),
		waitingResponseMutex: sync.Mutex{},
	}
}
