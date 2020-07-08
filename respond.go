package awaitws

func (client *Client) Respond(in map[string]interface{}, out interface{}) (err error) {
	idInterface, ok := in["id"]
	if !ok {
		err = ErrNoID
		return
	}

	idFloat, ok := idInterface.(float64)
	if !ok {
		err = ErrorBadID
		return
	}

	id := int(idFloat)

	return client.conn.WriteJSON(message{
		ID:       id,
		Response: true,
		Message:  out,
	})
}
