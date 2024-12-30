package expalias

import (
	"golang.org/x/exp/maps" // want "Import statement 'golang.org/x/exp/maps' can be replaced by 'maps'"
	aliasMaps "golang.org/x/exp/maps"
	"golang.org/x/exp/slices" // want "Import statement 'golang.org/x/exp/slices' can be replaced by 'slices'"
	aliasSlices "golang.org/x/exp/slices"
)

func _(m map[string]string, a, b []string) {
	aliasMaps.Clone(m)
	maps.Clone(m) // want `golang.org/x/exp/maps\.Clone\(\) can be replaced by maps\.Clone\(\)`
	aliasSlices.Equal(a, b)
	slices.Equal(a, b)
}
