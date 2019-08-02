package openapi3

type OperationObject struct {
	Tags        []string                    `yaml:"tags,omitempty"`
	Summary     string                      `yaml:"summary,omitempty"`
	Description string                      `yaml:"description,omitempty"`
	ID          string                      `yaml:"operationId,omitempty"`
	Parameters  []Parameter                 `yaml:"parameters,omitempty"`
	RequestBody RequestBody                 `yaml:"requestBody,omitempty"`
	Responses   map[string]Response         `yaml:"responses"`
	Callbacks   map[string]Callback         `yaml:"callbacks,omitempty"`
	Security    []SecurityRequirementObject `yaml:"security,omitempty"`
}

func NewOperationObject() *OperationObject {
	return &OperationObject{
		Responses: map[string]Response{},
		Callbacks: map[string]Callback{},
	}
}
