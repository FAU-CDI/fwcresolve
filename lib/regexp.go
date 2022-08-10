package lib

import (
	"regexp"
	"sync"
)

// RegexpCache caches compilation of regular expressions.
// Except for modification of Size, it can be used concurrently.
type RegexpCache struct {
	// Size is the maximum size of the cache
	Size int

	m     sync.Mutex // m protects cache
	cache map[string]*regexp.Regexp
}

// Compile returns a (possibly cached) compiled version of src.
// When src is invalid it is never cached.
//
// See [regexp.Compile] for description of regular expressions.
func (rr *RegexpCache) Compile(src string) (exp *regexp.Regexp, err error) {
	rr.m.Lock()
	defer rr.m.Unlock()

	// set a default size
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

	// make sure we have enough space
	if len(rr.cache) >= rr.Size {
		rr.evict()
	}

	// store the new element in the cache, and return!
	rr.cache[src] = new
	return new, nil
}

// evictOne will evict an unspecified expression from the cache
func (rr *RegexpCache) evict() {
	// TODO: Use a random strategy for this?
	for key := range rr.cache {
		delete(rr.cache, key)
		break
	}
}

// Reset removes all cached regular expressions
func (rr *RegexpCache) Reset() {
	rr.m.Lock()
	defer rr.m.Unlock()

	rr.cache = nil
}
