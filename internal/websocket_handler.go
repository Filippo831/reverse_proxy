package internal

import (
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

func handle_websocket(w http.ResponseWriter, r *http.Request) {
	r.URL.Scheme = "ws"
	fmt.Printf("request to websocket protocol\n")
	conn_to_server, _, err_to_server := websocket.DefaultDialer.Dial(r.URL.String(), nil)
	conn_to_client, err_to_client := upgrader.Upgrade(w, r, nil)

	if err_to_server != nil {
		log.Println(err_to_server)
		return
	}
	if err_to_client != nil {
		log.Println(err_to_client)
		return
	}
	go func() {
		for {
			msgType, msg, err := conn_to_server.ReadMessage()
			fmt.Printf("received message from server\n")
			if err != nil {
				log.Println(err)
				break
			}
			conn_to_client.WriteMessage(msgType, msg)
		}
	}()

	go func() {
		for {
			msgType, msg, err := conn_to_client.ReadMessage()
			fmt.Printf("received message from client\n")
			if err != nil {
				log.Println(err)
				break
			}
			conn_to_server.WriteMessage(msgType, msg)
		}
	}()
}
