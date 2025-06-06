package expboth

import (
	"golang.org/x/exp/maps"   // want "Import statement 'golang.org/x/exp/maps' can be replaced by 'maps'"
	"golang.org/x/exp/slices" // want "Import statement 'golang.org/x/exp/slices' can be replaced by 'slices'"
)

func _(m map[string]string, a, b []string) {
	maps.Clone(m) // want `golang.org/x/exp/maps.Clone\(\) can be replaced by maps.Clone\(\)`
	slices.Equal(a, b)
}
