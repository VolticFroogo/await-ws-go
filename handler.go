package awaitws

func (awaitWS *Server) HandleResponse(res map[string]interface{}) bool {
	if response, ok := res["response"]; ok {
		if !response.(bool) {
			return false
		}
	}

	idFloat, ok := res["id"].(float64)
	if !ok {
		return true
	}

	id := int(idFloat)

	for _, waitingRes := range awaitWS.waitingResponses {
		if waitingRes.ID != id {
			continue
		}

		waitingRes.Chan <- AwaitedResponse{
			Message: res["message"],
			Err:     nil,
		}
	}

	delete(awaitWS.waitingResponses, id)

	return true
}
