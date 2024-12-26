// Package main contains the basic runnable version of the linter.
package main

import (
	"github.com/ldez/exptostd"
	"golang.org/x/tools/go/analysis/singlechecker"
)

func main() {
	singlechecker.Main(exptostd.NewAnalyzer())
}
