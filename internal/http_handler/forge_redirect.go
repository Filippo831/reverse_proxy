package http_handler

import (
	"net"
	"net/http"
	"net/url"
)

func HttpRedirect(redirectURL *url.URL, r *http.Request) *http.Request {
	r.Host = redirectURL.Host
	r.URL.Host = redirectURL.Host
	r.URL.Scheme = redirectURL.Scheme
	r.RequestURI = ""

	client_host, _, _ := net.SplitHostPort(r.RemoteAddr)
	r.Header.Set("X-Forwarded-For", client_host)

	return r
}
