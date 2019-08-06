package openapi3

type Callback interface{}

type CallbackObject map[string]PathItemObject

func NewCallbackObject() *CallbackObject {
	return &CallbackObject{}
}

func (cb *CallbackObject) SetPath(path string, item PathItemObject) {
	(*cb)[path] = item
}

func (cb *CallbackObject) GetPath(path string) PathItemObject {
	return (*cb)[path]
}
