package builder

import (
	"github.com/gostaticanalysis/analysisutil"
	"github.com/gostaticanalysis/valgraph"
	"golang.org/x/tools/go/callgraph"
	"golang.org/x/tools/go/ssa"
)

func ByCallGraph(prog *ssa.Program, cg *callgraph.Graph) *valgraph.Graph {
	b := &byCallGraph{
		g:    valgraph.New(),
		cg:   cg,
		prog: prog,
	}
	b.traverse()
	return b.g
}

type byCallGraph struct {
	g    *valgraph.Graph
	cg   *callgraph.Graph
	prog *ssa.Program
}

func (b *byCallGraph) traverse() {
	for _, pkg := range b.prog.AllPackages() {
		b.traversePkg(pkg)
	}
}

func (b *byCallGraph) traversePkg(pkg *ssa.Package) {
	for _, m := range pkg.Members {
		b.traverseMember(m)
	}
}

func (b *byCallGraph) traverseMember(m ssa.Member) {
	switch m := m.(type) {
	case *ssa.Function:
		b.traverseFunc(nil, nil, m)
	case ssa.Value:
		b.traverseValue(nil, nil, m)
	}
}

func (b *byCallGraph) traverseFunc(from *valgraph.Node, site ssa.Instruction, f *ssa.Function) {
	if f == nil || b.g.Nodes[f] != nil {
		return
	}

	to := b.g.CreateNode(f)
	if from != nil {
		valgraph.AddEdge(from, site, to)
	}

	rets := analysisutil.Returns(f)
	for _, ret := range rets {
		for _, v := range ret.Results {
			b.traverseValue(to, ret, v)
		}
	}
}

func (b *byCallGraph) traverseValue(from *valgraph.Node, site ssa.Instruction, v ssa.Value) {
	if v == nil || b.g.Nodes[v] != nil {
		return
	}

	to := b.g.CreateNode(v)
	if from != nil {
		valgraph.AddEdge(from, site, to)
	}

	switch v := v.(type) {
	case *ssa.Call:
		fun := b.callee(v)
		b.traverseFunc(to, v, fun)
		for _, o := range v.Operands(nil) {
			if o != nil {
				b.traverseValue(to, v, *o)
			}
		}
	}
}

func (b *byCallGraph) callee(call ssa.CallInstruction) *ssa.Function {
	node := b.cg.Nodes[call.Parent()]
	if node == nil {
		return nil
	}

	for _, edge := range node.Out {
		if edge.Site == call {
			return edge.Callee.Func
		}
	}

	return nil
}
