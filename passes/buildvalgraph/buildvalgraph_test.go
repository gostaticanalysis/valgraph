package buildvalgraph_test

import (
	"fmt"
	"testing"

	"github.com/gostaticanalysis/testutil"
	"github.com/gostaticanalysis/valgraph"
	"github.com/gostaticanalysis/valgraph/passes/buildvalgraph"
	"golang.org/x/tools/go/analysis"
	"golang.org/x/tools/go/analysis/analysistest"
	"golang.org/x/tools/go/analysis/passes/buildssa"
	"golang.org/x/tools/go/ssa"
)

// TestAnalyzer is a test for Analyzer.
func TestAnalyzer(t *testing.T) {
	a := &analysis.Analyzer{
		Name: "test",
		Doc:  "test",
		Requires: []*analysis.Analyzer{
			buildssa.Analyzer,
			buildvalgraph.Analyzer,
		},
		Run: func(pass *analysis.Pass) (interface{}, error) {
			s := pass.ResultOf[buildssa.Analyzer].(*buildssa.SSA)
			g := pass.ResultOf[buildvalgraph.Analyzer].(*valgraph.Graph)
			for _, m := range s.Pkg.Members {
				v, _ := m.(ssa.Value)
				n := g.Nodes[v]
				if n != nil {
					fmt.Println(n)
				}
			}
			return nil, nil
		},
	}
	testdata := testutil.WithModules(t, analysistest.TestData(), nil)
	analysistest.Run(t, testdata, a, "a")
}
