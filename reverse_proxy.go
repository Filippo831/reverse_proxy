package main

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
)

var PORT int = 8081

func main() {
	fmt.Sprintln("starting reverse proxy at port %d", PORT)

    demoURL, err := url.Parse("http://127.17.0.1:8080")
    if err != nil {
        log.Fatal(err)
    }

    proxy := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        fmt.Println("incoming request")
        r.Host = demoURL.Host
        r.URL.Host = demoURL.Host
        r.URL.Scheme = demoURL.Scheme
        r.RequestURI = ""

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

        io.Copy(w, resp.Body)

    })

	if err := http.ListenAndServe(":8081", proxy); err != nil {
		log.Fatal(err)
	}
}
