package awaitws

type awaitMessage struct {
	ID      int         `json:"id"`
	Message interface{} `json:"message"`
}

func (awaitWS *Server) Message(msg interface{}) (wait chan AwaitedResponse, err error) {
	wait = make(chan AwaitedResponse)

	id := awaitWS.newWaitingResponse(waitingResponse{
		Chan: wait,
	})

	err = awaitWS.conn.WriteJSON(awaitMessage{
		ID:      id,
		Message: msg,
	})

	return
}
