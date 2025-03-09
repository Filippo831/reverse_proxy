package reverseproxy

import (
	"log"
	"net/http"
	"net/url"
	"slices"
	"sync"

	"github.com/Filippo831/reverse_proxy/internal/http_handler"
	"github.com/Filippo831/reverse_proxy/internal/read_configuration"
	"github.com/Filippo831/reverse_proxy/internal/run_server"
	"github.com/Filippo831/reverse_proxy/internal/websocket_handler"
)

func RunReverseProxy(conf_path string) {
	log.Printf("starting reverse proxy\n")

	var wg sync.WaitGroup

	readconfiguration.ReadConfiguration(conf_path)

	for _, server := range readconfiguration.Conf.Http {
		// redirectURL, err := url.Parse("https://127.0.0.1:8089")
		redirectURL, err := url.Parse(server.Location.To)

		if err != nil {
			log.Fatal(err)
		}

		proxy := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			r = http_handler.HttpRedirect(redirectURL, r)

			// check if the client wants to change the connection to a websocket
			if r.Header.Get("Upgrade") == "websocket" && slices.Contains(r.Header.Values("Connection"), "Upgrade") {
				websocket_handler.Handle_websocket(w, r, server.SslToClient, redirectURL.Scheme == "https")
			} else {
				http_handler.HttpHandler(w, r, server.MaxRedirect)
			}

		})

		wg.Add(1)

		go runserver.RunServer(proxy, server.Port, 5, 10, 60, server.SslCertificate, server.SslCertificateKey, &wg)
	}

	wg.Wait()

}
