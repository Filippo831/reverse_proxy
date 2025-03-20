package http_handler

import (
	// "errors"
	"errors"
	"fmt"
	"io"
	"log"

	// "log"
	"net/http"
	"strings"
	"time"

	readconfiguration "github.com/Filippo831/reverse_proxy/internal/read_configuration"
	"golang.org/x/net/http2"
)

func HttpHandler(w http.ResponseWriter, r *http.Request, conf readconfiguration.Server) {
	http2.ConfigureTransport(http.DefaultTransport.(*http.Transport))

	client := &http.Client{Timeout: 10 * time.Second, CheckRedirect: func(req *http.Request, via []*http.Request) error {

		// default value of max redirect
		maxRedirect := 10

		if conf.MaxRedirect != 0 {
			maxRedirect = conf.MaxRedirect
		}

		if len(via) >= maxRedirect {
			log.Printf("stopped after %d redirects\n", maxRedirect)
			return errors.New(fmt.Sprintf("stopped after %d redirects\n", maxRedirect))
		}

		// if one of the intermidiate response send a cookie to set,
		// write it on the final answer to the client, otherwise it get lost inside the redirects
		if req.Response.Header.Get("Set-Cookie") != "" {
			w.Header().Add("Set-Cookie", req.Response.Header.Get("Set-Cookie"))
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
	if conf.ChunkEncoding {
		go func() {
			for {
				select {
                // TODO: define how to make the stream work, if based on time like here or based on blocksize
				case <-time.Tick(time.Duration(conf.ChunkTimeout) * time.Millisecond):
					w.(http.Flusher).Flush()
				case <-done:
					return
				}
			}
		}()
	}

	trailerKeys := []string{}

	for key := range resp.Trailer {
		trailerKeys = append(trailerKeys, key)
	}

	if len(trailerKeys) > 0 {
		w.Header().Set("Trailer", strings.Join(trailerKeys, ","))
	}

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

	w.WriteHeader(resp.StatusCode)

	io.Copy(w, resp.Body)

	close(done)
}
