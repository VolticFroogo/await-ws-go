package awaitws

type awaitMessage struct {
	ID      int         `json:"id"`
	Message interface{} `json:"message"`
}

func (server *Server) Message(msg interface{}) (wait chan AwaitedResponse, err error) {
	wait = make(chan AwaitedResponse)

	id := server.newWaitingResponse(waitingResponse{
		Chan: wait,
	})

	err = server.conn.WriteJSON(awaitMessage{
		ID:      id,
		Message: msg,
	})

	return
}
