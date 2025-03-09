package test

import (
	"fmt"
	"os"
	"testing"

	reverseproxy "github.com/Filippo831/reverse_proxy/internal/reverse_proxy"
)

func TestMissingSSLParameter(t *testing.T) {
	configuration := ` { "servers": [ { "port": 8082, "server_name": "127.0.0.1", "ssl_to_client": true, "ssl_certificate": "../reverse_proxy.com+3.pem", "ssl_certificate_key": "","max_redirect": 10 ,"Location": { "path": "/", "to": "http://127.0.0.1:8080" } }] } `
	os.WriteFile("configuration_test.json", []byte(configuration), 0644)

	errs := make(chan error, 1)
	go func() {
		errs <- reverseproxy.RunReverseProxy("configuration_test.json")
	}()
	err := <-errs
	if fmt.Sprint(err) != "missing ssl parameter(s)" {
		t.Error(err)
	}
}

func TestKeyCorrectWrongBool(t *testing.T) {
	configuration := ` { "servers": [ { "port": 8082, "server_name": "127.0.0.1", "ssl_to_client": false, "ssl_certificate": "../reverse_proxy.com+3.pem", "ssl_certificate_key": "../reverse_proxy.com+3-key.pem","max_redirect": 10 ,"Location": { "path": "/", "to": "http://127.0.0.1:8080" } }] } `
	os.WriteFile("configuration_test.json", []byte(configuration), 0644)

	errs := make(chan error, 1)
	go func() {
		errs <- reverseproxy.RunReverseProxy("configuration_test.json")
	}()
	err := <-errs
	if fmt.Sprint(err) != "ssl parameters set even if ssl is selected to false" {
		t.Error(err)
	}
}
