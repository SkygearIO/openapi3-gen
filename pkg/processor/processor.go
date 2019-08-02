package processor

import (
	"fmt"
	"go/ast"
	"go/token"
	"strings"

	"github.com/skygeario/openapi3-gen/pkg/openapi3"
)

type processorError struct {
	inner    error
	position token.Position
}

func (err processorError) Unwrap() error {
	return err.inner
}

func (err processorError) Error() string {
	return fmt.Sprintf("%v: %v", err.position, err.inner)
}

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
	annotations := ParseAnnotations(docLines)

	errs := psr.processAnnotations(node, annotations)
	for _, err := range errs {
		err = processorError{inner: err, position: fset.Position(node.Pos())}
		psr.errs = append(psr.errs, err)
	}
}

func (psr *Processor) processAnnotations(node ast.Node, annotations []Annotation) (errs []error) {
	ctx := newContext(psr.oapi, node)
	for _, annotation := range annotations {
		err := ctx.Consume(annotation)
		if err != nil {
			errs = append(errs, err)
		}
	}
	return
}
