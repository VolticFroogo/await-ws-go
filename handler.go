package awaitws

func (client *Client) HandleResponse(res map[string]interface{}) (handled bool) {
	if responseInterface, ok := res["response"]; ok {
		response, ok := responseInterface.(bool)
		if !ok || !response {
			return
		}
	} else {
		return
	}

	handled = true
	idInterface, ok := res["id"]
	if !ok {
		return
	}

	idFloat, ok := idInterface.(float64)
	if !ok {
		return
	}

	id := int(idFloat)

	for _, waitingRes := range client.waitingResponses {
		if waitingRes.ID != id {
			continue
		}

		waitingRes.Chan <- AwaitedResponse{
			Message: res["message"],
			Err:     nil,
		}
	}

	delete(client.waitingResponses, id)

	return
}
