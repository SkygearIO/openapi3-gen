package openapi3

type ReferenceObject interface{}

func MakeSchemaRef(id string) ReferenceObject {
	return map[string]string{
		"$ref": "#/components/schemas/" + id,
	}
}

func MakeParameterRef(id string) ReferenceObject {
	return map[string]string{
		"$ref": "#/components/parameters/" + id,
	}
}

func MakeRequestBodyRef(id string) ReferenceObject {
	return map[string]string{
		"$ref": "#/components/requestBodies/" + id,
	}
}

func MakeResponseRef(id string) ReferenceObject {
	return map[string]string{
		"$ref": "#/components/responses/" + id,
	}
}

func MakeCallbackRef(id string) ReferenceObject {
	return map[string]string{
		"$ref": "#/components/callbacks/" + id,
	}
}
