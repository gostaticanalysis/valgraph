package buildvalgraph

import (
	"reflect"

	"github.com/gostaticanalysis/valgraph"
	"github.com/gostaticanalysis/valgraph/builder"
	"golang.org/x/tools/go/analysis"
	"golang.org/x/tools/go/analysis/passes/buildssa"
	"golang.org/x/tools/go/callgraph/cha"
	"golang.org/x/tools/go/callgraph/vta"
	"golang.org/x/tools/go/ssa/ssautil"
)

const doc = "buildvalgraph is ..."

var Analyzer = &analysis.Analyzer{
	Name: "buildvalgraph",
	Doc:  doc,
	Run:  run,
	Requires: []*analysis.Analyzer{
		buildssa.Analyzer,
	},
	ResultType: reflect.TypeOf((*valgraph.Graph)(nil)),
}

func run(pass *analysis.Pass) (interface{}, error) {
	s := pass.ResultOf[buildssa.Analyzer].(*buildssa.SSA)
	cg := vta.CallGraph(ssautil.AllFunctions(s.Pkg.Prog), cha.CallGraph(s.Pkg.Prog))
	return builder.ByCallGraph(s.Pkg.Prog, cg), nil
}
