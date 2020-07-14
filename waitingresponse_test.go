package awaitws

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestClient_newWaitingResponse(t *testing.T) {
	t.Run("Matching chan", func(t *testing.T) {
		client := NewClient(nil)

		wait := make(chan AwaitedResponse)
		id := client.newWaitingResponse(wait)

		assert.Equal(t, wait, client.waitingResponses[id])
	})

	t.Run("Incrementing ID", func(t *testing.T) {
		client := NewClient(nil)

		firstID := client.newWaitingResponse(make(chan AwaitedResponse))
		secondID := client.newWaitingResponse(make(chan AwaitedResponse))

		assert.Equal(t, 0, firstID)
		assert.Equal(t, 1, secondID)
	})
}
