package graph

import "github.com/goccy/go-graphviz"

type DependencyGraphGenerator struct {
	// MaximumDepth represents the maximum allowed parse depth.
	MaximumDepth uint

	gv       *graphviz.Graphviz
	rootNode *GraphNode
	nodes    map[string]*GraphNode
}

func (g *DependencyGraphGenerator) getOrCreateNode(id string) *GraphNode {
	if node, ok := g.nodes[id]; ok {
		return node
	}

	node := &GraphNode{
		Path:         id,
		Dependencies: []*GraphNode{},
	}
	g.nodes[id] = node
	return node
}

func NewImportGraphGenerator(maxDepth uint) *DependencyGraphGenerator {
	return &DependencyGraphGenerator{
		MaximumDepth: maxDepth,
		gv:           graphviz.New(),
	}
}

func (g *DependencyGraphGenerator) GenerateGraph(rootNode string) {
	graph, err := g.gv.Graph()
	if err != nil {
		panic(err)
	}

	g.rootNode = &GraphNode{Path: rootNode}

	stack := []*GraphNode{}

}
