package processor

import (
	"go/ast"
	"go/token"
	"strconv"
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

func translateJSONSchemaReference(json map[string]interface{}) {
	if ref, hasRef := json["$ref"].(string); len(json) == 1 && hasRef && ref[0] == '#' {
		ref = "#/components/schemas/" + ref[1:]
		json["$ref"] = ref
	} else {
		for _, value := range json {
			if subJSON, ok := value.(map[string]interface{}); ok {
				translateJSONSchemaReference(subJSON)
			}
		}
	}
}
