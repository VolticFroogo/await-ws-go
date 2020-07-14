package awaitws

import (
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func TestClient_HandleResponse(t *testing.T) {
	t.Run("Valid", func(t *testing.T) {
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
	})

	// NOTE: I use the line:
	// client.waitingResponses = nil
	// to throw an error in case it is written to when it shouldn't be.

	t.Run("No response", func(t *testing.T) {
		client := NewClient(nil)
		client.waitingResponses = nil

		msg := map[string]interface{}{}
		assert.False(t, client.HandleResponse(msg))
	})

	t.Run("Invalid response type", func(t *testing.T) {
		client := NewClient(nil)
		client.waitingResponses = nil

		msg := map[string]interface{}{
			"response": "badType",
		}
		assert.False(t, client.HandleResponse(msg))
	})

	t.Run("Response false", func(t *testing.T) {
		client := NewClient(nil)
		client.waitingResponses = nil

		msg := map[string]interface{}{
			"response": false,
		}
		assert.False(t, client.HandleResponse(msg))
	})

	t.Run("No ID", func(t *testing.T) {
		client := NewClient(nil)
		client.waitingResponses = nil

		msg := map[string]interface{}{
			"response": true,
		}
		assert.True(t, client.HandleResponse(msg))
	})

	t.Run("Invalid ID type", func(t *testing.T) {
		client := NewClient(nil)
		client.waitingResponses = nil

		msg := map[string]interface{}{
			"response": true,
			"id":       "badType",
		}
		assert.True(t, client.HandleResponse(msg))
	})
}
