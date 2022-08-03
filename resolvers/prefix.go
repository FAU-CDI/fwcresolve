package resolvers

import (
	"bufio"
	"errors"
	"io"
	"strings"

	"github.com/tkw1536/wdresolve"
)

// Prefix implements [wdresolve.Resolver] on the basis of a longest prefix match
type Prefix struct {
	Data map[string]string
}

func (prefix Prefix) Prefixes() map[string]string {
	return prefix.Data
}

func (prefix Prefix) Target(uri string) (url string) {
	return wdresolve.PrefixTarget(prefix, uri)
}

var ErrNoDestination = errors.New("encountered prefix without a destination")

// ReadPrefixes reads a set of prefixes from the provided reader
//
// The reader should return data that contains one prefix per line.
// Leading and trailing spaces are ignored.
// A destination is indicated by a line ending in ":".
// Blank lines and those starting with "#" are ignored.
// For example:
//
//    # I am ignored
//
//    https://wisski.example.com:
//      http://wisski.example.com/
//      https://wisski.example.com/
//    https://wisski.other.com:
//      http://wisski.other.com/
//      https://wisski.other.com/
//
func ReadPrefixes(reader io.Reader) (p Prefix, err error) {
	// read the file line by line
	scanner := bufio.NewScanner(reader)
	scanner.Split(bufio.ScanLines)

	var dest string
	data := make(map[string]string)

	for scanner.Scan() {
		s := strings.TrimSpace(scanner.Text())

		// found a new destination
		if strings.HasSuffix(s, ":") {
			dest = s[:len(s)-1]
			continue
		}

		// ignore empty and comment lines
		if len(s) == 0 || s[0] == '#' {
			continue
		}

		if len(dest) == 0 {
			return p, err
		}

		// store the destination
		data[s] = dest
	}

	// return with error
	if err = scanner.Err(); err != nil {
		return
	}

	// return with proper data
	p.Data = data
	return
}
