// Package exptostd It is an analyzer that detects functions from golang.org/x/exp/ that can be replaced by std functions.
package exptostd

import (
	"go/ast"
	"go/build"
	"go/types"
	"os"
	"slices"
	"strconv"
	"strings"

	"golang.org/x/tools/go/analysis"
	"golang.org/x/tools/go/analysis/passes/inspect"
	"golang.org/x/tools/go/ast/inspector"
)

const (
	go123   = 123
	go121   = 121
	goDevel = 666
)

type stdReplacement struct {
	MinGo int
	Text  string
}

type analyzer struct {
	mapsPkgReplacements   map[string]stdReplacement
	slicesPkgReplacements map[string]stdReplacement

	skipGoVersionDetection bool
	goVersion              int
}

// NewAnalyzer create a new Analyzer.
func NewAnalyzer() *analysis.Analyzer {
	_, skip := os.LookupEnv("EXPTOSTD_SKIP_GO_VERSION_CHECK")

	l := &analyzer{
		skipGoVersionDetection: skip,
		mapsPkgReplacements: map[string]stdReplacement{
			"Keys":       {MinGo: go123, Text: "slices.Collect(maps.Keys())"},
			"Values":     {MinGo: go123, Text: "slices.Collect(maps.Keys())"},
			"Equal":      {MinGo: go121, Text: "maps.Equal()"},
			"EqualFunc":  {MinGo: go121, Text: "maps.EqualFunc()"},
			"Clone":      {MinGo: go121, Text: "maps.Clone()"},
			"Copy":       {MinGo: go121, Text: "maps.Copy()"},
			"DeleteFunc": {MinGo: go121, Text: "maps.DeleteFunc()"},
			"Clear":      {MinGo: go121, Text: "clear()"},
		},
		slicesPkgReplacements: map[string]stdReplacement{
			"Equal":        {MinGo: go121, Text: "slices.Equal()"},
			"EqualFunc":    {MinGo: go121, Text: "slices.EqualFunc()"},
			"Compare":      {MinGo: go121, Text: "slices.Compare()"},
			"CompareFunc":  {MinGo: go121, Text: "slices.CompareFunc()"},
			"Index":        {MinGo: go121, Text: "slices.Index()"},
			"IndexFunc":    {MinGo: go121, Text: "slices.IndexFunc()"},
			"Contains":     {MinGo: go121, Text: "slices.Contains()"},
			"ContainsFunc": {MinGo: go121, Text: "slices.ContainsFunc()"},
			"Insert":       {MinGo: go121, Text: "slices.Insert()"},
			"Delete":       {MinGo: go121, Text: "slices.Delete()"},
			"DeleteFunc":   {MinGo: go121, Text: "slices.DeleteFunc()"},
			"Replace":      {MinGo: go121, Text: "slices.Replace()"},
			"Clone":        {MinGo: go121, Text: "slices.Clone()"},
			"Compact":      {MinGo: go121, Text: "slices.Compact()"},
			"CompactFunc":  {MinGo: go121, Text: "slices.CompactFunc()"},
			"Grow":         {MinGo: go121, Text: "slices.Grow()"},
			"Clip":         {MinGo: go121, Text: "slices.Clip()"},
			"Reverse":      {MinGo: go121, Text: "slices.Reverse()"},

			"Sort":             {MinGo: go121, Text: "slices.Sort()"},
			"SortFunc":         {MinGo: go121, Text: "slices.SortFunc()"},
			"SortStableFunc":   {MinGo: go121, Text: "slices.SortStableFunc()"},
			"IsSorted":         {MinGo: go121, Text: "slices.IsSorted()"},
			"IsSortedFunc":     {MinGo: go121, Text: "slices.IsSortedFunc()"},
			"Min":              {MinGo: go121, Text: "slices.Min()"},
			"MinFunc":          {MinGo: go121, Text: "slices.MinFunc()"},
			"Max":              {MinGo: go121, Text: "slices.Max()"},
			"MaxFunc":          {MinGo: go121, Text: "slices.MaxFunc()"},
			"BinarySearch":     {MinGo: go121, Text: "slices.BinarySearch()"},
			"BinarySearchFunc": {MinGo: go121, Text: "slices.BinarySearchFunc()"},
		},
	}

	return &analysis.Analyzer{
		Name:     "exptostd",
		Doc:      "Detects functions from golang.org/x/exp/ that can be replaced by std functions.",
		Run:      l.run,
		Requires: []*analysis.Analyzer{inspect.Analyzer},
	}
}

func (a *analyzer) run(pass *analysis.Pass) (interface{}, error) {
	insp, ok := pass.ResultOf[inspect.Analyzer].(*inspector.Inspector)
	if !ok {
		return nil, nil
	}

	a.goVersion = getGoVersion(pass)

	nodeFilter := []ast.Node{
		(*ast.CallExpr)(nil),
		(*ast.ImportSpec)(nil),
	}

	imports := map[string]*ast.ImportSpec{}
	importsReplaces := strings.NewReplacer("\"", "", "`", "")
	shouldKeepExpMaps := false
	shouldKeepExpSlices := false

	insp.Preorder(nodeFilter, func(n ast.Node) {
		if importSpec, ok := n.(*ast.ImportSpec); ok {
			cleanedPath := importsReplaces.Replace(importSpec.Path.Value)
			imports[cleanedPath] = importSpec

			return
		}

		callExpr, ok := n.(*ast.CallExpr)
		if !ok {
			return
		}

		selExpr, ok := callExpr.Fun.(*ast.SelectorExpr)
		if !ok {
			return
		}

		ident, ok := selExpr.X.(*ast.Ident)
		if !ok {
			return
		}

		switch ident.Name {
		case "maps":
			shouldKeepExpMaps = shouldKeepExpMaps ||
				!a.detectPackageUsage(pass, a.mapsPkgReplacements, selExpr, ident, callExpr, "golang.org/x/exp/maps")

		case "slices":
			shouldKeepExpSlices = shouldKeepExpSlices ||
				!a.detectPackageUsage(pass, a.slicesPkgReplacements, selExpr, ident, callExpr, "golang.org/x/exp/slices")
		}
	})

	maybeReplaceImport(pass, imports, shouldKeepExpMaps, shouldKeepExpSlices)

	return nil, nil
}

func maybeReplaceImport(
	pass *analysis.Pass,
	imports map[string]*ast.ImportSpec,
	shouldKeepExpMaps bool,
	shouldKeepExpSlices bool,
) {
	diagnostic := func(importSpec *ast.ImportSpec, replaceText string) analysis.Diagnostic {
		return analysis.Diagnostic{
			Pos:     importSpec.Pos(),
			End:     importSpec.End(),
			Message: "Import statement can drop `golang.org/x/exp` prefix",
			SuggestedFixes: []analysis.SuggestedFix{
				{
					TextEdits: []analysis.TextEdit{
						{
							Pos:     importSpec.Path.Pos(),
							End:     importSpec.Path.End(),
							NewText: []byte(replaceText),
						},
					},
				},
			},
		}
	}

	if imp, ok := imports["golang.org/x/exp/maps"]; ok && !shouldKeepExpMaps {
		if imp.Name == nil || imp.Name.Name == "" {
			pass.Report(diagnostic(imp, `"maps"`))
		}
	}

	if imp, ok := imports["golang.org/x/exp/slices"]; ok && !shouldKeepExpSlices {
		if imp.Name == nil || imp.Name.Name == "" {
			pass.Report(diagnostic(imp, `"slices"`))
		}
	}
}

func (a *analyzer) detectPackageUsage(pass *analysis.Pass,
	replacements map[string]stdReplacement,
	selExpr *ast.SelectorExpr, ident *ast.Ident, callExpr *ast.CallExpr,
	importPath string,
) bool {
	rp, ok := replacements[selExpr.Sel.Name]
	if !ok {
		return false
	}

	if !a.skipGoVersionDetection && rp.MinGo > a.goVersion {
		return false
	}

	obj := pass.TypesInfo.Uses[ident]
	if obj == nil {
		return false
	}

	pkg, ok := obj.(*types.PkgName)
	if !ok {
		return false
	}

	if pkg.Imported().Path() == importPath {
		pass.Reportf(callExpr.Pos(), "%s.%s() can be replaced by %s", importPath, selExpr.Sel.Name, rp.Text)
		return true
	}

	return false
}

func getGoVersion(pass *analysis.Pass) int {
	// Prior to go1.22, versions.FileVersion returns only the toolchain version,
	// which is of no use to us,
	// so disable this analyzer on earlier versions.
	if !slices.Contains(build.Default.ReleaseTags, "go1.22") {
		return 0 // false
	}

	pkgVersion := pass.Pkg.GoVersion()
	if pkgVersion == "" {
		// Empty means Go devel.
		return goDevel // true
	}

	vParts := strings.Split(strings.TrimPrefix(pkgVersion, "go"), ".")

	v, err := strconv.Atoi(strings.Join(vParts[:2], ""))
	if err != nil {
		v = 116
	}

	return v
}
