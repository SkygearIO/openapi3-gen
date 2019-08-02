package openapi3

type MediaTypeObject struct {
	Schema   Schema                   `yaml:"schema,omitempty"`
	Examples map[string]ExampleObject `yaml:"examples,omitempty"`
}

func NewMediaTypeObject() *MediaTypeObject {
	return &MediaTypeObject{
		Examples: map[string]ExampleObject{},
	}
}
