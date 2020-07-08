package awaitws

type message struct {
	ID       int         `json:"id"`
	Request  bool        `json:"request,omitempty"`
	Response bool        `json:"response,omitempty"`
	Message  interface{} `json:"message"`
}

func (client *Client) Request(msg interface{}) (wait chan AwaitedResponse, err error) {
	wait = make(chan AwaitedResponse)

	id := client.newWaitingResponse(waitingResponse{
		Chan: wait,
	})

	err = client.conn.WriteJSON(message{
		ID:      id,
		Request: true,
		Message: msg,
	})

	return
}

func (client *Client) IsRequest(msg map[string]interface{}) bool {
	requestInterface, ok := msg["request"]
	if !ok {
		return false
	}

	request, ok := requestInterface.(bool)
	if !ok {
		return false
	}

	return request
}
