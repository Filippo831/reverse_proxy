package readconfiguration

import (
	"encoding/json"
	"io"
	"log"
	"os"
)

type Configuration struct {
    // list of all server, possibly one server for http and one for https connection
    Http []Server `json:"servers"`

}
type Server struct {
    // select the server port
    Port int `json:"port"`

    // server base domain
    ServerName string `json:"server_name"`
    
    // define to which server a subdomain should connect to
    Location []Location `json:"location"`

    // chose if you want to activate ssl connection from reverse proxy to the client
    SslToClient bool `json:"ssl_to_client"`

    // ssl certificate file path
    SslCertificate string `json:"ssl_certificate"`

    // ssl certificate key file path
    SslCertificateKey string `json:"ssl_certificate_key"`

    // highest amount of redirect that the reverse proxy will solve. After this threshold it raise an error
    MaxRedirect int `json:"max_redirect"`
}

type Location struct {
    // subdomain.domain of the desired routing
    Domain string `json:"domain"`

    // server to connect to from this subdomain
    To string `json:"to"`
}

var Conf Configuration

//TODO: make actual check instead of raising only an error to understand why there is an error
//TODO: create test for this
func ReadConfiguration(filePath string) {
    jsonFile, err := os.Open(filePath)

    if err != nil {
        log.Fatal(err)
    }

    byteValue, err := io.ReadAll(jsonFile)

    if err != nil {
        log.Fatal(err)
    }

    readingJsonErr := json.Unmarshal(byteValue, &Conf)

    if readingJsonErr != nil {
        log.Fatal(err)
    }
}
