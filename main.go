package main

import (
	"dgraph/pkg/graph"
	"fmt"
	"github.com/blacktop/go-macho/types"
)

const (
	TEST_PATH = "/System/Applications/Books.app/Contents/MacOS/Books"
)

func main() {
	ggen := graph.NewImportGraphGenerator(10)

	graph := ggen.GenerateGraph(TEST_PATH, types.CPUArm64)

	fmt.Println(graph)
}
