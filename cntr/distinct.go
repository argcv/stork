package cntr

import (
	"strings"
)

// DistinctStrings will return distinct strings, case sensitive
func DistinctStrings(elems ...string) []string {
	m := map[string]bool{}
	var retElems []string
	for _, e := range elems {
		if _, ok := m[e]; !ok {
			m[e] = true
			retElems = append(retElems, e)
		}
	}
	return retElems
}

// DistinctStringsCaseInsensitive will return distinct strings, case insensitive
func DistinctStringsCaseInsensitive(elems ...string) []string {
	m := map[string]bool{}
	var retElems []string
	for _, e := range elems {
		le := strings.ToLower(e)
		if _, ok := m[le]; !ok {
			m[le] = true
			retElems = append(retElems, e)
		}
	}
	return retElems
}
