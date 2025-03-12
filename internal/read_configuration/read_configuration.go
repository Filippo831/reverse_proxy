package readconfiguration

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"os"
	"strings"
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


// TODO: make actual check instead of raising only an error to understand why there is an error
// TODO: create test for this
func ReadConfiguration(filePath string) error {

	jsonFile, err := os.Open(filePath)

	if err != nil {
		log.Fatal(err)
        return err
	}

	byteValue, err := io.ReadAll(jsonFile)

	if err != nil {
		log.Fatal(err)
        return err
	}

    newConf := Configuration{}

	readingJsonErr := json.Unmarshal(byteValue, &newConf)

    Conf = newConf

	if readingJsonErr != nil {
		log.Fatal(err)
        return err
	}
    err = checks(&Conf) 
    
    if err != nil {
        log.Fatal(err)
        return err
    }

	return nil
}

func checks(conf *Configuration) error {
	err := checkKeys(conf)
	if err != nil {
		return err
	}

	err = checkDomain(conf)
	if err != nil {
		return err
	}
	return nil
}

func checkKeys(conf *Configuration) error {
	for key, server := range conf.Http {
		if server.SslToClient && (server.SslCertificate == "" || server.SslCertificateKey == "") {
			log.Printf("missing ssl parameter(s) in server %d", key)
			return errors.New(fmt.Sprintf("missing ssl parameter(s) in server %d", key))

		}
		if !server.SslToClient && (server.SslCertificate != "" || server.SslCertificateKey != "") {
			log.Printf("ssl parameters set even if ssl is selected to false in server %d", key)
			return errors.New(fmt.Sprintf("ssl parameters set even if ssl is selected to false in server %d", key))
		}
	}
	return nil
}

func checkDomain(conf *Configuration) error {
	for serverKey, server := range conf.Http {
        subdomains := make(map[string]bool)
		for _, location := range server.Location {
            _, ok := subdomains[location.Domain]
            if ok {
                return errors.New(fmt.Sprintf("conflicting domain in server %d: %s\n", serverKey, location.Domain))
            } else {
                subdomains[location.Domain] = true
            }

			locationDomainArray := strings.Split(location.Domain, ".")
			locationDomain := locationDomainArray[len(locationDomainArray)-1]

			if locationDomain != server.ServerName {
				return errors.New(fmt.Sprintf("server number %d domain: %s\nlocation domain: %s\n", serverKey, server.ServerName, locationDomain))
			}
		}
	}
	return nil
}
