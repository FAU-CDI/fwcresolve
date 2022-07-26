package resolveproxy

import (
	"fmt"
	"net/http"
	"regexp"
	"sync"
)

// Proxy implements [http.Handler] and implements a resolver for WissKI URIs.
//
// It should contain a map from regular expressions for URIs to target WissKI deployments.
// Keys are encoded as regular expressions, Values may contain the $n$ syntax to replace matched groups.
//
// The behaviour when multiple keys match is non-deterministic, as any value might be used as a response.
type Proxy map[string]string

// ServerHTTP implements the http.Handler interface
func (p Proxy) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	// find the uri which we should resolve
	uri := p.WissKIURI(r)
	if uri == "" {
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte("Not Found"))
		return
	}

	// determine which wisski instance
	target := p.TargetWissKI(uri)
	if target == "" {
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte("Not Found"))
		return
	}

	// append the resolver!
	final := fmt.Sprintf("%s/wisski/get?uri=%s", target, uri)
	http.Redirect(w, r, final, http.StatusFound)
}

// TargetWissKI finds the target uri that a specific URI should be redirected to.
func (p Proxy) TargetWissKI(uri string) string {
	for match, value := range p {
		r := compileCached(match)
		if r == nil {
			continue
		}
		if match := r.FindStringSubmatch(uri); match != nil {
			return r.ReplaceAllString(match[0], value)
		}
	}
	return ""
}

// WissKIURI determines the URI to redirect from.
func (Proxy) WissKIURI(r *http.Request) string {
	if r.Method != http.MethodGet {
		return ""
	}
	if r.URL.Path != "/" {
		return ""
	}
	return r.URL.Query().Get("uri")
}

const maxCacheSize = 1000

var cache = make(map[string]*regexp.Regexp, maxCacheSize)
var cacheM sync.Mutex

// compileCache compiles src into a regular expression and potentially returns a cached expression.
// It is safe to be used concurrently.
func compileCached(src string) *regexp.Regexp {
	cacheM.Lock()
	defer cacheM.Unlock()

	// cache hit
	compiled, ok := cache[src]
	if ok {
		return compiled
	}

	// compile fresh
	new, err := regexp.Compile(src)
	if err != nil {
		return nil
	}

	// if we have don't have enough space, eject a random element from the cache
	if len(cache) >= maxCacheSize {
		for key := range cache {
			delete(cache, key)
			break
		}
	}

	// store the new element in the cache, and return!
	cache[src] = new
	return new
}
