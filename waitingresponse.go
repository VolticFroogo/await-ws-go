package awaitws

type AwaitedResponse struct {
	Message interface{}
	Err     error
}

type waitingResponse struct {
	ID   int
	Chan chan AwaitedResponse
}

func (server *Server) newWaitingResponse(response waitingResponse) int {
	server.waitingResponseMutex.Lock()
	defer server.waitingResponseMutex.Unlock()

	id := server.nextID
	server.nextID++

	response.ID = id
	server.waitingResponses[id] = response

	return id
}
