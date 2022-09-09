package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"regexp"
	"strings"

	"github.com/FAU-CDI/wdresolve"
	"github.com/FAU-CDI/wdresolve/resolvers"
)

func main() {
	p := wdresolve.ResolveHandler{
		TrustXForwardedProto: trustXForwardedProto,
	}

	fallback := &resolvers.Regexp{
		Data: map[string]string{},
	}

	if domainName != "" {
		fallback.Data[fmt.Sprintf("^https?://(.*)\\.%s", regexp.QuoteMeta(domainName))] = fmt.Sprintf("https://$1.%s", domainName)
		log.Printf("registering default domain %s\n", domainName)
	}
	if legacyDomainName != "" && domainName != "" {
		for _, domain := range strings.Split(legacyDomainName, ",") {
			fallback.Data[fmt.Sprintf("^https?://(.*)\\.%s", regexp.QuoteMeta(domain))] = fmt.Sprintf("https://$1.%s", domainName)
			log.Printf("registering legacy domain %s\n", domain)
		}
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
var trustXForwardedProto bool = false
var prefixFile string = os.Getenv("PREFIX_FILE")
var domainName string = os.Getenv("DEFAULT_DOMAIN")
var legacyDomainName string = os.Getenv("LEGACY_DOMAIN")

func init() {
	defer flag.Parse()

	flag.StringVar(&listenAddress, "listen", listenAddress, "Address to listen on")
	flag.StringVar(&prefixFile, "prefix", prefixFile, "Prefix file to read")
	flag.StringVar(&domainName, "domain", domainName, "Distillery domain name to use as a fallback")
	flag.StringVar(&legacyDomainName, "legacy-domain", legacyDomainName, "Distillery legacy domain name(s) to read")
	flag.BoolVar(&trustXForwardedProto, "trust-proxy", trustXForwardedProto, "Trust any X-Forwarded-Proto Header")
}
