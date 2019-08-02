package internal

import (
	"fmt"
	"go/ast"
	"go/token"

	"github.com/skygeario/openapi3-gen/pkg/scanner"
)

func Run(baseDir string, patterns []string, outputFile string) error {
	scn := scanner.New(func(fset *token.FileSet, file *ast.File) error {
		fmt.Println(file.Name)
		return nil
	})

	err := scn.Scan(baseDir, patterns)
	if err != nil {
		return err
	}
	return nil
}
