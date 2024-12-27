package expboth

import (
	"golang.org/x/exp/maps"   // want "Import statement can drop `golang.org/x/exp` prefix"
	"golang.org/x/exp/slices" // want "Import statement can drop `golang.org/x/exp` prefix"
)

func _(m map[string]string, a, b []string) {
	maps.Clone(m)      // want `golang.org/x/exp/maps.Clone\(\) can be replaced by maps.Clone\(\)`
	slices.Equal(a, b) // want `golang.org/x/exp/slices\.Equal\(\) can be replaced by slices\.Equal\(\)`
}
