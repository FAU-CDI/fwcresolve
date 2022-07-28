package resolvers

import (
	"regexp"
	"sync"
)

// Regexp implements [fwcresolver.Resolver] on the basis of a regexp mapping.
//
// It should contain a map from regular expressions for URIs to target WissKI deployments.
// Keys are encoded as regular expressions, Values may contain the $n$ syntax to replace matched groups.
//
// The behaviour when multiple keys match is non-deterministic, as any value might be used as a response.
type Regexp struct {
	Data map[string]string

	cache RegexpCache
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

// RegexpCache caches compilation of regular expressions.
// Except for modification of Size, it can be used concurrently.
type RegexpCache struct {
	Size int

	m     sync.Mutex                // m protects cache
	cache map[string]*regexp.Regexp // cache
}

// Compile returns a (possibly cached) compiled version of src
// When
func (rr *RegexpCache) Compile(src string) (exp *regexp.Regexp, err error) {
	rr.m.Lock()
	defer rr.m.Unlock()

	if rr.Size == 0 {
		rr.Size = 1000
	}

	// create cache object if it doesn't exist
	if rr.cache == nil {
		rr.cache = make(map[string]*regexp.Regexp, rr.Size)
	}

	// cache hit
	compiled, ok := rr.cache[src]
	if ok {
		return compiled, nil
	}

	// compile fresh
	new, err := regexp.Compile(src)
	if err != nil {
		return nil, err
	}

	// if we have don't have enough space, eject a random element from the cache
	if len(rr.cache) >= rr.Size {
		for key := range rr.cache {
			delete(rr.cache, key)
			break
		}
	}

	// store the new element in the cache, and return!
	rr.cache[src] = new
	return new, nil
}
