package http_handler

import (
	"net"
	"net/http"
	"net/url"
)

/*
   build the redirect request to send to backend server
   add X-Forwarded-For to keep the client ip. This is required from some server to keep track of the users
*/

func HttpRedirect(redirectURL *url.URL, r *http.Request) *http.Request {
	r.Host = redirectURL.Host
	r.URL.Host = redirectURL.Host
	r.URL.Scheme = redirectURL.Scheme
	r.RequestURI = ""

	client_host, _, _ := net.SplitHostPort(r.RemoteAddr)
	r.Header.Set("X-Forwarded-For", client_host)

	return r
}
