package awaitws

type AwaitedResponse struct {
	Message interface{}
	Err     error
}

type waitingResponse struct {
	ID   int
	Chan chan AwaitedResponse
}

func (client *Client) newWaitingResponse(response waitingResponse) int {
	client.waitingResponseMutex.Lock()
	defer client.waitingResponseMutex.Unlock()

	id := client.nextID
	client.nextID++

	response.ID = id
	client.waitingResponses[id] = response

	return id
}
