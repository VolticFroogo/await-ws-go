package main

import (
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
)

func ws(w http.ResponseWriter, r *http.Request) {
	c, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Printf("Upgrading connection error: %s", err)
		return
	}

	defer c.Close()

	awaitServer := awaitws.NewServer(c)

	go func() {
		wait, _ := awaitServer.Message("Hey!")
		response := <-wait

		log.Print(response)
	}()

	for {
		var res map[string]interface{}
		err = c.ReadJSON(&res)
		if err != nil {
			if websocket.IsCloseError(err) || websocket.IsUnexpectedCloseError(err) {
				break
			}

			log.Print(err)
			continue
		}

		if awaitServer.HandleResponse(res) {
			continue
		}
	}
}

func main() {
	http.HandleFunc("/ws", ws)

	log.Println("Listening to :8080 for incoming HTTP requests")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
