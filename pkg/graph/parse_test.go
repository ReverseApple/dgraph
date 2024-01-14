package graph

import (
	"github.com/blacktop/go-macho/types"
	"testing"
)

func TestParseDependencies(t *testing.T) {
	type testCase struct {
		name        string
		path        string
		cpu         types.CPU
		expectError bool
		expectNone  bool
	}

	testCases := []testCase{
		{
			"Valid fat binary with correct CPU",
			"/System/Applications/Books.app/Contents/MacOS/Books",
			types.CPUArm64,
			false,
			false,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			result, err := ParseDependencies(tc.path, tc.cpu)

			if (err != nil) != tc.expectError {
				t.Fatalf("ParseDependencies(%q, %d) error = %v, expectError %v", tc.path, tc.cpu, err, tc.expectError)
			}

			if len(result) == 0 && !tc.expectNone {
				t.Errorf("ParseDependencies(%q, %d) returned empty result", tc.path, tc.cpu)
			}

			for _, dep := range result {
				t.Logf("Dependency: %s", dep)
			}
		})
	}
}

// NOTE: The test samples use hypothetical paths, you should replace them with appropriate actual files when running these tests.
