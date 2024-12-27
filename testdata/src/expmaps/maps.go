package expmaps

import (
	`golang.org/x/exp/maps`
)

func _(m, a map[string]string) {
	maps.Clone(m) // want `golang.org/x/exp/maps.Clone\(\) can be replaced by maps.Clone\(\)`

	maps.Keys(m) // want `golang.org/x/exp/maps.Keys\(\) can be replaced by slices.Collect\(maps.Keys\(\)\)`

	maps.Values(m) // want `golang.org/x/exp/maps.Values\(\) can be replaced by slices.Collect\(maps.Keys\(\)\)`

	maps.Equal(m, a) // want `golang.org/x/exp/maps.Equal\(\) can be replaced by maps.Equal\(\)`

	maps.EqualFunc(m, a, func(i, j string) bool { // want `golang.org/x/exp/maps.EqualFunc\(\) can be replaced by maps.EqualFunc\(\)`
		return true
	})

	maps.Copy(m, a) // want `golang.org/x/exp/maps.Copy\(\) can be replaced by maps.Copy\(\)`

	maps.DeleteFunc(m, func(_, _ string) bool { // want `golang.org/x/exp/maps.DeleteFunc\(\) can be replaced by maps.DeleteFunc\(\)`
		return true
	})

	maps.Clear(m) // want `golang.org/x/exp/maps.Clear\(\) can be replaced by clear\(\)`
}
