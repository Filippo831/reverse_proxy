package runserver

import (
	"errors"
	"log"
	"net/http"
	"sync"
)

func RunServer(proxy http.HandlerFunc, port int, readTimeout int, writeTimeout int, idleTimeout int, sslCertificate string, sslKey string, sslActivate bool, wg *sync.WaitGroup) error {

	if sslCertificate == "" && sslKey == "" && !sslActivate {
		RunHttp(proxy, port, readTimeout, writeTimeout, idleTimeout, wg)
	} else if sslCertificate != "" && sslKey != "" && sslActivate {
		RunHttps(proxy, port, readTimeout, writeTimeout, idleTimeout, sslCertificate, sslKey, wg)
	}
	return nil
}
