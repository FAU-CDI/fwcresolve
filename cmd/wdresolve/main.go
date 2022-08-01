package main

import (
	"flag"
	"log"
	"net/http"
	"os"

	"github.com/tkw1536/wdresolve"
	"github.com/tkw1536/wdresolve/resolvers"
)

func main() {
	var p wdresolve.ResolveHandler

	fallback := &resolvers.Regexp{
		Data: map[string]string{
			"^https://(.*)\\.wisski\\.agfd\\.fau\\.de/": "https://$1.wisski.data.fau.de",
			"^https://(.*)\\.wisski\\.data\\.fau\\.de/": "https://$1.wisski.data.fau.de",
			"^http://(.*)\\.wisski\\.agfd\\.fau\\.de/":  "https://$1.wisski.data.fau.de",
			"^http://(.*)\\.wisski\\.data\\.fau\\.de/":  "https://$1.wisski.data.fau.de",
		},
	}

	if prefixFile != "" {
		p.Resolver = resolvers.InOrder{
			func() resolvers.Prefix {
				fs, err := os.Open(prefixFile)
				log.Println("loading prefixes from ", prefixFile)
				if err != nil {
					log.Fatal(err)
				}
				defer fs.Close()

				prefixes, err := resolvers.ReadPrefixes(fs)
				if err != nil {
					log.Fatal(err)
				}

				return prefixes
			}(),
			fallback,
		}
	} else {
		p.Resolver = fallback
	}

	log.Printf("Listening on %s", listenAddress)
	log.Fatal(http.ListenAndServe(listenAddress, p))
}

var listenAddress string = "0.0.0.0:8080"
var prefixFile string

func init() {
	defer flag.Parse()

	flag.StringVar(&listenAddress, "listen", listenAddress, "Address to listen on")
	flag.StringVar(&prefixFile, "prefix", prefixFile, "Prefix file to read")
}
