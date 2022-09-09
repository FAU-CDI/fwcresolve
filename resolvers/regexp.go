package resolvers

import "github.com/FAU-CDI/wdresolve/lib"

// Regexp implements [wdresolver.Resolver] on the basis of a regexp mapping.
//
// It should contain a map from regular expressions for URIs to target WissKI deployments.
// Keys are encoded as regular expressions, Values may contain the $n$ syntax to replace matched groups.
//
// The behaviour when multiple keys match is non-deterministic, as any value might be used as a response.
type Regexp struct {
	Data map[string]string

	cache lib.RegexpCache
}

func (rr *Regexp) Target(uri string) string {
	for match, value := range rr.Data {
		r, err := rr.cache.Compile(match)
		if r == nil || err != nil {
			continue
		}
		if match := r.FindStringSubmatch(uri); match != nil {
			return r.ReplaceAllString(match[0], value)
		}
	}
	return ""
}
