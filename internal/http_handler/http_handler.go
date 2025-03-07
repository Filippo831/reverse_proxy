package http_handler

import (
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"strings"
	"time"

	"golang.org/x/net/http2"
)

func HttpHandler(w http.ResponseWriter, r *http.Request, configurationRedirect int) {
	http2.ConfigureTransport(http.DefaultTransport.(*http.Transport))

	client := &http.Client{Timeout: 10 * time.Second, CheckRedirect: func(req *http.Request, via []*http.Request) error {
		maxRedirect := 10
		if configurationRedirect != 0 {
			maxRedirect = configurationRedirect
		}
		if len(via) >= maxRedirect {
			log.Printf("stopped after %d redirects\n", maxRedirect)
            return errors.New(fmt.Sprintf("stopped after %d redirects\n", maxRedirect))
		}
		return nil
	}}
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
