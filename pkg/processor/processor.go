package processor

import (
	"go/ast"
	"go/token"
	"strings"

	"github.com/skygeario/openapi3-gen/pkg/openapi3"
)

type Processor struct {
	oapi *openapi3.OpenAPIObject
	errs []error
}

func New() *Processor {
	return &Processor{
		oapi: openapi3.NewOpenAPIObject(),
	}
}

func (psr *Processor) End() (*openapi3.OpenAPIObject, []error) {
	return psr.oapi, psr.errs
}

func (psr *Processor) Process(fset *token.FileSet, file *ast.File) {
	ast.Inspect(file, func(n ast.Node) bool {
		if decl, ok := n.(*ast.FuncDecl); ok {
			psr.processNode(fset, n, decl.Doc)
		} else if decl, ok := n.(*ast.GenDecl); ok {
			psr.processNode(fset, n, decl.Doc)
		}
		return true
	})
}

func (psr *Processor) processNode(fset *token.FileSet, node ast.Node, doc *ast.CommentGroup) {
	docLines := strings.Split(doc.Text(), "\n")
	_ = docLines
	// TODO: parse & process annotations
}
