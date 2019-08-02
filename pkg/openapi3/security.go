package openapi3

type SecuritySchemeType string

const (
	SecuritySchemeTypeAPIKey = "apiKey"
	SecuritySchemeTypeHTTP   = "http"
)

func (t SecuritySchemeType) Validate() bool {
	return t == SecuritySchemeTypeAPIKey || t == SecuritySchemeTypeHTTP
}

type SecuritySchemeAPIKeyLocation string

const (
	SecuritySchemeAPIKeyLocationQuery  = "query"
	SecuritySchemeAPIKeyLocationHeader = "header"
	SecuritySchemeAPIKeyLocationCookie = "cookie"
)

func (l SecuritySchemeAPIKeyLocation) Validate() bool {
	return l == SecuritySchemeTypeAPIKey ||
		l == SecuritySchemeAPIKeyLocationHeader ||
		l == SecuritySchemeAPIKeyLocationCookie
}

type SecuritySchemeObject struct {
	Type             SecuritySchemeType           `yaml:"type"`
	Description      string                       `yaml:"description,omitempty"`
	APIKeyName       string                       `yaml:"name,omitempty"`
	APIKeyLocation   SecuritySchemeAPIKeyLocation `yaml:"in,omitempty"`
	HTTPAuthScheme   string                       `yaml:"scheme,omitempty"`
	HTTPBearerFormat string                       `yaml:"bearerFormat,omitempty"`
}

type SecurityRequirementObject map[string][]string
