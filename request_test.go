package awaitws

import (
	"github.com/gorilla/websocket"
	"github.com/stretchr/testify/assert"
	"log"
	"net/http"
	"net/url"
	"testing"
	"time"
)

func setupTestClient_Request() (srv *http.Server, msgChan chan map[string]interface{}) {
	upgrader := websocket.Upgrader{
		CheckOrigin: func(r *http.Request) bool {
			return true
		},
	}

	r := http.NewServeMux()

	msgChan = make(chan map[string]interface{})

	r.HandleFunc("/ws", func(w http.ResponseWriter, r *http.Request) {
		// Upgrade the connection to a WebSocket.
		c, err := upgrader.Upgrade(w, r, nil)
		if err != nil {
			log.Printf("Upgrading connection error: %s", err)
			return
		}

		// Safely close the WebSocket when this function dies.
		defer c.Close()

		for {
			var msg map[string]interface{}
			err = c.ReadJSON(&msg)
			if err != nil && (websocket.IsCloseError(err) || websocket.IsUnexpectedCloseError(err)) {
				break
			}

			msgChan <- msg
		}
	})

	srv = &http.Server{
		Addr:    "localhost:8080",
		Handler: r,
	}

	go srv.ListenAndServe()

	return
}

func cleanupTestClient_Request(srv *http.Server) {
	err := srv.Close()
	if err != nil {
		log.Panic(err)
		return
	}
}

func TestClient_Request(t *testing.T) {
	srv, msgChan := setupTestClient_Request()

	u := url.URL{Scheme: "ws", Host: srv.Addr, Path: "/ws"}
	conn, _, err := websocket.DefaultDialer.Dial(u.String(), nil)
	if err != nil {
		log.Panic(err)
		return
	}

	client := NewClient(conn)

	msgContent := "msgContent"
	wait, err := client.Request(msgContent)
	if err != nil {
		log.Panic(err)
		return
	}

	var msg map[string]interface{}

	select {
	case msg = <-msgChan:
		break

	case <-time.After(testTimeout):
		t.Error("Timeout")
		return
	}

	assert.Equal(t, float64(0), msg["id"])
	assert.Equal(t, true, msg["request"])
	assert.Equal(t, msgContent, msg["message"])
	assert.Equal(t, wait, client.waitingResponses[0])

	wait, err = client.Request(msgContent)
	if err != nil {
		log.Panic(err)
		return
	}

	select {
	case msg = <-msgChan:
		break

	case <-time.After(testTimeout):
		t.Error("Timeout")
		return
	}

	assert.Equal(t, float64(1), msg["id"])
	assert.Equal(t, true, msg["request"])
	assert.Equal(t, msgContent, msg["message"])
	assert.Equal(t, wait, client.waitingResponses[1])

	err = conn.Close()
	if err != nil {
		log.Panic(err)
		return
	}

	cleanupTestClient_Request(srv)
}

func TestClient_IsRequest(t *testing.T) {
	client := NewClient(nil)

	tests := []struct {
		name     string
		arg      map[string]interface{}
		expected bool
	}{
		{
			name: "True",
			arg: map[string]interface{}{
				"request": true,
			},
			expected: true,
		},
		{
			name: "False",
			arg: map[string]interface{}{
				"request": false,
			},
			expected: false,
		},
		{
			name:     "Empty",
			arg:      map[string]interface{}{},
			expected: false,
		},
		{
			name: "Invalid type",
			arg: map[string]interface{}{
				"request": "badType",
			},
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.expected, client.IsRequest(tt.arg))
		})
	}
}
