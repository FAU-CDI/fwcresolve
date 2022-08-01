package wdresolve

import (
	"fmt"
	"net/http"
)

// ResolveHandler implements [http.Handler] and resolves WissKI URIs to individual WissKI Resolve URIs.
type ResolveHandler struct {
	Resolver Resolver
}

// Resolver resolves URIs to WissKI base URIs
type Resolver interface {
	Target(uri string) string
}

// ServerHTTP implements the http.Handler interface
func (rh ResolveHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	// find the uri which we should resolve
	uri := rh.URI(r)
	if uri == "" {
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte("Not Found"))
		return
	}

	// determine which wisski instance
	target := rh.Resolver.Target(uri)
	if target == "" {
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte("Not Found"))
		return
	}

	// append the resolver!
	final := rh.ResolverURL(target, uri)
	http.Redirect(w, r, final, http.StatusFound)
}

// URI extracts the uri to resolve from a request.
// When no uri exists, returns the empty string.
func (ResolveHandler) URI(r *http.Request) string {
	if r.Method != http.MethodGet {
		return ""
	}
	return r.URL.Query().Get("uri")
}

func (ResolveHandler) ResolverURL(target, uri string) string {
	return fmt.Sprintf("%s/wisski/get?uri=%s", target, uri)
}
