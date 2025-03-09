package main

import (
	"log"
	"os"

	"github.com/Filippo831/reverse_proxy/internal/reverse_proxy"
)

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
