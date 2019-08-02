package openapi3

type ComponentsObject struct {
	Schemas         map[string]*Schema               `yaml:"schemas,omitempty"`
	Parameters      map[string]*ParameterObject      `yaml:"parameters,omitempty"`
	RequestBodies   map[string]*RequestBodyObject    `yaml:"requestBodies,omitempty"`
	Responses       map[string]*ResponseObject       `yaml:"responses,omitempty"`
	SecuritySchemes map[string]*SecuritySchemeObject `yaml:"securitySchemes,omitempty"`
	Callbacks       map[string]*CallbackObject       `yaml:"callbacks,omitempty"`
}

func NewComponentsObject() *ComponentsObject {
	return &ComponentsObject{
		Schemas:         map[string]*Schema{},
		Parameters:      map[string]*ParameterObject{},
		RequestBodies:   map[string]*RequestBodyObject{},
		Responses:       map[string]*ResponseObject{},
		SecuritySchemes: map[string]*SecuritySchemeObject{},
		Callbacks:       map[string]*CallbackObject{},
	}
}
