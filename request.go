package awaitws

type message struct {
	ID       int         `json:"id"`
	Request  bool        `json:"request,omitempty"`
	Response bool        `json:"response,omitempty"`
	Message  interface{} `json:"message"`
}

// Request sends a request via the WebSocket connection and returns
// an AwaitedResponse channel which can be awaited upon.
func (client *Client) Request(msg interface{}) (wait chan AwaitedResponse, err error) {
	// Create a channel to be awaited upon.
	wait = make(chan AwaitedResponse)

	// Add the channel to the waiting response map.
	id := client.newWaitingResponse(wait)

	// Write the request the WebSocket connection.
	err = client.conn.WriteJSON(message{
		ID:      id,
		Request: true,
		Message: msg,
	})

	return
}

// IsRequest returns whether a message is a request or not.
func (client *Client) IsRequest(msg map[string]interface{}) bool {
	// Get the request as an interface.
	requestInterface, ok := msg["request"]
	if !ok {
		return false
	}

	// Assert that the request interface is a bool.
	request, ok := requestInterface.(bool)
	if !ok {
		return false
	}

	return request
}
