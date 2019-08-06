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

// e.g. master_key: X-API-Key in header
var apiKeySecArgFormat = regexp.MustCompile(`^([^\s]+)\s*:\s*([^\s]+)\s+in\s+([^\s]+)$`)

// e.g. access_token: JWT Bearer
var httpSecArgFormat = regexp.MustCompile(`^([^\s]+)\s*:\s*([^\s]+)(?:\s+([^\s]+))?$`)

// e.g. GET /me - Get current user
var operationArgFormat = regexp.MustCompile(`^([^\s]+)\s+([^\s]+)\s+-\s+(.+)$`)

// e.g. provider_name in query as ProviderName
var parameterArgFormat = regexp.MustCompile(`^([^\s]+)\s+in\s+([^\s]+)(?:\s+as\s+([^\s]+))?$`)

// e.g. DisableUserExpiring - Disable a user with expiry
var exampleArgFormat = regexp.MustCompile(`^([^\s]+)\s+-\s+(.+)$`)

// e.g. 200 {AuthURLResponse}
var responseRefArgFormat = regexp.MustCompile(`^([^\s]+)\s+{([^\s]+)}$`)

// e.g. user_update {UserUpdateEvent}
var callbackArgFormat = regexp.MustCompile(`^([^\s]+)\s+{([^\s]+)}$`)

var handlers map[AnnotationType]annotationHandler = map[AnnotationType]annotationHandler{
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
	AnnotationTypeSecurity: func(ctx *context, arg string, body string) error {
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
	AnnotationTypeSecurityAPIKey: func(ctx *context, arg string, body string) error {
		matches := apiKeySecArgFormat.FindStringSubmatch(arg)
		if len(matches) != 4 {
			return fmt.Errorf("must provide scheme name, parameter name and location")
		}

		name := matches[1]
		apiKeyName := matches[2]
		apiKeyLocation := openapi3.SecuritySchemeAPIKeyLocation(matches[3])
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
	AnnotationTypeSecurityHTTP: func(ctx *context, arg string, body string) error {
		matches := httpSecArgFormat.FindStringSubmatch(arg)
		if len(matches) < 3 {
			return fmt.Errorf("must provide scheme name and HTTP auth scheme")
		}

		name := matches[1]
		authScheme := strings.ToLower(matches[2])
		var bearerFormat string
		if authScheme == "bearer" {
			if len(matches) < 4 {
				return fmt.Errorf("must provide bearer token format")
			}
			bearerFormat = matches[3]
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
		matches := operationArgFormat.FindStringSubmatch(arg)
		if len(matches) < 3 {
			return fmt.Errorf("must provide HTTP method and path")
		}

		method := matches[1]
		path := matches[2]
		var summary string
		if len(matches) == 4 {
			summary = matches[3]
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
		matches := refArgFormat.FindStringSubmatch(arg)
		if len(matches) == 2 {
			if ctx.operation == nil {
				return fmt.Errorf("must be used with Operation")
			}
			id := matches[1]
			ctx.operation.Parameters = append(ctx.operation.Parameters, openapi3.MakeParameterRef(id))
			return nil
		}

		matches = parameterArgFormat.FindStringSubmatch(arg)
		if len(matches) < 3 {
			return fmt.Errorf("must provide parameter name and location")
		}

		name := matches[1]
		location := openapi3.ParameterLocation(matches[2])
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
			if len(matches) < 4 {
				return fmt.Errorf("must provide reference ID")
			}
			id := matches[3]
			ctx.oapi.Components.Parameters[id] = parameter
		}

		ctx.setContextObject(parameter)

		return nil
	},
	AnnotationTypeRequestBody: func(ctx *context, arg string, body string) error {
		matches := refArgFormat.FindStringSubmatch(arg)
		if len(matches) == 2 {
			if ctx.operation == nil {
				return fmt.Errorf("must be used with Operation")
			}
			id := matches[1]
			ctx.operation.RequestBody = openapi3.MakeRequestBodyRef(id)
			return nil
		}

		requestBody := openapi3.NewRequestBodyObject()
		requestBody.Description = body
		if ctx.operation != nil {
			ctx.operation.RequestBody = requestBody
		} else {
			id := arg
			ctx.oapi.Components.RequestBodies[id] = requestBody
		}

		ctx.setContextObject(requestBody)

		return nil
	},
	AnnotationTypeResponse: func(ctx *context, arg string, body string) error {
		if ctx.operation == nil {
			response := openapi3.NewResponseObject()
			response.Description = body
			ctx.oapi.Components.Responses[arg] = response
			ctx.response = response
		} else {
			var response openapi3.Response
			var responseKey string
			matches := responseRefArgFormat.FindStringSubmatch(arg)
			if len(matches) == 3 {
				responseKey = matches[1]
				id := matches[2]
				response = openapi3.MakeResponseRef(id)
			} else {
				responseObj := openapi3.NewResponseObject()
				responseObj.Description = body
				responseKey = arg
				response = responseObj
				ctx.setContextObject(responseObj)
			}

			ctx.operation.Responses[responseKey] = response
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
			valid := len(body) > 0
			if !valid {
				schemaValue, valid = extractConstValue(ctx.node)
			}
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

		matches := exampleArgFormat.FindStringSubmatch(arg)
		if len(matches) != 3 {
			return fmt.Errorf("must provide example name and summary")
		}
		name := matches[1]
		summary := matches[2]

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
			callback := openapi3.NewCallbackObject()
			ctx.oapi.Components.Callbacks[arg] = callback
			ctx.setContextObject(callback)
		} else {
			var callback openapi3.Callback
			var callbackKey string
			matches := callbackArgFormat.FindStringSubmatch(arg)
			if len(matches) == 3 {
				callbackKey = matches[1]
				id := matches[2]
				callback = openapi3.MakeCallbackRef(id)
			} else {
				callback = openapi3.NewCallbackObject()
				callbackKey = arg
				ctx.setContextObject(callback)
			}

			ctx.operation.Callbacks[callbackKey] = callback
		}

		return nil
	},
}
