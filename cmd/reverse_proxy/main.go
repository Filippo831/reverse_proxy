package main

import (
	"fmt"
	"log"
	"net/http"
	"net/url"
    "sync"

	"github.com/Filippo831/reverse_proxy/internal/http_handler"
	"github.com/Filippo831/reverse_proxy/internal/read_configuration"
	"github.com/Filippo831/reverse_proxy/internal/run_server"
	"github.com/Filippo831/reverse_proxy/internal/websocket_handler"
)

func main() {
	fmt.Sprintln("starting reverse proxy")

    var wg sync.WaitGroup

	configuration := readconfiguration.ReadConfiguration("configuration.json")


	for _, server := range configuration.Http {
        fmt.Printf(server.SslCertificate)
        fmt.Printf("\n")

        fmt.Printf("loop\n")
		// redirectURL, err := url.Parse("https://127.0.0.1:8089")
		redirectURL, err := url.Parse(server.Location.To)

		if err != nil {
			log.Fatal(err)
		}

		proxy := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			r = http_handler.HttpRedirect(redirectURL, r)

			// if websocket enter here
			if r.Header.Get("Upgrade") == "websocket" {
				websocket_handler.Handle_websocket(w, r)
			} else {
				http_handler.HttpHandler(w, r)
			}

		})

        wg.Add(1)

        go runserver.RunServer(proxy, server.Port, 5, 10, 60, server.SslCertificate, server.SslCertificateKey, &wg)
	}

    wg.Wait()

	// runserver.RunHttps(proxy, 8082, 5, 10, 60, "reverse_proxy.com+3.pem", "reverse_proxy.com+3-key.pem")
}
