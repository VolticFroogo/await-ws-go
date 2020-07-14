package awaitws

// Respond will respond to a request with a specified output.
func (client *Client) Respond(in map[string]interface{}, out interface{}) (err error) {
	// Get the ID as an interface.
	idInterface, ok := in["id"]
	if !ok {
		err = ErrNoID
		return
	}

	// Assert that the ID interface is a float64.
	idFloat, ok := idInterface.(float64)
	if !ok {
		err = ErrBadID
		return
	}

	// Convert the float64 to an int.
	id := int(idFloat)

	// Write the message as JSON via the WS connection including
	// the request's ID and the output message.
	return client.conn.WriteJSON(message{
		ID:       id,
		Response: true,
		Message:  out,
	})
}
