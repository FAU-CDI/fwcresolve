package resolvers

import (
	"github.com/FAU-CDI/wdresolve"
)

// InOrder implements [wdresolve.Resolver] and queries several child resolvers in order
type InOrder []wdresolve.Resolver

func (io InOrder) Target(uri string) string {
	for _, r := range io {
		if target := r.Target(uri); target != "" {
			return target
		}
	}
	return ""
}

func (io InOrder) Prefixes() map[string]string {
	result := make(map[string]string)
	for _, r := range io {
		pr, isPr := r.(wdresolve.PrefixResolver)
		if !isPr {
			continue
		}
		for key, value := range pr.Prefixes() {
			result[key] = value
		}
	}
	return result
}
