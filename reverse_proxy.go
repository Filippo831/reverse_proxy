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

	"golang.org/x/net/http2"
)

var PORT int = 8081

func main() {
	fmt.Sprintln("starting reverse proxy at port %d", PORT)

	demoURL, err := url.Parse("https://127.17.0.1:8096")
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

		resp, err := http.DefaultClient.Do(r)

		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			fmt.Fprint(w, err)
			return
		}

		for key, values := range resp.Header {
			for _, value := range values {
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

		w.WriteHeader(resp.StatusCode)
		io.Copy(w, resp.Body)

		close(done)
	})

	if err := http.ListenAndServeTLS(":8081", "reverse_proxy.rsa.crt", "reverse_proxy.rsa.key", proxy); err != nil {
		log.Fatal(err)
	}
}
