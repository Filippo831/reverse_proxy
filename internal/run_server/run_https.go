package runserver

import (
	"fmt"
	"log"
	"net/http"
	"sync"
	"time"
)

func RunHttps(proxy http.HandlerFunc, port int, readTimeout int, writeTimeout int, idleTimeout int, sslCertificate string, sslKey string, wg *sync.WaitGroup) {
    defer wg.Done()

	server := &http.Server{Addr: fmt.Sprintf(":%d", port), ReadTimeout: time.Duration(readTimeout) * time.Second, WriteTimeout: time.Duration(writeTimeout) * time.Second, IdleTimeout: time.Duration(idleTimeout) * time.Second, Handler: proxy}

	if err := server.ListenAndServeTLS(sslCertificate, sslKey); err != nil {
		log.Fatal(err)
	}
}
