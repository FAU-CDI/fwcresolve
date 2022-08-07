package wdresolve

import (
	_ "embed"
	"encoding/json"
	"fmt"
	"net/http"
)

// ResolveHandler implements [http.Handler] and resolves WissKI URIs to individual WissKI Resolve URIs.
type ResolveHandler struct {
	Resolver Resolver
}

//go:embed index.html
var indexHTML []byte

// HandlerAction is an action the handler should perform
type HandlerAction int

const (
	DefaultAction  HandlerAction = iota
	IndexAction                  // Show the user an index page
	ResolveAction                // resolve a URI
	PrefixesAction               // list known prefixes
)

// ServerHTTP implements the http.Handler interface
func (rh ResolveHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	// find the uri which we should resolve
	action, uri := rh.Action(r)
	switch action {
	case DefaultAction:
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte("Not Found"))
	case IndexAction:
		w.Header().Set("Content-Type", "text/html")
		w.Write(indexHTML)
	case ResolveAction:
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
	case PrefixesAction:
		// return prefixes as json
		w.WriteHeader(http.StatusOK)
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(rh.prefixes())
	default:
		panic("never reached")
	}
}

// prefixes computes the prefixes
func (rh ResolveHandler) prefixes() (prefixes map[string]string) {
	if presolver, ok := rh.Resolver.(PrefixResolver); ok {
		prefixes = presolver.Prefixes()
	}
	if prefixes == nil {
		prefixes = make(map[string]string, 0)
	}
	return
}

// Action extracts the action and uri to perform the action on
// When no uri exists, returns the empty string.
func (ResolveHandler) Action(r *http.Request) (action HandlerAction, uri string) {
	if r.Method != http.MethodGet {
		return DefaultAction, ""
	}

	// no parameters => return index
	if r.URL.RawQuery == "" {
		return IndexAction, ""
	}

	query := r.URL.Query()

	if uri = query.Get("uri"); uri != "" {
		return ResolveAction, uri
	}

	if query.Has("prefixes") {
		return PrefixesAction, ""
	}

	// unknown
	return DefaultAction, ""
}

func (ResolveHandler) ResolverURL(target, uri string) string {
	return fmt.Sprintf("%s/wisski/get?uri=%s", target, uri)
}
