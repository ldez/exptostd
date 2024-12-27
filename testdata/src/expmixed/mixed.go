package expmixed

import (
	"expmixed/maps"

	"golang.org/x/exp/slices" // want "Import statement can drop `golang.org/x/exp` prefix"
)

func _(m map[string]string, a, b []string) {
	maps.Clone(m)
	slices.Equal(a, b) // want `golang.org/x/exp/slices\.Equal\(\) can be replaced by slices\.Equal\(\)`
}
