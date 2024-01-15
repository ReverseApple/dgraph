package graph

import (
	"errors"
	"fmt"
	"github.com/blacktop/go-macho"
	"github.com/blacktop/go-macho/types"
	"path/filepath"
	"strings"
)

func systemBinaryResidesInDyldCache(path string) bool {
	// this is a pretty dumb classification function, but it's fine for now.

	if strings.HasPrefix(path, "/System/Library/PrivateFrameworks") ||
		strings.HasPrefix(path, "/System/Library/Frameworks") ||
		strings.HasPrefix(path, "/System/iOSSupport/System/Library/PrivateFrameworks") ||
		strings.HasPrefix(path, "/usr/lib/") {

		_, file := filepath.Split(path)

		switch file {
		case "CoreFP", "CoreKE", "CoreLSKD", "CoreADI":
			return false
		default:
			return true
		}
	}

	return false
}

func isPathLoadEntry(load macho.Load) bool {
	loadParts := strings.Split(load.String(), "/")

	switch strings.TrimSpace(loadParts[0]) {
	case "@executable_path", "@loader_path":
		return false
	case "@rpath":
		return true
	default:
		return strings.HasPrefix(load.String(), "/")
	}
}

func getDependencies(file *macho.File) []string {
	var result []string

	for _, load := range file.Loads {
		if isPathLoadEntry(load) {
			loadSpl := strings.Split(load.String(), " ")
			var loadFixed string
			if len(loadSpl) == 2 {
				loadFixed = strings.TrimSpace(loadSpl[:len(loadSpl)-1][0])
			} else {
				loadFixed = load.String()
			}

			result = append(result, loadFixed)
		}
	}

	return result
}

func ParseDependencies(path string, cpu types.CPU) ([]string, error) {

	if systemBinaryResidesInDyldCache(path) {

		// todo: implement this
		return []string{}, nil
	} else {
		//fmt.Printf("parsing: %q\n", path)
		// if we aren't in the dyld cache, try to get the correct binary architecture for parsing.
		// reference: https://github.com/blacktop/ipsw/blob/1805c3a59e21aeac55e1da44c14d0fb98a8e005d/cmd/ipsw/cmd/macho/macho_info.go#L233C30-L233C30
		fat, err := macho.OpenFat(path)
		if err != nil && !errors.Is(err, macho.ErrNotFat) {
			return nil, err
		}

		if errors.Is(err, macho.ErrNotFat) {
			file, err := macho.Open(path)
			if err != nil {
				return nil, err
			}
			defer file.Close()

			if file.CPU != cpu {
				return nil, fmt.Errorf("parsed CPU arch for %q (%q) does not meet requirement arch: (%q)", path, file.CPU.String(), cpu.String())
			}

			return getDependencies(file), nil
		}

		archs := fat.Arches
		for _, arch := range archs {
			if arch.CPU == cpu {
				file := arch.File
				return getDependencies(file), nil
			}
		}

		return nil, fmt.Errorf("no matching architecture found for %q", path)
	}

}
