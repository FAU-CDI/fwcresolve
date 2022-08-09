package wdresolve

import (
	_ "embed"
	"fmt"
	"html/template"
	"net/http"
	"strings"

	"golang.org/x/exp/maps"
	"golang.org/x/exp/slices"
)

// ResolveHandler implements [http.Handler] and resolves WissKI URIs to individual WissKI Resolve URIs.
type ResolveHandler struct {
	Resolver             Resolver
	TrustXForwardedProto bool
}

//go:embed index.html.tpl
var indexHTML string
var indexTemplate = template.Must(template.New("index.html").Parse(indexHTML))

// HandlerAction is an action the handler should perform
type HandlerAction int

const (
	DefaultAction HandlerAction = iota
	IndexAction                 // Show the user an index page
	ResolveAction               // resolve a URI
)

// ServerHTTP implements the http.Handler interface
func (rh ResolveHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	// find the uri which we should resolve
	action, uri := rh.Action(r)
	switch action {
	case DefaultAction:
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte("Not Found"))
	// render index html
	case IndexAction:
		w.Header().Set("Content-Type", "text/html")
		indexTemplate.Execute(w, struct {
			Prefixes [][2]string
			URL      string
		}{
			Prefixes: rh.prefixes(),
			URL:      rh.url(r),
		})
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
	default:
		panic("never reached")
	}
}

func (rh ResolveHandler) url(r *http.Request) string {
	proto := "http"
	if r.TLS != nil {
		proto = "https://"
	}
	if rh.TrustXForwardedProto {
		if p := r.Header.Get("X-Forwarded-Proto"); p != "" {
			proto = strings.ToLower(p)
		}
	}
	return fmt.Sprintf("%s://%s%s", proto, r.Host, r.URL.Path)
}

// prefixes computes the prefixes
func (rh ResolveHandler) prefixes() (pfxary [][2]string) {

	// determine the prefixes
	var prefixes map[string]string
	if presolver, ok := rh.Resolver.(PrefixResolver); ok {
		prefixes = presolver.Prefixes()
	}

	// get the prefix keys
	keys := maps.Keys(prefixes)
	slices.Sort(keys)

	// make them an array (for clean iteration)
	pfxary = make([][2]string, len(keys))
	for i, k := range keys {
		pfxary[i] = [2]string{k, prefixes[k]}
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

	// unknown
	return DefaultAction, ""
}

func (ResolveHandler) ResolverURL(target, uri string) string {
	return fmt.Sprintf("%s/wisski/get?uri=%s", target, uri)
}
