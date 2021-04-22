package main

import (
	"fmt"
	"go/parser"
	"go/token"
	"io"
	"log"
	"os"
	"strconv"
	"strings"
)

func main() {
	input, err := io.ReadAll(os.Stdin)
	if err != nil {
		log.Fatal(err)
	}

	for _, filename := range strings.Split(string(input), "\n") {
		if filename == "" {
			continue
		}

		fset := token.NewFileSet()
		file, err := parser.ParseFile(fset, filename, nil, parser.ImportsOnly)
		if err != nil {
			log.Fatal(err)
		}

		var runs [][]string
		var run []string
		for j, spec := range file.Imports {
			path, err := strconv.Unquote(spec.Path.Value)
			if err != nil {
				log.Fatal(err)
			}

			if j > 0 && lineAt(fset, spec.Pos()) > 1+lineAt(fset, file.Imports[j-1].End()) {
				// j begins a new run. End this one.
				runs = append(runs, run)
				run = nil
			}
			run = append(run, path)
		}
		if run != nil {
			runs = append(runs, run)
		}

		fmt.Println(classify(runs))
	}
}

func lineAt(fset *token.FileSet, pos token.Pos) int {
	return fset.PositionFor(pos, false).Line
}

func classify(runs [][]string) string {
	var std, cmd bool

	for _, run := range runs {
		for _, path := range run {
			if strings.HasPrefix(path, "cmd/") {
				cmd = true
			} else {
				if cmd {
					return "std+cmd mixed"
				}
				std = true
			}
		}
	}

	if std && cmd {
		if len(runs) <= 1 {
			return "std+cmd mixed"
		}

		return "std+cmd split"
	}

	if std {
		return "std only"
	}
	if cmd {
		return "cmd only"
	}

	return "none"
}
