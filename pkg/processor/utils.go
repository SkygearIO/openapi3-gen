package processor

import (
	"go/ast"
	"go/token"
	"regexp"
	"strconv"

	"github.com/skygeario/openapi3-gen/pkg/openapi3"
)

const jsonMediaType = "application/json"

func extractDeclName(n ast.Node) (name string, ok bool) {
	switch typedNode := n.(type) {
	case *ast.FuncDecl:
		return typedNode.Name.Name, true

	case *ast.GenDecl:
		for _, spec := range typedNode.Specs {
			name, ok = extractDeclName(spec)
			if ok {
				return
			}
		}

	case *ast.ValueSpec:
		return typedNode.Names[0].Name, true

	case *ast.TypeSpec:
		return typedNode.Name.Name, true
	}

	return
}

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
		if len(typedNode.Values) == 0 {
			return
		}
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

func matchRegex(str string, re *regexp.Regexp) (matches []string, success bool) {
	matches = re.FindStringSubmatch(str)
	if len(matches) == 0 {
		success = false
		return
	} else {
		matches = matches[1:]
		success = true
		return
	}
}
