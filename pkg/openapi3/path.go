package openapi3

import "net/http"

type Paths interface {
	SetPath(path string, item PathItemObject)
	GetPath(path string) PathItemObject
}

type PathsObject map[string]PathItemObject

func NewPathsObject() *PathsObject {
	return &PathsObject{}
}

func (paths *PathsObject) SetPath(path string, item PathItemObject) {
	(*paths)[path] = item
}

func (paths *PathsObject) GetPath(path string) PathItemObject {
	return (*paths)[path]
}

type PathItemObject struct {
	Summary     string            `yaml:"summary,omitempty"`
	Description string            `yaml:"description,omitempty"`
	Parameters  []ParameterObject `yaml:"parameters,omitempty"`
	Get         *OperationObject  `yaml:"get,omitempty"`
	Put         *OperationObject  `yaml:"put,omitempty"`
	Post        *OperationObject  `yaml:"post,omitempty"`
	Delete      *OperationObject  `yaml:"delete,omitempty"`
	Options     *OperationObject  `yaml:"options,omitempty"`
	Head        *OperationObject  `yaml:"head,omitempty"`
	Patch       *OperationObject  `yaml:"patch,omitempty"`
	Trace       *OperationObject  `yaml:"trace,omitempty"`
}

func (path *PathItemObject) SetOperation(method string, op *OperationObject) bool {
	switch method {
	case http.MethodGet:
		path.Get = op
	case http.MethodPut:
		path.Put = op
	case http.MethodPost:
		path.Post = op
	case http.MethodDelete:
		path.Delete = op
	case http.MethodOptions:
		path.Options = op
	case http.MethodHead:
		path.Head = op
	case http.MethodPatch:
		path.Patch = op
	case http.MethodTrace:
		path.Trace = op
	default:
		return false
	}
	return true
}
