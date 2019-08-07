package processor

import (
	"go/parser"
	"go/token"
	"testing"

	"github.com/skygeario/openapi3-gen/pkg/openapi3"
	. "github.com/smartystreets/goconvey/convey"
)

func TestProcessor(t *testing.T) {
	process := func(sources ...string) (*openapi3.OpenAPIObject, []error) {
		psr := New()
		fset := token.NewFileSet()
		for _, src := range sources {
			file, _ := parser.ParseFile(fset, "", src, parser.ParseComments)
			psr.Process(fset, file)
		}
		return psr.oapi, psr.errs
	}

	Convey("Processor", t, func() {
		Convey("should process top-level annotations", func() {
			oapi, errs := process(`
				package main

				/*
					@API Test API
					@Version 1.0.0
					@Server https://example.com/
						Default API Server
					@Server https://{app}.example.com/
						API Server
						@Variable app my_app app1 app2 app3
							App ID

					@SecurityAPIKey api_key header X-API-Key
						API Key
					@SecurityHTTP access_token Bearer JWT
						Access Token
					@Security api_key
					@Security access_token
				
					@Tag User
						User APIs
					@Tag Admin
						Admin APIs
				*/
				func main() {}
			`)

			So(errs, ShouldBeEmpty)
			So(oapi.Info.Title, ShouldEqual, "Test API")
			So(oapi.Info.Version, ShouldEqual, "1.0.0")
			So(oapi.Servers, ShouldResemble, []openapi3.ServerObject{
				openapi3.ServerObject{
					URL:         "https://example.com/",
					Description: "Default API Server",
					Variables:   map[string]openapi3.ServerVariable{},
				},
				openapi3.ServerObject{
					URL:         "https://{app}.example.com/",
					Description: "API Server",
					Variables: map[string]openapi3.ServerVariable{
						"app": openapi3.ServerVariable{
							Description: "App ID",
							Default:     "my_app",
							Enum:        []string{"app1", "app2", "app3"},
						},
					},
				},
			})
			So(oapi.Components.SecuritySchemes, ShouldResemble, map[string]*openapi3.SecuritySchemeObject{
				"api_key": &openapi3.SecuritySchemeObject{
					Type:           openapi3.SecuritySchemeTypeAPIKey,
					Description:    "API Key",
					APIKeyName:     "X-API-Key",
					APIKeyLocation: openapi3.SecuritySchemeAPIKeyLocationHeader,
				},
				"access_token": &openapi3.SecuritySchemeObject{
					Type:             openapi3.SecuritySchemeTypeHTTP,
					Description:      "Access Token",
					HTTPAuthScheme:   "bearer",
					HTTPBearerFormat: "JWT",
				},
			})
			So(oapi.Security, ShouldResemble, []openapi3.SecurityRequirementObject{
				openapi3.SecurityRequirementObject{"api_key": []string{}},
				openapi3.SecurityRequirementObject{"access_token": []string{}},
			})
			So(oapi.Tags, ShouldResemble, []openapi3.TagObject{
				openapi3.TagObject{Name: "User", Description: "User APIs"},
				openapi3.TagObject{Name: "Admin", Description: "Admin APIs"},
			})
		})

		Convey("should process component annotations", func() {
			oapi, errs := process(`
				package main

				// @JSONSchema
				const TestSchema = ` + "`" + `
				{
					"$id": "#TestSchema",
					"type": "object"
				}
				` + "`" + `

				/*
					@RequestBody
						Test Request.
				*/
				type TestRequest struct {}

				/*
					@Response
						Test Response.
				*/
				type TestResponse struct {}

				/*
					@Callback
						Test Callback.
						@Operation POST /test-1 - Test callback 1
						@Operation PUT /test-1 - Test callback 2
				*/
				type TestCallback struct {}
				
				/*
					@Parameter test query
						Test Parameter.
				*/
				type TestParameter string
			`)

			schema := openapi3.Schema(map[string]interface{}{
				"type": "object",
			})
			requestBody := openapi3.NewRequestBodyObject()
			requestBody.Description = "Test Request."
			response := openapi3.NewResponseObject()
			response.Description = "Test Response."
			callbackPost := openapi3.NewOperationObject()
			callbackPost.Summary = "Test callback 1"
			callbackPut := openapi3.NewOperationObject()
			callbackPut.Summary = "Test callback 2"
			callback := openapi3.NewCallbackObject()
			(*callback)["/test-1"] = openapi3.PathItemObject{
				Post: callbackPost,
				Put:  callbackPut,
			}
			param := openapi3.NewParameterObject()
			param.Name = "test"
			param.Location = openapi3.ParameterLocationQuery
			param.Description = "Test Parameter."

			So(errs, ShouldBeEmpty)
			So(oapi.Components.Schemas, ShouldResemble, map[string]*openapi3.Schema{
				"TestSchema": &schema,
			})
			So(oapi.Components.RequestBodies, ShouldResemble, map[string]*openapi3.RequestBodyObject{
				"TestRequest": requestBody,
			})
			So(oapi.Components.Responses, ShouldResemble, map[string]*openapi3.ResponseObject{
				"TestResponse": response,
			})
			So(oapi.Components.Callbacks, ShouldResemble, map[string]*openapi3.CallbackObject{
				"TestCallback": callback,
			})
			So(oapi.Components.Parameters, ShouldResemble, map[string]*openapi3.ParameterObject{
				"TestParameter": param,
			})
		})

		Convey("should use specified component ID", func() {
			oapi, errs := process(`
				package main

				/*
					@ID EventContext
					@RequestBody
						Test Context.
				*/
				type Context struct {}
			`)

			requestBody := openapi3.NewRequestBodyObject()
			requestBody.Description = "Test Context."

			So(errs, ShouldBeEmpty)
			So(oapi.Components.RequestBodies, ShouldResemble, map[string]*openapi3.RequestBodyObject{
				"EventContext": requestBody,
			})
		})

		Convey("should translate JSON Schemas", func() {
			oapi, errs := process(`
				package main

				// @JSONSchema
				const TestSchema1 = ` + "`" + `
				{
					"$id": "#TestSchema1",
					"type": "string"
				}
				` + "`" + `

				// @JSONSchema
				const TestSchema2 = ` + "`" + `
				{
					"$id": "#TestSchema2",
					"type": "array",
					"items": [
						{ "type": "string" },
						{ "$ref": "#TestSchema1" }
					]
				}
				` + "`" + `
			`)

			schema1 := openapi3.Schema(map[string]interface{}{
				"type": "string",
			})
			schema2 := openapi3.Schema(map[string]interface{}{
				"type": "array",
				"items": []interface{}{
					map[string]interface{}{"type": "string"},
					map[string]interface{}{"$ref": "#/components/schemas/TestSchema1"},
				},
			})

			So(errs, ShouldBeEmpty)
			So(oapi.Components.Schemas, ShouldResemble, map[string]*openapi3.Schema{
				"TestSchema1": &schema1,
				"TestSchema2": &schema2,
			})
		})

		Convey("should process operation annotations", func() {
			Convey("using reference objects", func() {
				oapi, errs := process(`
					package main
	
					/*
						@Operation GET /user/{id} - Get User
							Return user with specific ID.
					
							If ` + "`" + `id` + "`" + ` parameter is "me", returns the calling user.
	
							@Tag User Object
						
							@Parameter {UserID}
							@Response 200 {UserResponse}
					
						@Operation POST /user - Create User
							Create new user with specified information.
	
							@Tag User Object
							@Security admin_key
						
							@RequestBody {CreateUserRequest}
							@Response default {ErrorResponse}
							@Response 200 {UserResponse}
						
							@Callback user_create {UserCreateEvent}
					*/
					type UserController struct {}
				`)

				getOp := openapi3.NewOperationObject()
				getOp.Summary = "Get User"
				getOp.Description = "Return user with specific ID.\n\nIf `id` parameter is \"me\", returns the calling user."
				getOp.Tags = []string{"User Object"}
				getOp.Parameters = []openapi3.Parameter{
					openapi3.ReferenceObject{"$ref": "#/components/parameters/UserID"},
				}
				getOp.Responses = map[string]openapi3.Response{
					"200": openapi3.ReferenceObject{"$ref": "#/components/responses/UserResponse"},
				}

				postOp := openapi3.NewOperationObject()
				postOp.Summary = "Create User"
				postOp.Description = "Create new user with specified information."
				postOp.Tags = []string{"User Object"}
				postOp.Security = []openapi3.SecurityRequirementObject{
					openapi3.SecurityRequirementObject{"admin_key": []string{}},
				}
				postOp.RequestBody = openapi3.ReferenceObject{
					"$ref": "#/components/requestBodies/CreateUserRequest",
				}
				postOp.Responses = map[string]openapi3.Response{
					"default": openapi3.ReferenceObject{"$ref": "#/components/responses/ErrorResponse"},
					"200":     openapi3.ReferenceObject{"$ref": "#/components/responses/UserResponse"},
				}
				postOp.Callbacks = map[string]openapi3.Callback{
					"user_create": openapi3.ReferenceObject{"$ref": "#/components/callbacks/UserCreateEvent"},
				}

				So(errs, ShouldBeEmpty)
				So(oapi.Paths, ShouldResemble, openapi3.PathsObject{
					"/user/{id}": openapi3.PathItemObject{Get: getOp},
					"/user":      openapi3.PathItemObject{Post: postOp},
				})
			})
			Convey("using inline objects", func() {
				oapi, errs := process(`
					package main
	
					/*
						@Operation PATCH /user/{id} - Update User
							Update new user with specified information.
	
							@Tag User Object
						
							@Parameter id path
								ID of user in UUID format.
								@JSONSchema
									{ "type": "string" }

							@RequestBody
								Describe the new information of user.
								@JSONSchema
									{
										"type": "object",
										"properties": {
											"name": { "type": "string" }
										}
									}
								@JSONExample UpdateName - Update name of user
									{
										"name": "Test"
									}

							@Response 200
								Updated the user successfully.

							@Callback user_updated
								@Operation POST /user_updated - User is updated
									A user is updated.

									@Response 200
										Acknowledge the event.
					*/
					func DeleteUser(userID string) error {
						return nil
					}
				`)

				patchOp := openapi3.NewOperationObject()
				patchOp.Summary = "Update User"
				patchOp.Description = "Update new user with specified information."
				patchOp.Tags = []string{"User Object"}

				param := openapi3.NewParameterObject()
				param.Name = "id"
				param.Location = openapi3.ParameterLocationPath
				param.Description = "ID of user in UUID format."
				param.Required = true
				param.Schema = openapi3.Schema(map[string]interface{}{"type": "string"})
				patchOp.Parameters = []openapi3.Parameter{param}

				req := openapi3.NewRequestBodyObject()
				reqMediaType := openapi3.NewMediaTypeObject()
				reqMediaType.Schema = openapi3.Schema(map[string]interface{}{
					"type": "object",
					"properties": map[string]interface{}{
						"name": map[string]interface{}{"type": "string"},
					},
				})
				reqMediaType.Examples = map[string]openapi3.ExampleObject{
					"UpdateName": openapi3.ExampleObject{
						Summary: "Update name of user",
						Value: map[string]interface{}{
							"name": "Test",
						},
					},
				}
				req.Description = "Describe the new information of user."
				req.Content = map[string]openapi3.MediaTypeObject{
					"application/json": *reqMediaType,
				}
				patchOp.RequestBody = req

				resp := openapi3.NewResponseObject()
				resp.Description = "Updated the user successfully."
				patchOp.Responses = map[string]openapi3.Response{
					"200": resp,
				}

				callback := openapi3.NewCallbackObject()
				callbackOp := openapi3.NewOperationObject()
				callbackOp.Summary = "User is updated"
				callbackOp.Description = "A user is updated."
				callbackOpResp := openapi3.NewResponseObject()
				callbackOpResp.Description = "Acknowledge the event."
				callbackOp.Responses = map[string]openapi3.Response{
					"200": callbackOpResp,
				}
				(*callback)["/user_updated"] = openapi3.PathItemObject{Post: callbackOp}
				patchOp.Callbacks = map[string]openapi3.Callback{
					"user_updated": callback,
				}

				So(errs, ShouldBeEmpty)
				So(oapi.Paths, ShouldResemble, openapi3.PathsObject{
					"/user/{id}": openapi3.PathItemObject{Patch: patchOp},
				})
			})
		})
	})
}
