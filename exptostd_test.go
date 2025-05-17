package exptostd_test

import (
	"os"
	"os/exec"
	"path/filepath"
	"testing"

	"github.com/ldez/exptostd"
	"golang.org/x/tools/go/analysis"
	"golang.org/x/tools/go/analysis/analysistest"
)

func TestAnalyzer(t *testing.T) {
	testCases := []struct {
		dir string
	}{
		{dir: "expmaps"},
		{dir: "expmaps122"},
		{dir: "expslices"},
		{dir: "expboth"},
		{dir: "expnone"},
		{dir: "expmixed"},
		{dir: "expalias"},
		{dir: "expconstraints"},
		{dir: "expconstraintskeep"},
	}

	for _, test := range testCases {
		t.Run(test.dir, func(t *testing.T) {
			runWithSuggestedFixes(t, exptostd.NewAnalyzer(), test.dir)
		})
	}
}

func runWithSuggestedFixes(t *testing.T, a *analysis.Analyzer, dir string, patterns ...string) []*analysistest.Result {
	t.Helper()

	wd, err := os.Getwd()
	if err != nil {
		t.Fatal(err)
	}

	defer func() { _ = os.Chdir(wd) }()

	tempDir := t.TempDir()

	// Needs to be run outside testdata.
	err = os.CopyFS(tempDir, os.DirFS(filepath.Join(analysistest.TestData(), "src")))
	if err != nil {
		t.Fatal(err)
	}

	// NOTE: analysistest does not yet support modules;
	// see https://github.com/golang/go/issues/37054 for details.

	srcPath := filepath.Join(tempDir, filepath.FromSlash(dir))

	err = os.Chdir(srcPath)
	if err != nil {
		t.Fatal(err)
	}

	output, err := exec.Command("go", "mod", "vendor").CombinedOutput()
	if err != nil {
		t.Log(string(output))
		t.Fatal(err)
	}

	return analysistest.RunWithSuggestedFixes(t, srcPath, a, patterns...)
}
