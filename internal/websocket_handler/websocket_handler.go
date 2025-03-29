package websocket_handler

import (
	"log"
	"net/http"

	"github.com/gorilla/websocket"
)

/*
utility to create a websocket connection out of a http connection
*/
var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

/*
estabalish communication with the client and with the server and pass results
- a server is created to listen to communications from the client
- a client connected to the server to forward reqests

then every time a response or a request comes just forward to the other end
*/

func Handle_websocket(w http.ResponseWriter, r *http.Request, sslToClient bool, sslToServer bool) {
	// if ssl to server enabled use wss otherwise use ws to encrypt or not the connection
	if sslToServer {
		r.URL.Scheme = "wss"
	} else {
		r.URL.Scheme = "ws"
	}

	log.Printf("request to websocket protocol\n")

    // estabalish the communication with the client
	conn_to_server, _, err_to_server := websocket.DefaultDialer.Dial(r.URL.String(), nil)

    // if ssl to client enabled use wss otherwise use ws to encrypt or not the connection
	if sslToClient {
		r.URL.Scheme = "wss"
	} else {
		r.URL.Scheme = "ws"
	}
	conn_to_client, err_to_client := upgrader.Upgrade(w, r, nil)

	if err_to_server != nil {
		log.Println(err_to_server)
		return
	}
	if err_to_client != nil {
		log.Println(err_to_client)
		return
	}

    
    // get a message from the server and send it to the client
	go func() {
		for {
			msgType, msg, err := conn_to_server.ReadMessage()
			if err != nil {
				log.Println(err)
				break
			}
			conn_to_client.WriteMessage(msgType, msg)
		}
	}()

    // get a message from the client and send it to the server
	go func() {
		for {
			msgType, msg, err := conn_to_client.ReadMessage()
			if err != nil {
				log.Println(err)
				break
			}
			conn_to_server.WriteMessage(msgType, msg)
		}
	}()
}
