package scanner

import (
	"go/ast"
	"go/parser"
	"go/token"

	"golang.org/x/tools/go/packages"
)

type ScannerHandler func(*token.FileSet, *ast.File) error

type Scanner struct {
	fset    *token.FileSet
	handler ScannerHandler
}

func New(handler ScannerHandler) *Scanner {
	return &Scanner{
		fset:    token.NewFileSet(),
		handler: handler,
	}
}

func (scn *Scanner) Scan(dir string, patterns []string) error {
	pkgConfig := packages.Config{
		Dir:  dir,
		Mode: packages.NeedFiles,
	}

	pkgs, err := packages.Load(&pkgConfig, patterns...)
	if err != nil {
		return err
	}

	for _, pkg := range pkgs {
		for _, file := range pkg.GoFiles {
			astFile, err := parser.ParseFile(scn.fset, file, nil, parser.ParseComments)
			if err != nil {
				return err
			}

			err = scn.handler(scn.fset, astFile)
			if err != nil {
				return err
			}
		}
	}

	return nil
}
