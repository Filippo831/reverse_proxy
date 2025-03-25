package reverseproxy

import (
	"fmt"
	"log"
	"net/http"
	"net/url"
	"slices"
	"strings"
	"sync"

	"github.com/quic-go/quic-go/http3"

	"github.com/Filippo831/reverse_proxy/internal/http_handler"
	"github.com/Filippo831/reverse_proxy/internal/read_configuration"
	"github.com/Filippo831/reverse_proxy/internal/run_server"
	"github.com/Filippo831/reverse_proxy/internal/websocket_handler"
)

func RunReverseProxy(conf_path string) error {
	log.Printf("starting reverse proxy\n")

	var wg sync.WaitGroup

	err := readconfiguration.ReadConfiguration(conf_path)
	if err != nil {
		log.Print(err)
		return err
	}

	for _, server := range readconfiguration.Conf.Http {

		proxy := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			redirectURL, _ := url.Parse("https://127.0.0.1:8089")
			found := false

			// get the host name from the request and based on the subdomain redirect to the right url
			domain := r.Host
			domain = strings.Split(domain, ":")[0]

			for _, location := range server.Location {
				if domain == location.Domain {
					redirectURL, _ = url.Parse(location.To)
					found = true
				}
			}

			// message to send if the subdomain requested is not define into location list
			if !found {
				fmt.Fprintf(w, "subdomain not found")
				return
			}

			redirect := http_handler.HttpRedirect(redirectURL, r)

			// check if the client wants to change the connection to a websocket
			if r.Header.Get("Upgrade") == "websocket" && slices.Contains(r.Header.Values("Connection"), "Upgrade") {
				websocket_handler.Handle_websocket(w, redirect, server.SslToClient, redirectURL.Scheme == "https")
			} else {
				http_handler.HttpHandler(w, redirect, server)
			}

		})

		// every server adds 1 to this counter to keep track of how many go routine are running
		wg.Add(1)

		// run the server in a goroutine
		errs := make(chan error, 1)
		go func() {
			log.Printf("running server under domain %s and port %d", server.ServerName, server.Port)

			if server.Http3Active {
				http3Server := &http3.Server{Addr: fmt.Sprintf(":%d", server.Port), Handler: proxy}
				errs <- http3Server.ListenAndServeTLS(server.SslCertificate, server.SslCertificateKey)
			} else {
				errs <- runserver.RunHttp2Server(proxy, server.Port, 5, 10, 60, server.SslCertificate, server.SslCertificateKey, server.SslToClient, &wg)
			}
			if err := <-errs; err != nil {
				log.Fatal(err)
			}
		}()
	}

	wg.Wait()
	return nil

}
