package reverseproxy

import (
	"crypto/tls"
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

	/*
	   create wait group because you can run multiple server in go routine.
	   without waitgroup you will have a goroutine for each server but the main execution will end a so the routines
	   with this you wait for all goroutine to end

	   this is only a counter that increments when a routine is run and decrement when it ends
	*/
	var wg sync.WaitGroup

	/*
	   given the configuration file run a function that create a go object out of it
	   the configuration object is a global variable in the same file read from other files
	*/

	err := readconfiguration.ReadConfiguration(conf_path)
	if err != nil {
		log.Print(err)
		return err
	}

	/*
	   for each server defined in the configuration run this function into the loop

	   LOOP TASK
	   - given the subdomain get the server to whom send the request
	   - build the request for the backend server
	   - check the header and decide whether use http connection or websocket connection
	*/
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

			// get request to redirect to backend server
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

			/*
			   TODO: for now this is a try to make http3 work
			*/
			if server.Http3Active {
				cert, err := tls.LoadX509KeyPair(server.SslCertificate, server.SslCertificateKey)
				errs <- err

				http3Server := http3.Server{Addr: ":8083", Handler: proxy, TLSConfig: &tls.Config{Certificates: []tls.Certificate{cert}, NextProtos: []string{"h3"}}}
				errs <- http3Server.ListenAndServe()

			} else {
				/*
				   run http2 server which will also work with http1
				   it will try to use the newest protocol
				*/
				errs <- runserver.RunHttp2Server(proxy, server.Port, 5, 10, 60, server.SslCertificate, server.SslCertificateKey, server.SslToClient, &wg)
			}
			if err := <-errs; err != nil {
				log.Fatal(err)
				return
			}
		}()
	}

    // wait until every server stop running before ending the execution
	wg.Wait()
	return nil

}
