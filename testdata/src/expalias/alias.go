package expalias

import (
	"golang.org/x/exp/maps"
	aliased "golang.org/x/exp/maps"
	"golang.org/x/exp/slices" // want "Import statement 'golang.org/x/exp/slices' can be replaced by 'slices'"
)

func _(m map[string]string, a, b []string) {
	aliased.Clone(m)
	maps.Clone(m)      // want `golang.org/x/exp/maps\.Clone\(\) can be replaced by maps\.Clone\(\)`
	slices.Equal(a, b) // want `golang.org/x/exp/slices\.Equal\(\) can be replaced by slices\.Equal\(\)`
}
