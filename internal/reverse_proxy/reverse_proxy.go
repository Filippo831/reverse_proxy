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

		// run the server in a goroutine
		errs := make(chan error, 1)

		if server.Http3Active {

			/*
				HTTP3 SERVER

				when Http3Active is true, create 2 parallel server on the same port (since 2 different
				protocols are used: UDP for HTTP/3 and TCP/IP for HTTP2). When the first request arrives
				inform the client that http3 is available with the "Alt-Svc" header

				alt-svc h3=":$PORT"; ma=86400

				ma=$SECONDS -> time to keep the information alive, after this time the client will make a http/2 request again

				FIX:
				--- http/3 works only if you craft the request for http/3, the switch does not work ---

			*/
			wg.Add(2)
			log.Printf("running server under domain %s and port %d", server.ServerName, server.Port)

			http3Server := http3.Server{Addr: fmt.Sprintf(":%d", server.Port), Handler: proxy}

			http2Server := http.Server{Addr: fmt.Sprintf(":%d", server.Port), Handler: http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				w.Header().Set("Alt-Svc", fmt.Sprintf(`h3=":%d"; ma=86400`, server.Port))
				proxy.ServeHTTP(w, r)
			})}

			go func() {
				defer wg.Done()
				errs <- http2Server.ListenAndServeTLS(server.SslCertificate, server.SslCertificateKey)
				if err := <-errs; err != nil {
					log.Fatal(err)
					return
				}
			}()
			go func() {
				defer wg.Done()
				errs <- http3Server.ListenAndServeTLS(server.SslCertificate, server.SslCertificateKey)
				if err := <-errs; err != nil {
					log.Fatal(err)
					return
				}
			}()

		} else {
			wg.Add(1)
			/*
			   run http2 server which will also work with http1
			   it will try to use the newest protocol
			*/
			go func() {
				defer wg.Done()
				log.Printf("running server under domain %s and port %d", server.ServerName, server.Port)
				errs <- runserver.RunHttp2Server(proxy, server.Port, 5, 10, 60, server.SslCertificate, server.SslCertificateKey, server.SslToClient, &wg)
				if err := <-errs; err != nil {
					log.Fatal(err)
					return
				}
			}()
		}
	}

	// wait until every server stop running before ending the execution
	wg.Wait()
	return nil

}
