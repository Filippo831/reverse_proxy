package runserver

import (
	"fmt"
	"log"
	"net/http"
	"sync"
	"time"
)

func RunHttp(proxy http.HandlerFunc, port int, readTimeout int, writeTimeout int, idleTimeout int, wg *sync.WaitGroup) {
    // when this block ends wg will decrease by 1. This means that the server ended its work and cann notify that it's not working anymore
	defer wg.Done()

	server := &http.Server{Addr: fmt.Sprintf(":%d", port), ReadTimeout: time.Duration(readTimeout) * time.Second, WriteTimeout: time.Duration(writeTimeout) * time.Second, IdleTimeout: time.Duration(idleTimeout) * time.Second, Handler: proxy}

	if err := server.ListenAndServe(); err != nil {
		log.Fatal(err)
	}
}
