package processor

import (
	"fmt"
	"go/ast"
	"strings"

	"github.com/skygeario/openapi3-gen/pkg/openapi3"
)

type context struct {
	astNodeName  string
	astNodeValue string
	componentID  string

	oapi        *openapi3.OpenAPIObject
	server      *openapi3.ServerObject
	operation   *openapi3.OperationObject
	parameter   *openapi3.ParameterObject
	requestBody *openapi3.RequestBodyObject
	response    *openapi3.ResponseObject
	callback    *openapi3.CallbackObject
}

func newContext(oapi *openapi3.OpenAPIObject, node ast.Node) *context {
	name, _ := extractDeclName(node)
	value, _ := extractConstValue(node)
	return &context{
		astNodeName:  name,
		astNodeValue: value,
		componentID:  name,
		oapi:         oapi,
	}
}

func (ctx *context) setContextObject(scope interface{}) {
	switch obj := scope.(type) {
	case *openapi3.ServerObject:
		ctx.server = obj
	case *openapi3.OperationObject:
		ctx.operation = obj
	case *openapi3.ParameterObject:
		ctx.parameter = obj
		ctx.requestBody = nil
		ctx.response = nil
	case *openapi3.RequestBodyObject:
		ctx.parameter = nil
		ctx.requestBody = obj
		ctx.response = nil
	case *openapi3.ResponseObject:
		ctx.parameter = nil
		ctx.requestBody = nil
		ctx.response = obj
	case *openapi3.CallbackObject:
		ctx.parameter = nil
		ctx.requestBody = nil
		ctx.response = nil
		ctx.callback = obj
	default:
		panic(fmt.Errorf("unknown contextual object: %T", scope))
	}
}

func (ctx *context) Consume(annotation Annotation) error {
	handler, exists := handlers[annotation.Type]
	if !exists {
		panic(fmt.Errorf("unknown annotation type: %v", annotation.Type))
	}

	body := strings.Join(annotation.Body, "\n")
	return handler(ctx, annotation.Argument, body)
}
