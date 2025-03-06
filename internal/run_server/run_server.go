package runserver

import (
	"net/http"
	"sync"
)

func RunServer(proxy http.HandlerFunc, port int, readTimeout int, writeTimeout int, idleTimeout int, sslCertificate string, sslKey string, wg *sync.WaitGroup) {
	if sslCertificate == "" && sslKey == "" {
		RunHttp(proxy, port, readTimeout, writeTimeout, idleTimeout, wg)
	} else {
		RunHttps(proxy, port, readTimeout, writeTimeout, idleTimeout, sslCertificate, sslKey, wg)
	}
}
