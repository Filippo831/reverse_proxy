package reverseproxy

import (
	"log"
	"net/http"
	"net/url"
	"slices"
	"strings"
	"sync"

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
        // TODO: make a base url for 404-ish error. This will be the base url that show up if no matches is found under location list
		redirectURL, _ := url.Parse("https://127.0.0.1:8089")
		proxy := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

            // get the host name from the request and based on the subdomain redirect to the right url

			domain := r.Host
			domain = strings.Split(domain, ":")[0]


			for _, location := range server.Location {
				if domain == location.Domain {
					redirectURL, _ = url.Parse(location.To)
				}
			}

			r = http_handler.HttpRedirect(redirectURL, r)

			// check if the client wants to change the connection to a websocket
			if r.Header.Get("Upgrade") == "websocket" && slices.Contains(r.Header.Values("Connection"), "Upgrade") {
				websocket_handler.Handle_websocket(w, r, server.SslToClient, redirectURL.Scheme == "https")
			} else {
				http_handler.HttpHandler(w, r, server.MaxRedirect)
			}

		})

        // every server adds 1 to this counter to keep track of how many go routine are running
		wg.Add(1)

        // run the server in a goroutine
		errs := make(chan error, 1)
		go func() {
			errs <- runserver.RunServer(proxy, server.Port, 5, 10, 60, server.SslCertificate, server.SslCertificateKey, server.SslToClient, &wg)
		}()
		if err := <-errs; err != nil {
			return err
		}
	}

	wg.Wait()
	return nil

}
