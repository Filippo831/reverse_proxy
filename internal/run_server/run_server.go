package runserver

import (
	"errors"
	"log"
	"net/http"
	"sync"
)

func RunServer(proxy http.HandlerFunc, port int, readTimeout int, writeTimeout int, idleTimeout int, sslCertificate string, sslKey string, sslActivate bool, wg *sync.WaitGroup) error {
	if sslActivate && (sslCertificate == "" || sslKey == "") {
		log.Printf("missing ssl parameter(s)")
		return errors.New("missing ssl parameter(s)")

	}
	if !sslActivate && (sslCertificate != "" || sslKey != "") {
		log.Printf("ssl parameters set even if ssl is selected to false")
		return errors.New("ssl parameters set even if ssl is selected to false")
	}

	if sslCertificate == "" && sslKey == "" && !sslActivate {
		RunHttp(proxy, port, readTimeout, writeTimeout, idleTimeout, wg)
	} else if sslCertificate != "" && sslKey != "" && sslActivate {
		RunHttps(proxy, port, readTimeout, writeTimeout, idleTimeout, sslCertificate, sslKey, wg)
	}
	return nil
}
