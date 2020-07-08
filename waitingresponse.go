package awaitws

type AwaitedResponse struct {
	Message interface{}
	Err     error
}

type waitingResponse struct {
	ID   int
	Chan chan AwaitedResponse
}

func (awaitWS *Server) newWaitingResponse(response waitingResponse) int {
	awaitWS.waitingResponseMutex.Lock()
	defer awaitWS.waitingResponseMutex.Unlock()

	id := awaitWS.nextID
	awaitWS.nextID++

	response.ID = id
	awaitWS.waitingResponses[id] = response

	return id
}
