package test

import (
	"io"
	"log"
	"net/http"
	"os"
	"strings"
	"testing"

	"github.com/Filippo831/reverse_proxy/internal/reverse_proxy"
)


func TestRedirectLessThanThreshold(t *testing.T) {
	configuration := ` { "servers": [ { "port": 8081, "server_name": "127.0.0.1","max_redirect": 10 ,"Location": { "path": "/", "to": "http://127.0.0.1:8080" } }] } `
	os.WriteFile("configuration_test.json", []byte(configuration), 0644)

    go reverseproxy.RunReverseProxy("configuration_test.json")    

	direct_resp, direct_err := http.Get("http://127.0.0.1:8080/redirect/9")

	if direct_err != nil {
		t.Errorf("%s", direct_err)
	}

	proxy_resp, proxy_err := http.Get("http://127.0.0.1:8081/redirect/9")

	if proxy_err != nil {
		t.Errorf("%s", direct_err)
	}
	if proxy_resp.Header.Get("Location") != direct_resp.Header.Get("Location") {
		t.Errorf("redirect Location does not correspond\nproxy: %s\ndirect: %s\n", proxy_resp.Header.Get("Location"), direct_resp.Header.Get("Location"))
	}
    // close(quit)
}

func TestRedirectMoreThanThreshold(t *testing.T) {
	configuration := ` { "servers": [ { "port": 8081, "server_name": "127.0.0.1","max_redirect": 10 ,"Location": { "path": "/", "to": "http://127.0.0.1:8080" } }] } `
	os.WriteFile("configuration_test.json", []byte(configuration), 0644)

    // go reverseproxy.RunReverseProxy("configuration_test.json")

	direct_resp, direct_err := http.Get("http://127.0.0.1:8080/redirect/11")

	if direct_err != nil {
		log.Printf("%s", direct_err)
	}

	proxy_resp, proxy_err := http.Get("http://127.0.0.1:8081/redirect/11")

	if proxy_err != nil {
		log.Printf("%s", direct_err)
	}
	if proxy_resp.Header.Get("Location") == direct_resp.Header.Get("Location") {
		t.Errorf("both location are equals\nproxy: %s\ndirect: %s\n", proxy_resp.Header.Get("Location"), direct_resp.Header.Get("Location"))
	}

    content, err := io.ReadAll(proxy_resp.Body)
    if err != nil {
        t.Error("error while reading the body")
    }
    bodyString := string(content)

    if strings.TrimSpace(bodyString) != `Get "/relative-redirect/1": stopped after 10 redirects` {
		t.Errorf(`expected output: Get "/relative-redirect/1": stopped after 10 redirects\ngot: %s`, bodyString)
	}
    // close(quit)
}
