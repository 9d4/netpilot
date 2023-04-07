// Package prefix contains some constants such as redis key
// that can be used across the program.
package prefix

import "fmt"

const (
	BoardPrefix        keyPrefix = "board:"
	BoardChannelPrefix keyPrefix = "ch:board:"
)

type keyPrefix string

// Get concat prefix with uuid and attributes
// e.g. prefix:uuid:attr1:attr2
// Every attribute will be concatenated with colon ':',
// except for the last one
func (k keyPrefix) Get(uuid string, attributes ...string) string {
	var attrs string
	for _, attr := range attributes {
		attrs += ":" + attr
	}

	return fmt.Sprintf("%s%s%s", k, uuid, attrs)
}

func (k keyPrefix) Status(uuid string) string {
	return k.Get(uuid, "status")
}
