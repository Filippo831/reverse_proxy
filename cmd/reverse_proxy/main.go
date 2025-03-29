package main

import (
	"log"
	"os"

	"github.com/Filippo831/reverse_proxy/internal/reverse_proxy"
)

/*
    run the reverse proxy and attach a file for logging

    LOG
    when creating docker image make the file writeable connected to a system file through docker configuration
    the file is in append mode

    REVERSE PROXY
    call a function that runs the reverse proxy following a configuration
    the configuration will also be connected to a system file through docker configuration

*/

func main() {
	f, err := os.OpenFile("test_log_file.log", os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		log.Fatalf("error opening file: %v", err)
	}
	defer f.Close()

	log.SetOutput(f)
	log.Println("init log")

	reverseproxy.RunReverseProxy("configuration.json")
}
