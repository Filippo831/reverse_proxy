package main

import (
	"crypto/tls"
	"fmt"
	"io"
	"log"
	"net"
	"net/http"
	"net/url"
	"strings"
	"time"

    "github.com/Filippo831/reverse_proxy/internal"

	"golang.org/x/net/http2"
)

var PORT int = 8081


func main() {
	fmt.Sprintln("starting reverse proxy at port %d", PORT)

	// jellyfin
	// demoURL, err := url.Parse("http://127.17.0.1:8096")

	// my demo server
	demoURL, err := url.Parse("http://127.1.0.1:8088")

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

		tr := &http.Transport{TLSClientConfig: &tls.Config{InsecureSkipVerify: true}}
		http2.ConfigureTransport(tr)

		client := &http.Client{Transport: tr, Timeout: 10 * time.Second}

		// if websocket enter here
		if r.Header.Get("Upgrade") == "websocket" {
            handle_websocket(w , r)
		} else {
            test_http_handler()
			resp, err := client.Do(r)
			if err != nil {
				w.WriteHeader(http.StatusInternalServerError)
				fmt.Fprint(w, err)
				return
			}
			// resp, err := http.DefaultClient.Do(r)
			for key, values := range resp.Header {
				for _, value := range values {
					// log.Printf("%s : %s", key, value)
					w.Header().Add(key, value)
				}
			}

			done := make(chan bool)
			go func() {
				for {
					select {
					case <-time.Tick(10 * time.Millisecond):
						w.(http.Flusher).Flush()
					case <-done:
						return
					}
				}
			}()

			trailerKeys := []string{}

			for key := range resp.Trailer {
				trailerKeys = append(trailerKeys, key)
			}

			w.Header().Set("Trailer", strings.Join(trailerKeys, ","))

			for key, values := range resp.Trailer {
				for _, value := range values {
					w.Header().Set(key, value)
				}
			}

			/*
			   if the url changed (redirect happened), write the field Location into the
			   response to make the client change the url as well
			*/

			if resp.Request.URL.String() != r.URL.String() {
				w.Header().Add("Location", resp.Request.URL.Path)
				w.WriteHeader(http.StatusSeeOther)
			} else {
				w.WriteHeader(resp.StatusCode)
			}

			io.Copy(w, resp.Body)

			close(done)
		}

	})

	// if err := http.ListenAndServeTLS(":8081", "reverse_proxy.rsa.crt", "reverse_proxy.rsa.key", proxy); err != nil {
	// 	log.Fatal(err)
	// }
	if err := http.ListenAndServe(":8081", proxy); err != nil {
		log.Fatal(err)
	}
}
