package examples

/*
	@API Test API
	@Version 1.0.0
	@Server https://example.com/
		Production API Server
	@Server https://{env}.example.com/
		Internal API Server
		@Variable env staging dev qa staging
			Environment ID

	@SecuritySchemeAPIKey admin_key header X-Admin-Key
		Key for administrative operations
	@SecuritySchemeHTTP access_token Bearer JWT
		Access Token
	@SecurityRequirement access_token

	@Tag User
		User APIs
*/
func Main() {
}

/*
	@Parameter id path
		ID of user in UUID format.
		@JSONSchema
			{ "type": "string" }
*/
type UserID string

// @JSONSchema
const UserSchema = `
{
	"$id": "#User",
	"type": "object",
	"properties": {
		"id": { "type": "string" },
		"name": { "type": "string" }
	}
}
`

type User struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

/*
	@Response UserResponse
		Describe user information.
		@JSONSchema
			{
				"type": "object",
				"properties": {
					"user": { "$ref": "#User" }
				}
			}
		@JSONExample User - Example User
			{
				"user": {
					"id": "C9D3C933-E1BC-46CC-B8F9-057951DD63B0",
					"name": "Test User"
				}
			}
*/
type UserResponse struct {
	User User `json:"user"`
}

/*
	@Operation GET /user/{id} - Get User
		Return user with specific ID.

		If `id` parameter is "me", returns the calling user.

		@Tag User

		@Parameter {UserID}
		@Response 200 {UserResponse}
*/
func GetUser(id UserID) UserResponse {
	return UserResponse{}
}

/*
	@Operation POST /user - Create User
		Create new user with specified information.

		@Tag User
		@SecurityRequirement admin_key

		@RequestBody
			Describe information of the user.
			@JSONSchema
				{
					"type": "object",
					"properties": {
						"name": { "type": "string" }
					}
				}
			@JSONExample TestUser - User with name 'Test'
				{
					"name": "Test"
				}
		@Response 200 {UserResponse}
		@Response 403
			Client is not authorized to create user.
*/
func CreateUser(user User) UserResponse {
	return UserResponse{}
}
