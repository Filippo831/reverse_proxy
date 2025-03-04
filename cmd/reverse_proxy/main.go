package main

import (
	"fmt"
	"log"
	"net"
	"net/http"
	"net/url"
	"time"

	"github.com/Filippo831/reverse_proxy/internal/http_handler"
	"github.com/Filippo831/reverse_proxy/internal/websocket_handler"
)

var PORT int = 8081

func run_http(proxy http.HandlerFunc) {
	server := &http.Server{Addr: ":8081", ReadTimeout: 5 * time.Second, WriteTimeout: 10 * time.Second, IdleTimeout: 60 * time.Second, Handler: proxy}

	server.SetKeepAlivesEnabled(false)

	if err := server.ListenAndServe(); err != nil {
		log.Fatal(err)
	}
	// if err := http.ListenAndServe(":8088", nil); err != nil {
	// 	log.Fatal(err)
	// }
}

func run_https(proxy http.HandlerFunc) {
	server := &http.Server{Addr: ":8082", ReadTimeout: 5 * time.Second, WriteTimeout: 10 * time.Second, IdleTimeout: 60 * time.Second, Handler: proxy}
	if err := server.ListenAndServeTLS("reverse_proxy.com+3.pem", "reverse_proxy.com+3-key.pem"); err != nil {
		log.Fatal(err)
	}
	// if err := http.ListenAndServeTLS(":8089", , nil); err != nil {
	// 	log.Fatal(err)
	// }
}

func main() {
	fmt.Sprintln("starting reverse proxy at port %d", PORT)

	// jellyfin
	// demoURL, err := url.Parse("http://127.17.0.1:8096")

	// my demo server
	demoURL, err := url.Parse("https://127.0.0.1:8089")

	if err != nil {
		log.Fatal(err)
	}

	proxy := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Println("incoming request")
		r.Host = demoURL.Host
		r.URL.Host = demoURL.Host
		r.URL.Scheme = demoURL.Scheme
		r.RequestURI = ""

		client_host, _, _ := net.SplitHostPort(r.RemoteAddr)
		r.Header.Set("X-Forwarded-For", client_host)

		// if websocket enter here
		if r.Header.Get("Upgrade") == "websocket" {
			websocket_handler.Handle_websocket(w, r)
		} else {
			http_handler.Http_handler(w, r)
		}

	})

	go run_http(proxy)
	run_https(proxy)
}
