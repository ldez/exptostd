package expconstraintstypes

import (
	"golang.org/x/exp/constraints" // want `Import statement 'golang.org/x/exp/constraints' may be replaced by 'cmp'`
)

type A interface {
	~string | constraints.Integer // want `golang.org/x/exp/constraints.Integer can be replaced by ~int | ~int8 | ~int16 | ~int32 | ~int64 | ~uint | ~uint8 | ~uint16 | ~uint32 | ~uint64 | ~uintptr`
}

type B interface {
	~string | constraints.Float // want `golang.org/x/exp/constraints.Float can be replaced by ~float32 | ~float64`
}

type C interface {
	~string | constraints.Signed // want `golang.org/x/exp/constraints.Float can be replaced by ~int | ~int8 | ~int16 | ~int32 | ~int64`
}

type E interface {
	~string | constraints.Unsigned // want `golang.org/x/exp/constraints.Float can be replaced by ~uint | ~uint8 | ~uint16 | ~uint32 | ~uint64 | ~uintptr`
}

type F interface {
	~string | constraints.Complex // want `golang.org/x/exp/constraints.Float can be replaced by ~complex64 | ~complex128`
}
