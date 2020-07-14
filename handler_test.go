package awaitws

import (
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func TestClient_HandleResponse(t *testing.T) {
	client := NewClient(nil)

	wait := make(chan AwaitedResponse)
	client.waitingResponses[0] = wait

	msg := map[string]interface{}{
		"id":       float64(0),
		"response": true,
		"message":  "msgContent",
	}
	go client.HandleResponse(msg)

	var response AwaitedResponse

	select {
	case response = <-wait:
		break

	case <-time.After(testTimeout):
		t.Error("Timeout")
		return
	}

	assert.Equal(t, msg["message"], response.Message)
	assert.Nil(t, response.Err)

	client.waitingResponseMutex.Lock()
	defer client.waitingResponseMutex.Unlock()

	_, ok := client.waitingResponses[0]
	assert.False(t, ok)
}
