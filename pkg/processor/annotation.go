package processor

import (
	"regexp"
	"strings"
)

//go:generate stringer -type=AnnotationType -trimprefix=AnnotationType
type AnnotationType int

const (
	// @ID <Component ID>
	// e.g.
	// @ID EventContext
	AnnotationTypeID AnnotationType = iota

	// @API <Title>
	// [<Description>]
	// e.g.
	// @API Test API
	//     This is a test API.
	AnnotationTypeAPI

	// @Version <Version>
	// e.g.
	// @Version 1.0.0
	AnnotationTypeVersion

	// @Server <URL>
	// [<Description>]
	// e.g.
	// @Server https://{env}.example.com
	//     Internal API Server
	AnnotationTypeServer

	// @Variable <Name> <Default Value> [<Possible Value>...]
	// [<Description>]
	// e.g.
	// @Variable env staging dev qa staging
	AnnotationTypeVariable

	// @Tag <Name>
	// [<Description>]
	// e.g.
	// @Tag User
	//     User APIs
	AnnotationTypeTag

	// @SecurityRequirement <Security Scheme ID>
	// e.g.
	// @SecurityRequirement access_token
	AnnotationTypeSecurityRequirement

	// @SecuritySchemeAPIKey <ID> <Field Name> <Field Location>
	// [<Description>]
	// e.g.
	// @SecuritySchemeAPIKey api_key X-API-Key header
	AnnotationTypeSecuritySchemeAPIKey

	// @SecuritySchemeHTTP <ID> <HTTP Auth Scheme> [<Bearer Token Format>]
	// [<Description>]
	// e.g.
	// @SecuritySchemeHTTP access_token Bearer JWT
	//     Access Token
	AnnotationTypeSecuritySchemeHTTP

	// @Operation <HTTP Method> <Path> - <Summary>
	// [<Description>]
	// e.g.
	// @Operation GET /user/{id} - Get User
	//     Return the user with specified ID.
	AnnotationTypeOperation

	// @Parameter [<Name> <Location>|{<Component ID>}]
	// [<Description>]
	// e.g.
	// @Parameter id path
	//     ID of user
	// @Parameter {UserID}
	AnnotationTypeParameter

	// @RequestBody [{<Component ID>}]
	// [<Description>]
	// e.g.
	// @RequestBody
	//     New user information
	// @RequestBody {UpdateUserRequest}
	AnnotationTypeRequestBody

	// @Response <Status Code> [{<Component ID>}]
	// [<Description>]
	// e.g.
	// @Response 200
	//     User is updated successfully
	// @Response default {ErrorResponse}
	AnnotationTypeResponse

	// @JSONSchema [{<Component ID>}]
	// <JSON Schema>
	// e.g.
	// @JSONSchema
	//     { "type": "object" }
	// @JSONSchema {User}
	AnnotationTypeJSONSchema

	// @JSONExample <Key> - <Summary>
	// <JSON>
	// e.g.
	// @JSONExample TestUser - Test user information
	//     { "id": "user-id" }
	AnnotationTypeJSONExample

	// @Callback <Key>
	// e.g.
	// @Callback UserCreated
	AnnotationTypeCallback

	AnnotationTypeMaximum
)

var annotationTypeMap map[string]AnnotationType

func init() {
	annotationTypeMap = map[string]AnnotationType{}
	for i := 0; i < int(AnnotationTypeMaximum); i++ {
		t := AnnotationType(i)
		annotationTypeMap[strings.ToLower(t.String())] = t
	}
}

type Annotation struct {
	Type     AnnotationType
	Argument string
	Body     []string
}

var annotationRegex = regexp.MustCompile(`^@([^\s]+)(?:\s+(.*))?$`)

func tryParseAnnotation(line string) (annotation Annotation, ok bool) {
	matches := annotationRegex.FindStringSubmatch(line)
	if len(matches) != 3 {
		return
	}

	t, validType := annotationTypeMap[strings.ToLower(matches[1])]
	if !validType {
		return
	}

	annotation = Annotation{Type: t, Argument: matches[2]}
	ok = true
	return
}

func trimEmptyLines(body []string) []string {
	start := 0
	for start < len(body) && body[start] == "" {
		start++
	}
	end := len(body) - 1
	for end >= start && body[end] == "" {
		end--
	}
	return body[start : end+1]
}

func ParseAnnotations(lines []string) []Annotation {
	var annotations []Annotation
	var current *Annotation

	for _, line := range lines {
		line = strings.TrimSpace(line)
		annotation, ok := tryParseAnnotation(line)
		if ok {
			if current != nil {
				current.Body = trimEmptyLines(current.Body)
				annotations = append(annotations, *current)
			}
			current = &annotation
		} else if current != nil {
			current.Body = append(current.Body, line)
		}
	}
	if current != nil {
		current.Body = trimEmptyLines(current.Body)
		annotations = append(annotations, *current)
	}
	return annotations
}
