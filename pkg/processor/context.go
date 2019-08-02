package processor

import (
	"go/ast"

	"github.com/skygeario/openapi3-gen/pkg/openapi3"
)

type context struct {
	oapi *openapi3.OpenAPIObject
	node ast.Node

	server      *openapi3.ServerObject
	operation   *openapi3.OperationObject
	parameter   *openapi3.ParameterObject
	requestBody *openapi3.RequestBodyObject
	response    *openapi3.ResponseObject
	callback    *openapi3.CallbackObject
}

func newContext(oapi *openapi3.OpenAPIObject, node ast.Node) *context {
	return &context{oapi: oapi, node: node}
}

func (ctx *context) Consume(annotation Annotation) error {
	// TODO: consume annotation and update context
	return nil
}
