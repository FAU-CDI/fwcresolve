package main

import (
	"flag"
	"log"
	"net/http"

	resolveproxy "github.com/tkw1536/resolve-proxy"
)

func main() {
	p := resolveproxy.Proxy(map[string]string{
		"^https://(.*)\\.wisski\\.agfd\\.fau\\.de/": "https://$1.wisski.data.fau.de",
		"^https://(.*)\\.wisski\\.data\\.fau\\.de/": "https://$1.wisski.data.fau.de",
		"^http://(.*)\\.wisski\\.agfd\\.fau\\.de/": "https://$1.wisski.data.fau.de",
		"^http://(.*)\\.wisski\\.data\\.fau\\.de/": "https://$1.wisski.data.fau.de",
	})
	log.Printf("Listening on %s", listenAddress)
	log.Fatal(http.ListenAndServe(listenAddress, p))
}

var listenAddress string = "0.0.0.0:8080"

func init() {
	flag.StringVar(&listenAddress, "listen", listenAddress, "Address to listen on")
}
