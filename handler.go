package awaitws

// HandleResponse will check if a message is a response, if the message is a response,
// it will pass the message to the original request where it is being awaited upon.
// Returns whether or not the message is a response.
func (client *Client) HandleResponse(msg map[string]interface{}) (handled bool) {
	// If the message is not a response, return.
	if responseInterface, ok := msg["response"]; ok {
		response, ok := responseInterface.(bool)
		if !ok || !response {
			return
		}
	} else {
		return
	}

	// If this point has been reached, the message must be a response;
	// always return true after this point.
	handled = true

	// Get the ID as an interface.
	idInterface, ok := msg["id"]
	if !ok {
		return
	}

	// Assert that the ID interface is a float64.
	idFloat, ok := idInterface.(float64)
	if !ok {
		return
	}

	// Convert the float64 to an int.
	id := int(idFloat)

	// Prevent write collisions.
	client.waitingResponseMutex.Lock()

	// Make a copy of the waitingResponse.
	wait := client.waitingResponses[id]

	// Remove the waitingResponse from the map.
	delete(client.waitingResponses, id)

	client.waitingResponseMutex.Unlock()

	// Send the message to the relevant channel.
	wait <- AwaitedResponse{
		Message: msg["message"],
		Err:     nil,
	}

	return
}
