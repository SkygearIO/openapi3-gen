package processor

import (
	"regexp"
	"strings"
)

//go:generate stringer -type=AnnotationType -trimprefix=AnnotationType
type AnnotationType int

const (
	AnnotationTypeAPI AnnotationType = iota
	AnnotationTypeVersion
	AnnotationTypeServer
	AnnotationTypeVariable
	AnnotationTypeTag
	AnnotationTypeSecurity
	AnnotationTypeSecurityAPIKey
	AnnotationTypeSecurityHTTP
	AnnotationTypeOperation
	AnnotationTypeParameter
	AnnotationTypeRequestBody
	AnnotationTypeResponse
	AnnotationTypeJSONSchema
	AnnotationTypeJSONExample
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

func ParseAnnotations(lines []string) []Annotation {
	var annotations []Annotation
	var current *Annotation

	for _, line := range lines {
		line = strings.TrimSpace(line)
		annotation, ok := tryParseAnnotation(line)
		if ok {
			if current != nil {
				annotations = append(annotations, *current)
			}
			current = &annotation
		} else if current != nil {
			current.Body = append(current.Body, line)
		}
	}
	if current != nil {
		annotations = append(annotations, *current)
	}
	return annotations
}
