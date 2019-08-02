package processor

type annotationHandler func(ctx *context, arg string, body string) error

var handlers map[AnnotationType]annotationHandler = map[AnnotationType]annotationHandler{
	AnnotationTypeAPI: func(ctx *context, arg string, body string) error {
		return nil
	},
	AnnotationTypeVersion: func(ctx *context, arg string, body string) error {
		return nil
	},
	AnnotationTypeServer: func(ctx *context, arg string, body string) error {
		return nil
	},
	AnnotationTypeVariable: func(ctx *context, arg string, body string) error {
		return nil
	},
	AnnotationTypeTag: func(ctx *context, arg string, body string) error {
		return nil
	},
	AnnotationTypeSecurity: func(ctx *context, arg string, body string) error {
		return nil
	},
	AnnotationTypeSecurityAPIKey: func(ctx *context, arg string, body string) error {
		return nil
	},
	AnnotationTypeSecurityHTTP: func(ctx *context, arg string, body string) error {
		return nil
	},
	AnnotationTypeOperation: func(ctx *context, arg string, body string) error {
		return nil
	},
	AnnotationTypeParameter: func(ctx *context, arg string, body string) error {
		return nil
	},
	AnnotationTypeRequestBody: func(ctx *context, arg string, body string) error {
		return nil
	},
	AnnotationTypeResponse: func(ctx *context, arg string, body string) error {
		return nil
	},
	AnnotationTypeJSONSchema: func(ctx *context, arg string, body string) error {
		return nil
	},
	AnnotationTypeJSONExample: func(ctx *context, arg string, body string) error {
		return nil
	},
	AnnotationTypeCallback: func(ctx *context, arg string, body string) error {
		return nil
	},
}
