package processor

import (
	"encoding/json"
	"fmt"
	"regexp"
	"strings"

	"github.com/pkg/errors"
	"github.com/skygeario/openapi3-gen/pkg/openapi3"
)

type annotationHandler func(ctx *context, arg string, body string) error

// e.g. {EmptyResponse}
var refArgFormat = regexp.MustCompile(`^{([^\s]+)}$`)

// e.g. GET /me - Get current user
var operationArgFormat = regexp.MustCompile(`^([^\s]+)\s+([^\s]+)\s+-\s+(.+)$`)

// e.g. DisableUserExpiring - Disable a user with expiry
var exampleArgFormat = regexp.MustCompile(`^([^\s]+)\s+-\s+(.+)$`)

var handlers map[AnnotationType]annotationHandler = map[AnnotationType]annotationHandler{
	AnnotationTypeID: func(ctx *context, arg string, body string) error {
		ctx.componentID = arg
		return nil
	},
	AnnotationTypeAPI: func(ctx *context, arg string, body string) error {
		ctx.oapi.Info.Title = arg
		ctx.oapi.Info.Description = body
		return nil
	},
	AnnotationTypeVersion: func(ctx *context, arg string, body string) error {
		ctx.oapi.Info.Version = arg
		return nil
	},
	AnnotationTypeServer: func(ctx *context, arg string, body string) error {
		server := openapi3.NewServerObject()
		server.URL = arg
		server.Description = body

		servers := append(ctx.oapi.Servers, *server)
		ctx.oapi.Servers = servers
		ctx.setContextObject(&servers[len(servers)-1])
		return nil
	},
	AnnotationTypeVariable: func(ctx *context, arg string, body string) error {
		if ctx.server == nil {
			return fmt.Errorf("must be used with Server")
		}

		parts := strings.Fields(arg)
		if len(parts) < 2 {
			return fmt.Errorf("must provide name and at least one value")
		}

		variable := openapi3.ServerVariable{
			Description: body,
			Default:     parts[1],
			Enum:        parts[2:],
		}
		ctx.server.Variables[parts[0]] = variable
		return nil
	},
	AnnotationTypeTag: func(ctx *context, arg string, body string) error {
		if len(body) == 0 {
			if ctx.operation == nil {
				return fmt.Errorf("must be used with Operation")
			}
			ctx.operation.Tags = append(ctx.operation.Tags, arg)
		} else {
			tag := openapi3.TagObject{
				Name:        arg,
				Description: body,
			}
			ctx.oapi.Tags = append(ctx.oapi.Tags, tag)
		}
		return nil
	},
	AnnotationTypeSecurityRequirement: func(ctx *context, arg string, body string) error {
		args := strings.Fields(arg)
		if len(args) < 1 {
			return fmt.Errorf("must provide security scheme name")
		}
		requirement := openapi3.SecurityRequirementObject{
			args[0]: args[1:],
		}

		if ctx.operation != nil {
			ctx.operation.Security = append(ctx.operation.Security, requirement)
		} else {
			ctx.oapi.Security = append(ctx.oapi.Security, requirement)
		}
		return nil
	},
	AnnotationTypeSecuritySchemeAPIKey: func(ctx *context, arg string, body string) error {
		fields := strings.Fields(arg)
		if len(fields) != 3 {
			return fmt.Errorf("must provide scheme name, parameter name and location")
		}

		name := fields[0]
		apiKeyLocation := openapi3.SecuritySchemeAPIKeyLocation(fields[1])
		apiKeyName := fields[2]
		if !apiKeyLocation.Validate() {
			return fmt.Errorf("invalid API key location: %v", apiKeyLocation)
		}

		scheme := &openapi3.SecuritySchemeObject{
			Type:           openapi3.SecuritySchemeTypeAPIKey,
			APIKeyName:     apiKeyName,
			APIKeyLocation: apiKeyLocation,
			Description:    body,
		}
		ctx.oapi.Components.SecuritySchemes[name] = scheme
		return nil
	},
	AnnotationTypeSecuritySchemeHTTP: func(ctx *context, arg string, body string) error {
		fields := strings.Fields(arg)
		if len(fields) < 2 {
			return fmt.Errorf("must provide scheme name and HTTP auth scheme")
		}

		name := fields[0]
		authScheme := strings.ToLower(fields[1])
		var bearerFormat string
		if authScheme == "bearer" {
			if len(fields) < 3 {
				return fmt.Errorf("must provide bearer token format")
			}
			bearerFormat = fields[2]
		}

		scheme := &openapi3.SecuritySchemeObject{
			Type:             openapi3.SecuritySchemeTypeHTTP,
			HTTPAuthScheme:   authScheme,
			HTTPBearerFormat: bearerFormat,
			Description:      body,
		}
		ctx.oapi.Components.SecuritySchemes[name] = scheme
		return nil
	},
	AnnotationTypeOperation: func(ctx *context, arg string, body string) error {
		matches, success := matchRegex(arg, operationArgFormat)
		if !success {
			return fmt.Errorf("must provide HTTP method and path")
		}

		method := matches[0]
		path := matches[1]
		var summary string
		if len(matches) == 3 {
			summary = matches[2]
		}

		operation := openapi3.NewOperationObject()
		operation.Summary = summary
		operation.Description = body

		var paths openapi3.Paths
		if ctx.callback != nil {
			paths = ctx.callback
		} else {
			paths = &ctx.oapi.Paths
		}

		pathItem := paths.GetPath(path)
		if !pathItem.SetOperation(method, operation) {
			return fmt.Errorf("invalid HTTP method: %v", method)
		}

		paths.SetPath(path, pathItem)
		ctx.setContextObject(operation)

		return nil
	},
	AnnotationTypeParameter: func(ctx *context, arg string, body string) error {
		matches, isRef := matchRegex(arg, refArgFormat)
		if isRef {
			if ctx.operation == nil {
				return fmt.Errorf("must be used with Operation")
			}
			id := matches[0]
			ctx.operation.Parameters = append(ctx.operation.Parameters, openapi3.MakeParameterRef(id))
			return nil
		}

		fields := strings.Fields(arg)
		if len(fields) != 2 {
			return fmt.Errorf("must provide parameter name and location")
		}

		name := fields[0]
		location := openapi3.ParameterLocation(fields[1])
		if !location.Validate() {
			return fmt.Errorf("invalid parameter location: %v", location)
		}

		parameter := openapi3.NewParameterObject()
		parameter.Name = name
		parameter.Location = location
		parameter.Required = location == openapi3.ParameterLocationPath
		parameter.Description = body

		if ctx.operation != nil {
			ctx.operation.Parameters = append(ctx.operation.Parameters, parameter)
		} else {
			if ctx.componentID == "" {
				return fmt.Errorf("must provide component ID")
			}
			ctx.oapi.Components.Parameters[ctx.componentID] = parameter
			ctx.componentID = ""
		}

		ctx.setContextObject(parameter)

		return nil
	},
	AnnotationTypeRequestBody: func(ctx *context, arg string, body string) error {
		matches, isRef := matchRegex(arg, refArgFormat)
		if isRef {
			if ctx.operation == nil {
				return fmt.Errorf("must be used with Operation")
			}
			id := matches[0]
			ctx.operation.RequestBody = openapi3.MakeRequestBodyRef(id)
			return nil
		}

		requestBody := openapi3.NewRequestBodyObject()
		requestBody.Description = body
		if ctx.operation != nil {
			ctx.operation.RequestBody = requestBody
		} else {
			if ctx.componentID == "" {
				return fmt.Errorf("must provide component ID")
			}
			ctx.oapi.Components.RequestBodies[ctx.componentID] = requestBody
			ctx.componentID = ""
		}

		ctx.setContextObject(requestBody)

		return nil
	},
	AnnotationTypeResponse: func(ctx *context, arg string, body string) error {
		if ctx.operation == nil {
			response := openapi3.NewResponseObject()
			response.Description = body

			if ctx.componentID == "" {
				return fmt.Errorf("must provide component ID")
			}
			ctx.oapi.Components.Responses[ctx.componentID] = response
			ctx.componentID = ""

			ctx.setContextObject(response)
		} else {
			var response openapi3.Response
			var statusCode string
			fields := strings.Fields(arg)
			switch len(fields) {
			case 1:
				responseObj := openapi3.NewResponseObject()
				responseObj.Description = body
				statusCode = fields[0]
				response = responseObj
				ctx.setContextObject(responseObj)
			case 2:
				statusCode = fields[0]
				matches, success := matchRegex(fields[1], refArgFormat)
				if !success {
					return fmt.Errorf("invalid object reference format")
				}
				id := matches[0]
				response = openapi3.MakeResponseRef(id)
			default:
				return fmt.Errorf("invalid response annotation format")
			}

			ctx.operation.Responses[statusCode] = response
		}

		return nil
	},
	AnnotationTypeJSONSchema: func(ctx *context, arg string, body string) error {
		var schema openapi3.Schema
		var id string
		matches := refArgFormat.FindStringSubmatch(arg)
		isRef := len(matches) == 2
		if isRef {
			id = matches[1]
			schema = openapi3.MakeSchemaRef(id)
		} else {
			schemaValue := body
			if len(schemaValue) == 0 {
				schemaValue = ctx.astNodeValue
			}

			valid := len(schemaValue) > 0
			if !valid {
				return fmt.Errorf("invalid json schema declaration")
			}

			var jsonSchema map[string]interface{}
			err := json.Unmarshal([]byte(schemaValue), &jsonSchema)
			if err != nil {
				return errors.Wrap(err, "invalid json schema")
			}

			jsonSchema = translateJSONSchema(jsonSchema).(map[string]interface{})

			id, _ = jsonSchema["$id"].(string)
			if len(id) > 0 {
				if id[0] != '#' {
					return fmt.Errorf("json schema ID must start with #")
				}
				id = id[1:]
			}
			delete(jsonSchema, "$id")

			schema = jsonSchema
		}

		if ctx.parameter != nil {
			ctx.parameter.Schema = schema
		} else if ctx.requestBody != nil {
			mediaType, exists := ctx.requestBody.Content[jsonMediaType]
			if !exists {
				mediaType = *openapi3.NewMediaTypeObject()
			}
			mediaType.Schema = schema
			ctx.requestBody.Content[jsonMediaType] = mediaType
		} else if ctx.response != nil {
			mediaType, exists := ctx.response.Content[jsonMediaType]
			if !exists {
				mediaType = *openapi3.NewMediaTypeObject()
			}
			mediaType.Schema = schema
			ctx.response.Content[jsonMediaType] = mediaType
		} else {
			if isRef {
				return fmt.Errorf("invalid annotation usage")
			}
			if id == "" {
				return fmt.Errorf("schema must contains non-empty top-level '$id' property")
			}
			ctx.oapi.Components.Schemas[id] = &schema
		}

		return nil
	},
	AnnotationTypeJSONExample: func(ctx *context, arg string, body string) error {
		var value interface{}

		err := json.Unmarshal([]byte(body), &value)
		if err != nil {
			return errors.Wrap(err, "invalid json example")
		}

		matches, success := matchRegex(arg, exampleArgFormat)
		if !success {
			return fmt.Errorf("must provide example name and summary")
		}
		name := matches[0]
		summary := matches[1]

		example := openapi3.ExampleObject{
			Summary: summary,
			Value:   value,
		}
		if ctx.parameter != nil {
			ctx.parameter.Examples[name] = example
		} else if ctx.requestBody != nil {
			mediaType, exists := ctx.requestBody.Content[jsonMediaType]
			if !exists {
				mediaType = *openapi3.NewMediaTypeObject()
			}
			mediaType.Examples[name] = example
			ctx.requestBody.Content[jsonMediaType] = mediaType
		} else if ctx.response != nil {
			mediaType, exists := ctx.response.Content[jsonMediaType]
			if !exists {
				mediaType = *openapi3.NewMediaTypeObject()
			}
			mediaType.Examples[name] = example
			ctx.response.Content[jsonMediaType] = mediaType
		} else {
			return fmt.Errorf("invalid annotation usage")
		}

		return nil
	},
	AnnotationTypeCallback: func(ctx *context, arg string, body string) error {
		if ctx.operation == nil {
			if ctx.componentID == "" {
				return fmt.Errorf("must provide component ID")
			}

			callback := openapi3.NewCallbackObject()
			ctx.oapi.Components.Callbacks[ctx.componentID] = callback
			ctx.componentID = ""
			ctx.setContextObject(callback)
		} else {
			var callback openapi3.Callback
			var callbackKey string
			fields := strings.Fields(arg)
			switch len(fields) {
			case 1:
				callback = openapi3.NewCallbackObject()
				callbackKey = fields[0]
				ctx.setContextObject(callback)
			case 2:
				callbackKey = fields[0]
				matches, success := matchRegex(fields[1], refArgFormat)
				if !success {
					return fmt.Errorf("invalid object reference format")
				}
				id := matches[0]
				callback = openapi3.MakeCallbackRef(id)
			default:
				return fmt.Errorf("invalid callback annotation format")
			}

			ctx.operation.Callbacks[callbackKey] = callback
		}

		return nil
	},
}
