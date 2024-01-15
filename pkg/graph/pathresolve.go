package graph

import (
	"fmt"
	"log"
	"os"
	"strings"
)

func getPackageRootPath(path string) (string, error) {
	// get the directory containing `info.plist`
	stack := strings.Split(path, "/")

	fileInfo, err := os.Stat(path)
	if err != nil {
		return "", err
	}
	if !fileInfo.IsDir() {
		stack = stack[:len(stack)-1]
	}

	for len(stack) > 0 {
		entries, err := os.ReadDir(strings.Join(stack, "/"))
		if err != nil {
			log.Fatalf("Failed to read directory: %v", err)
		}

		for _, entry := range entries {
			if entry.Name() == "Info.plist" {
				return strings.Join(stack, "/"), nil
			}
		}

		stack = stack[:len(stack)-1]
	}

	return "", fmt.Errorf("could not find package root for %q", path)
}

func resolvePath(path string, contextPath string) (string, error) {
	pathParts := strings.Split(path, "/")

	switch pathParts[0] {
	case "@rpath":
		// as a temporary, "good enough" solution, replace @rpath with "root package"/Frameworks/
		packageRoot, err := getPackageRootPath(contextPath)
		if err != nil {
			return "", err
		}

		return packageRoot + "/Frameworks/" + strings.Join(pathParts[1:], "/"), nil
	case "@loader_path", "@executable_path":
		return "", fmt.Errorf("path type is not supported (yet): %q", path)
	default:
		return path, nil
	}

}
