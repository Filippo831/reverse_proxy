package readconfiguration

import (
	"encoding/json"
	"io"
	"log"
	"os"
)

type Configuration struct {
    Http []Server `json:"servers"`

}
type Server struct {
    Port int `json:"port"`
    ServerName string `json:"server_name"`
    Location Location `json:"location"`
    SslToClient bool `json:"ssl_to_client"`
    SslCertificate string `json:"ssl_certificate"`
    SslCertificateKey string `json:"ssl_certificate_key"`
    MaxRedirect int `json:"max_redirect"`
}

type Location struct {
    Path string `json:"path"`
    To string `json:"to"`
}

var Conf Configuration

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
