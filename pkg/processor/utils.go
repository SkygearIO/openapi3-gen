package processor

import (
	"go/ast"
	"go/token"
	"strconv"

	"github.com/skygeario/openapi3-gen/pkg/openapi3"
)

const jsonMediaType = "application/json"

func extractConstValue(n ast.Node) (value string, ok bool) {
	switch typedNode := n.(type) {
	case *ast.GenDecl:
		for _, spec := range typedNode.Specs {
			value, ok = extractConstValue(spec)
			if ok {
				return
			}
		}

	case *ast.ValueSpec:
		if lit, valid := typedNode.Values[0].(*ast.BasicLit); valid && lit.Kind == token.STRING {
			litValue, err := strconv.Unquote(lit.Value)
			if err == nil {
				return litValue, true
			}
		}
	}

	return
}

func translateJSONSchema(json interface{}) interface{} {
	switch typedJSON := json.(type) {
	case map[string]interface{}:
		if ref, hasRef := typedJSON["$ref"].(string); len(typedJSON) == 1 && hasRef {
			if len(ref) > 1 && ref[0] == '#' {
				id := ref[1:]
				return map[string]interface{}(openapi3.MakeSchemaRef(id))
			}
		}
		result := map[string]interface{}{}
		for key, value := range typedJSON {
			result[key] = translateJSONSchema(value)
		}
		return result

	case []interface{}:
		result := make([]interface{}, len(typedJSON))
		for i, value := range typedJSON {
			result[i] = translateJSONSchema(value)
		}
		return result

	default:
		return json
	}
}
