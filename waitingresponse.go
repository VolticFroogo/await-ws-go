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
	// Use a mutex to prevent write collisions.
	client.waitingResponseMutex.Lock()
	defer client.waitingResponseMutex.Unlock()

	// Get the next ID and increment it.
	id := client.nextID
	client.nextID++

	// Set the ID to the response and add the response to the map.
	response.ID = id
	client.waitingResponses[id] = response

	return id
}
