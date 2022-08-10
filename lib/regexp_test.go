package lib

import (
	"fmt"
	"unsafe"
)

func ExampleRegexpCache() {
	// create a cache of size 2
	cache := &RegexpCache{
		Size: 2,
	}

	// compile 'hello' regexp
	first, err := cache.Compile("hello")
	fmt.Println(first, err)

	// compile 'world' regexp
	second, err := cache.Compile("world")
	fmt.Println(second, err)

	// compile 'hello' regexp again
	third, err := cache.Compile("hello")
	fmt.Println(third, err)

	// this is re-using the same regular expression
	fmt.Println(unsafe.Pointer(first) == unsafe.Pointer(third))

	// Output: hello <nil>
	// world <nil>
	// hello <nil>
	// true
}
