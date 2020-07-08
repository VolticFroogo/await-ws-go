package awaitws

type AwaitedResponse struct {
	Message interface{}
	Err     error
}

func (client *Client) newWaitingResponse(response chan AwaitedResponse) int {
	// Use a mutex to prevent write collisions.
	client.waitingResponseMutex.Lock()
	defer client.waitingResponseMutex.Unlock()

	// Get the next ID and increment it.
	id := client.nextID
	client.nextID++

	// Add the response to the map.
	client.waitingResponses[id] = response

	return id
}
