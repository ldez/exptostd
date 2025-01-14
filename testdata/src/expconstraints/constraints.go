package expconstraints

import (
	"time"

	"golang.org/x/exp/constraints" // want "Import statement 'golang.org/x/exp/constraints' can be replaced by 'cmp'"
)

type _[T constraints.Ordered] []T // want `golang.org/x/exp/constraints\.Ordered can be replaced by cmp\.Ordered`

type _ interface {
	constraints.Ordered | time.Time // want `golang.org/x/exp/constraints\.Ordered can be replaced by cmp\.Ordered`
}

type _ interface {
	constraints.Ordered // want `golang.org/x/exp/constraints\.Ordered can be replaced by cmp\.Ordered`
	comparable
}

type _[K constraints.Ordered, V any] struct { // want `golang.org/x/exp/constraints\.Ordered can be replaced by cmp\.Ordered`
}

func _[T constraints.Ordered](_ []T) bool { // want `golang.org/x/exp/constraints\.Ordered can be replaced by cmp\.Ordered`
	return true
}

func _[S []E, E constraints.Ordered](s S) E { // want `golang.org/x/exp/constraints\.Ordered can be replaced by cmp\.Ordered`
	return s[len(s)-1]
}
