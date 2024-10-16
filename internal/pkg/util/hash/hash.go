// Package hash provides utils to calculate hashes.
package hash

import (
	"fmt"
	"hash/fnv"

	"github.com/davecgh/go-spew/spew"
)

// Compute computes the hash of an object using the spew library.
// Note: spew follows pointers and prints actual values of the nested objects
// ensuring the hash does not change when a pointer changes.
func Compute(obj any) (uint32, error) {
	hasher := fnv.New32a()

	printer := spew.ConfigState{
		Indent:         " ",
		SortKeys:       true,
		DisableMethods: true,
		SpewKeys:       true,
	}

	if _, err := printer.Fprintf(hasher, "%#v", obj); err != nil {
		return 0, fmt.Errorf("failed to calculate hash")
	}

	return hasher.Sum32(), nil
}
