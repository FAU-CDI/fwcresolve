package main

import (
	"flag"
	"log"
	"net/http"

	"github.com/tkw1536/fwcresolve"
	"github.com/tkw1536/fwcresolve/resolvers"
)

func main() {
	r := &resolvers.Regexp{
		Data: map[string]string{
			"^https://(.*)\\.wisski\\.agfd\\.fau\\.de/": "https://$1.wisski.data.fau.de",
			"^https://(.*)\\.wisski\\.data\\.fau\\.de/": "https://$1.wisski.data.fau.de",
			"^http://(.*)\\.wisski\\.agfd\\.fau\\.de/":  "https://$1.wisski.data.fau.de",
			"^http://(.*)\\.wisski\\.data\\.fau\\.de/":  "https://$1.wisski.data.fau.de",
		},
	}
	p := fwcresolve.ResolveHandler{
		Resolver: r,
	}
	log.Printf("Listening on %s", listenAddress)
	log.Fatal(http.ListenAndServe(listenAddress, p))
}

var listenAddress string = "0.0.0.0:8080"

func init() {
	flag.StringVar(&listenAddress, "listen", listenAddress, "Address to listen on")
}
