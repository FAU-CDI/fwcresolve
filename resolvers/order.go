package resolvers

import (
	"github.com/tkw1536/wdresolve"
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
