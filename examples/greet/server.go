package main

import (
	"flag"
	"github.com/VolticFroogo/await-ws-go"
	"log"
	"net/http"

	"github.com/gorilla/websocket"
)

var (
	upgrader = websocket.Upgrader{
		CheckOrigin: func(r *http.Request) bool {
			return true
		},
	}
	addr = flag.String("addr", ":8080", "http service address")
)

func ws(w http.ResponseWriter, r *http.Request) {
	// Upgrade the connection to a WebSocket.
	c, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Printf("Upgrading connection error: %s", err)
		return
	}

	// Safely close the WebSocket when this function dies.
	defer c.Close()

	// Create a new awaitws client with the connection.
	awaitClient := awaitws.NewClient(c)

	// Create a new goroutine for the initial message.
	go func() {
		// Send a message to the WS client.
		wait, err := awaitClient.Request("server")
		if err != nil {
			log.Print(err)
			return
		}

		// Wait for a response.
		response := <-wait

		// Print the response message.
		// Expected output: "Hello server"
		log.Print(response.Message)
	}()

	for {
		// Decode the incoming message as a mapstructure.
		var msg map[string]interface{}
		err = c.ReadJSON(&msg)
		if err != nil {
			if websocket.IsCloseError(err) || websocket.IsUnexpectedCloseError(err) {
				break
			}

			log.Print(err)
			continue
		}

		// If the message is a response, await-ws will handle it, just continue.
		if awaitClient.HandleResponse(msg) {
			continue
		}

		// Verify that the message is a request.
		if awaitClient.IsRequest(msg) {
			// Respond to the request prepending "Hello " to their message.
			// In a lot of scenarios, you'll have complex types here, and not just a string as the message.
			// To decode a complex type, check out this library: https://github.com/mitchellh/mapstructure
			err = awaitClient.Respond(msg, "Hello "+msg["message"].(string))
			if err != nil {
				log.Print(err)
				continue
			}
		}
	}
}

func main() {
	flag.Parse()

	http.HandleFunc("/ws", ws)

	log.Printf("Listening to %s for incoming HTTP requests", *addr)
	log.Fatal(http.ListenAndServe(*addr, nil))
}
