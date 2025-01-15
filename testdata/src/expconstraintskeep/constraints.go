package expconstraintskeep

import (
	"golang.org/x/exp/constraints"
)

type _ interface {
	constraints.Ordered | constraints.Float // want "golang.org/x/exp/constraints.Ordered can be replaced by cmp.Ordered"
}
