package graph

type GraphNode struct {
	Path         string
	Dependencies []*GraphNode
}
