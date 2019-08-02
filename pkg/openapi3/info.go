package openapi3

type InfoObject struct {
	Title       string `yaml:"title"`
	Description string `yaml:"description,omitempty"`
	Version     string `yaml:"version,omitempty"`
}
