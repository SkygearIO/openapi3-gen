package openapi3

type ReferenceObject map[string]interface{}

func MakeSchemaRef(id string) ReferenceObject {
	return map[string]interface{}{
		"$ref": "#/components/schemas/" + id,
	}
}

func MakeParameterRef(id string) ReferenceObject {
	return map[string]interface{}{
		"$ref": "#/components/parameters/" + id,
	}
}

func MakeRequestBodyRef(id string) ReferenceObject {
	return map[string]interface{}{
		"$ref": "#/components/requestBodies/" + id,
	}
}

func MakeResponseRef(id string) ReferenceObject {
	return map[string]interface{}{
		"$ref": "#/components/responses/" + id,
	}
}

func MakeCallbackRef(id string) ReferenceObject {
	return map[string]interface{}{
		"$ref": "#/components/callbacks/" + id,
	}
}
