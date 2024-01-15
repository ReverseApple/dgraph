package main

import (
	"fmt"
	"github.com/ReverseApple/dgraph/pkg/graph"
	"github.com/blacktop/go-macho/types"
	"io"
	"os"
)

const (
	TEST_PATH = "/System/Applications/Books.app/Contents/MacOS/Books"
)

func main() {
	ggen := graph.NewImportGraphGenerator(10)

	graph, err := ggen.GenerateGraph(TEST_PATH, types.CPUArm64)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error generating graph: %v\n", err)
		os.Exit(1)
	}

	io.Copy(os.Stdout, graph)

	//fmt.Println(graph)
}
