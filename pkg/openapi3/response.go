package openapi3

type Response interface{}
type ResponseObject struct {
	Description string                     `yaml:"description"`
	Content     map[string]MediaTypeObject `yaml:"content,omitempty"`
}

func NewResponseObject() *ResponseObject {
	return &ResponseObject{
		Content: map[string]MediaTypeObject{},
	}
}
