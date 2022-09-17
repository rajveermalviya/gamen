package main

import (
	"fmt"
	"io/fs"
	"os"
	"path/filepath"

	"github.com/iancoleman/strcase"
)

var packagePrefix = os.Args[1]

func main() {
	packagePaths := []string{}

	err := filepath.WalkDir(".", func(path string, d fs.DirEntry, err error) error {
		if err != nil || !d.IsDir() || path == "." {
			return err
		}

		packagePaths = append(packagePaths, packagePrefix+path)
		packageName := strcase.ToSnake(filepath.Base(path))

		f, err := os.Create(filepath.Join(path, "dummy.go"))
		if err != nil {
			return err
		}
		defer f.Close()
		fmt.Fprintf(f, "//go:build dummy\n\n")
		fmt.Fprintf(f, "package %s\n", packageName)
		return nil
	})
	if err != nil {
		panic(err)
	}

	f, err := os.Create("dummy.go")
	if err != nil {
		panic(err)
	}
	defer f.Close()

	fmt.Fprintf(f, "//go:build dummy\n\n")
	fmt.Fprintf(f, "package include\n\n")
	fmt.Fprintf(f, "import (\n")
	for _, v := range packagePaths {
		fmt.Fprintf(f, "\t_ \"%s\"\n", v)
	}
	fmt.Fprintf(f, ")\n")
}
