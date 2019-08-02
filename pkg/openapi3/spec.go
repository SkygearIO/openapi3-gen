package openapi3

type OpenAPIObject struct {
	Version    string                      `yaml:"openapi"`
	Info       InfoObject                  `yaml:"info"`
	Servers    []ServerObject              `yaml:"servers,omitempty"`
	Paths      PathsObject                 `yaml:"paths"`
	Components ComponentsObject            `yaml:"components,omitempty"`
	Security   []SecurityRequirementObject `yaml:"security,omitempty"`
	Tags       []TagObject                 `yaml:"tags,omitempty"`
}

func NewOpenAPIObject() *OpenAPIObject {
	return &OpenAPIObject{
		Version:    "3.0.0",
		Paths:      *NewPathsObject(),
		Components: *NewComponentsObject(),
	}
}
