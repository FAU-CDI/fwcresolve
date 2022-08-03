package wdresolve

import "strings"

// Resolver resolves URIs to WissKI base URIs
type Resolver interface {
	Target(uri string) string
}

// PrefixResolver is a resolver that resolves WissKI base URIs using a longest prefix match.
// Resolver should call PrefixTarget()
type PrefixResolver interface {
	Resolver
	Prefixes() map[string]string
}

// PrefixTarget implemen s Target for PrefixResolvers
func PrefixTarget(resolver PrefixResolver, uri string) (url string) {
	var match string
	for key, value := range resolver.Prefixes() {
		// check if we have a longer prefix on the string
		if strings.HasPrefix(uri, key) && len(key) > len(match) {
			match = key
			url = value
		}
	}
	return
}
