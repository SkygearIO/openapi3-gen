package openapi3

type ServerObject struct {
	URL         string                    `yaml:"url"`
	Description string                    `yaml:"description,omitempty"`
	Variables   map[string]ServerVariable `yaml:"variables,omitempty"`
}

type ServerVariable struct {
	Enum        []string `yaml:"enum,omitempty"`
	Default     string   `yaml:"default"`
	Description string   `yaml:"description,omitempty"`
}

func NewServerObject() *ServerObject {
	return &ServerObject{
		Variables: map[string]ServerVariable{},
	}
}
