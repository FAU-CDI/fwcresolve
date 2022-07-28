package resolvers

import "strings"

// Prefix implements [fwcresolver.Resolver] on the basis of a longest prefix match
type Prefix struct {
	Data map[string]string
}

func (prefix Prefix) Target(uri string) (url string) {
	var match string
	for key, value := range prefix.Data {
		// check if we have a longer prefix on the string
		if strings.HasPrefix(uri, key) && len(key) > len(match) {
			match = key
			url = value
		}
	}
	return
}
