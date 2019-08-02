package openapi3

type ParameterLocation string

const (
	ParameterLocationQuery  = "query"
	ParameterLocationHeader = "header"
	ParameterLocationPath   = "path"
	ParameterLocationCookie = "cookie"
)

func (l ParameterLocation) Validate() bool {
	return l == ParameterLocationQuery ||
		l == ParameterLocationHeader ||
		l == ParameterLocationPath ||
		l == ParameterLocationCookie
}

type Parameter interface{}
type ParameterObject struct {
	Name        string                   `yaml:"name"`
	Location    ParameterLocation        `yaml:"in"`
	Description string                   `yaml:"description,omitempty"`
	Required    bool                     `yaml:"required,omitempty"`
	Schema      Schema                   `yaml:"schema,omitempty"`
	Examples    map[string]ExampleObject `yaml:"examples,omitempty"`
}

func NewParameterObject() *ParameterObject {
	return &ParameterObject{
		Examples: map[string]ExampleObject{},
	}
}
