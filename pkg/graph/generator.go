package graph

import (
	"bytes"
	"github.com/blacktop/go-macho/types"
	"github.com/goccy/go-graphviz"
	"log"
	"strings"
)

type DependencyGraphGenerator struct {
	// MaximumDepth represents the maximum allowed parse depth.
	MaximumDepth uint

	gv       *graphviz.Graphviz
	rootNode *GraphNode
	nodes    map[string]*GraphNode
}

func NewImportGraphGenerator(maxDepth uint) *DependencyGraphGenerator {
	return &DependencyGraphGenerator{
		MaximumDepth: maxDepth,
		gv:           graphviz.New(),
	}
}

func (g *DependencyGraphGenerator) GenerateGraph(rootNode string, cpu types.CPU) string {
	graph, err := g.gv.Graph()
	if err != nil {
		panic(err)
	}

	graph.CreateNode(rootNode)

	stack := []string{rootNode}
	visited := map[string]bool{}

	for len(stack) > 0 {
		// pop element from the stack...
		current := stack[len(stack)-1]

		//log.Println(current, len(stack))

		currentNode, err := graph.CreateNode(current)
		if err != nil {
			println(err)
		}

		stack = stack[:len(stack)-1]

		dependencies, err := ParseDependencies(current, cpu)
		if err != nil {
			//log.Println(err)
			continue
		}

		for _, dep := range dependencies {

			// if we have already visited the node, skip it.
			if _, ok := visited[dep]; ok {
				continue
			}

			if strings.HasPrefix(dep, "@executable_path") || strings.HasPrefix(dep, "@loader_path") {
				visited[dep] = true
				continue
			}

			depResolved, resolveErr := resolvePath(dep, current)
			if resolveErr == nil && strings.HasSuffix(current, depResolved) {
				continue
			}

			// add the dependency node to the graph
			depNode, _ := graph.CreateNode(dep)

			// create an edge between the current node and the dependency node
			graph.CreateEdge("test", currentNode, depNode)

			// push the dependency onto the stack for further processing
			if resolveErr != nil {
				stack = append(stack, dep)
			} else {
				stack = append(stack, depResolved)
			}
		}
		// mark the dependency as visited
		visited[current] = true
	}

	var buf bytes.Buffer
	if err := g.gv.Render(graph, "dot", &buf); err != nil {
		log.Fatal(err)
	}

	return buf.String()
}
