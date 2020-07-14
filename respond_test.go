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

func setupTestClient_Respond() (srv *http.Server, msgChan chan map[string]interface{}) {
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

func cleanupTestClient_Respond(srv *http.Server) {
	err := srv.Close()
	if err != nil {
		log.Panic(err)
		return
	}
}

func TestClient_Respond(t *testing.T) {
	srv, msgChan := setupTestClient_Respond()

	t.Run("Valid", func(t *testing.T) {
		u := url.URL{Scheme: "ws", Host: srv.Addr, Path: "/ws"}
		conn, _, err := websocket.DefaultDialer.Dial(u.String(), nil)
		if err != nil {
			log.Panic(err)
			return
		}

		client := NewClient(conn)

		msgContent := "msgContent"
		err = client.Respond(map[string]interface{}{
			"id": float64(0),
		}, msgContent)
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
		assert.Equal(t, true, msg["response"])
		assert.Equal(t, msgContent, msg["message"])
	})

	t.Run("No ID", func(t *testing.T) {
		client := NewClient(nil)

		msgContent := "msgContent"
		err := client.Respond(map[string]interface{}{}, msgContent)

		assert.Equal(t, ErrNoID, err)
	})

	t.Run("Invalid ID type", func(t *testing.T) {
		client := NewClient(nil)

		msgContent := "msgContent"
		err := client.Respond(map[string]interface{}{
			"id": "badType",
		}, msgContent)

		assert.Equal(t, ErrBadID, err)
	})

	cleanupTestClient_Respond(srv)
}
