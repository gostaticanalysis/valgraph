package valgraph

import (
	"fmt"
	"go/token"

	"golang.org/x/tools/go/ssa"
)

type Graph struct {
	Nodes map[ssa.Value]*Node
}

func New() *Graph {
	return &Graph{Nodes: make(map[ssa.Value]*Node)}
}

func (g *Graph) CreateNode(v ssa.Value) *Node {
	n, _ := g.Nodes[v]
	if n != nil {
		return n
	}
	g.Nodes[v] = &Node{Value: v, ID: len(g.Nodes)}
	return g.Nodes[v]
}

type Node struct {
	ID    int
	Value ssa.Value
	In    []*Edge
	Out   []*Edge
}

func (n *Node) String() string {
	return fmt.Sprintf("n%d:%s", n.ID, n.Value)
}

type Edge struct {
	From *Node
	Site ssa.Instruction
	To   *Node
}

func (e Edge) String() string {
	return fmt.Sprintf("%s --> %s", e.From, e.To)
}

func (e Edge) Description() string {
	return e.Site.String()
}

func (e Edge) Pos() token.Pos {
	if e.Site == nil {
		return token.NoPos
	}
	return e.Site.Pos()
}

func AddEdge(from *Node, site ssa.Instruction, to *Node) {
	e := &Edge{From: from, Site: site, To: to}
	from.Out = append(from.Out, e)
	to.In = append(to.In, e)
}
