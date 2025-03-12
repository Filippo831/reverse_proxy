package websocket_handler 

import (
	"log"
	"net/http"

	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

func Handle_websocket(w http.ResponseWriter, r *http.Request, sslToClient bool, sslToServer bool) {
    if sslToServer{
        r.URL.Scheme = "wss"
    } else {
        r.URL.Scheme = "ws"
    }

	log.Printf("request to websocket protocol\n")
	conn_to_server, _, err_to_server := websocket.DefaultDialer.Dial(r.URL.String(), nil)

    if sslToClient{
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
