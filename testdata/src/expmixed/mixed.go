package expmixed

import (
	"expmixed/maps"

	"golang.org/x/exp/slices" // want "Import statement 'golang.org/x/exp/slices' can be replaced by 'slices'"
)

func _(m map[string]string, a, b []string) {
	maps.Clone(m)
	slices.Equal(a, b)
}
