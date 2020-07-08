package main

import (
	"flag"
	awaitws "github.com/VolticFroogo/await-ws-go"
	"github.com/gorilla/websocket"
	"log"
	"net/url"
)

var (
	addr = flag.String("addr", ":8080", "http service address")
)

func main() {
	flag.Parse()

	u := url.URL{Scheme: "ws", Host: *addr, Path: "/ws"}
	c, _, err := websocket.DefaultDialer.Dial(u.String(), nil)
	if err != nil {
		log.Fatal(err)
	}

	defer c.Close()

	awaitClient := awaitws.NewClient(c)

	// Create a new goroutine for the initial message.
	go func() {
		// Send a message to the WS server.
		wait, err := awaitClient.Request("client")
		if err != nil {
			log.Print(err)
			return
		}

		// Wait for a response.
		response := <-wait

		// Print the response message.
		// Expected output: "Hello client"
		log.Print(response.Message)
	}()

	for {
		var msg map[string]interface{}
		err = c.ReadJSON(&msg)
		if err != nil {
			if websocket.IsCloseError(err) || websocket.IsUnexpectedCloseError(err) {
				break
			}

			log.Print(err)
			continue
		}

		if awaitClient.HandleResponse(msg) {
			continue
		}

		if awaitClient.IsRequest(msg) {
			err = awaitClient.Respond(msg, "Hello "+msg["message"].(string))
			if err != nil {
				log.Print(err)
				continue
			}
		}
	}
}
