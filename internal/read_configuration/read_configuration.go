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
    SslCertificate string `json:"ssl_certificate"`
    SslCertificateKey string `json:"ssl_certificate_key"`
}

type Location struct {
    Path string `json:"path"`
    To string `json:"to"`
}

func ReadConfiguration(filePath string) Configuration {
    jsonFile, err := os.Open(filePath)

    if err != nil {
        log.Fatal(err)
    }

    byteValue, err := io.ReadAll(jsonFile)

    if err != nil {
        log.Fatal(err)
    }

    var configuration Configuration

    readingJsonErr := json.Unmarshal(byteValue, &configuration)

    if readingJsonErr != nil {
        log.Fatal(err)
    }

    
    return configuration
}
