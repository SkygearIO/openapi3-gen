package openapi3

type RequestBody interface{}
type RequestBodyObject struct {
	Description string                     `yaml:"description,omitempty"`
	Content     map[string]MediaTypeObject `yaml:"content"`
	Required    bool                       `yaml:"required,omitempty"`
}

func NewRequestBodyObject() *RequestBodyObject {
	return &RequestBodyObject{
		Content: map[string]MediaTypeObject{},
	}
}
