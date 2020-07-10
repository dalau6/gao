/*
This is the design file. It contains the API specification, methods, inputs, and outputs using Goa DSL code. The objective is to use this as a single source of truth for the entire API source code.
*/
package design

import (
	. "goa.design/goa/v3/dsl"
	cors "goa.design/plugins/v3/cors/dsl"
)

// Main API declaration
var _ = API("clients", func() {
	Title("An api for clients")
	Description("This api manages clients with CRUD operations")
	cors.Origin("/.*localhost.*/", func() {
		cors.Headers("X-Authorization", "X-Time", "X-Api-Version",
			"Content-Type", "Origin", "Authorization")
		cors.Methods("GET", "POST", "OPTIONS")
		cors.Expose("Content-Type", "Origin")
		cors.MaxAge(100)
		cors.Credentials()
	})
	Server("clients", func() {
		Host("localhost", func() {
			URI("http://localhost:8080/api/v1")
		})
	})
})

// Client Service declaration with two methods and Swagger API specification file
var _ = Service("client", func() {
	Description("The Client service allows access to client members")
	Error("unauthorized", String, "Credentials are invalid")
	HTTP(func() {
		Response("unauthorized", StatusUnauthorized)
	})
	Method("add", func() {
		Payload(func() {
			TokenField(1, "token", String, func() {
				Description("JWT used for authentication")
			})
			Field(2, "ClientID", String, "Client ID")
			Field(3, "ClientName", String, "Client ID")
			Field(4, "ContactName", String, "Contact Name")
			Field(5, "ContactEmail", String, "Contact Email")
			Field(6, "ContactMobile", Int, "Contact Mobile Number")
			Required("token",
				"ClientID", "ClientName", "ContactName",
				"ContactEmail", "ContactMobile")
		})
		Security(JWTAuth, func() {
			Scope("api:write")
		})
		Result(Empty)
		Error("invalid-scopes", String, "Token scopes are invalid")
		Error("not_found", NotFound, "Client not found")
		HTTP(func() {
			POST("/api/v1/client/{ClientID}")
			Header("token:X-Authorization")
			Response("invalid-scopes", StatusForbidden)
			Response(StatusCreated)
		})
	})

	Method("get", func() {
		Payload(func() {
			TokenField(1, "token", String, func() {
				Description("JWT used for authentication")
			})
			Field(2, "ClientID", String, "Client ID")
			Required("token", "ClientID")
		})
		Security(JWTAuth, func() {
			Scope("api:read")
		})
		Result(ClientManagement)
		Error("invalid-scopes", String, "Token scopes are invalid")
		Error("not_found", NotFound, "Client not found")
		HTTP(func() {
			GET("/api/v1/client/{ClientID}")
			Header("token:X-Authorization")
			Response("invalid-scopes", StatusForbidden)
			Response(StatusOK)
		})
	})

	Method("show", func() {
		Payload(func() {
			TokenField(1, "token", String, func() {
				Description("JWT used for authentication")
			})
			Required("token")
		})
		Security(JWTAuth, func() {
			Scope("api:read")
		})
		Result(CollectionOf(ClientManagement))
		Error("invalid-scopes", String, "Token scopes are invalid")
		HTTP(func() {
			GET("/api/v1/client")
			Header("token:X-Authorization")
			Response("invalid-scopes", StatusForbidden)
			Response(StatusOK)
		})
	})
	Files("/openapi.json", "./gen/http/openapi.json")
})

// ClientManagement is a custom ResultType used to
// configure views for our custom type
var ClientManagement = ResultType("application/vnd.client", func() {
	Description("A ClientManagement type describes a Client of company.")
	Reference(Client)
	TypeName("ClientManagement")

	Attributes(func() {
		Attribute("ClientID", String, "ID is the unique id of the Client.", func() {
			Example("ABCDEF12356890")
		})
		Field(2, "ClientName")
		Attribute("ContactName", String, "Name of the Contact.", func() {
			Example("John Doe")
		})
		Field(4, "ContactEmail")
		Field(5, "ContactMobile")
	})

	View("default", func() {
		Attribute("ClientID")
		Attribute("ClientName")
		Attribute("ContactName")
		Attribute("ContactEmail")
		Attribute("ContactMobile")
	})

	Required("ClientID")
})

// Client is the custom type for clients in our database
var Client = Type("Client", func() {
	Description("Client describes a customer of company.")
	Attribute("ClientID", String, "ID is the unique id of the Client Member.", func() {
		Example("ABCDEF12356890")
	})
	Attribute("ClientName", String, "Name of the Client", func() {
		Example("John Doe Limited")
	})
	Attribute("ContactName", String, "Name of the Client Contact.", func() {
		Example("John Doe")
	})
	Attribute("ContactEmail", String, "Email of the Client Contact", func() {
		Example("john.doe@johndoe.com")
	})
	Attribute("ContactMobile", Int, "Mobile number of the Client Contact", func() {
		Example(12365474235)
	})
	Required("ClientID", "ClientName", "ContactName", "ContactEmail", "ContactMobile")
})

// NotFound is a custom type where we add the queried field in the response
var NotFound = Type("NotFound", func() {
	Description("NotFound is the type returned " +
		"when the requested data that does not exist.")
	Attribute("message", String, "Message of error", func() {
		Example("Client ABCDEF12356890 not found")
	})
	Field(2, "id", String, "ID of missing data")
	Required("message", "id")
})

// Creds is a custom type for replying Tokens
var Creds = Type("Creds", func() {
	Field(1, "jwt", String, "JWT token", func() {
		Example("eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9." +
			"eyJzdWIiOiIxMjM0NTY3ODkwIiwibmFtZSI6IkpvaG4gRG9" +
			"lIiwiYWRtaW4iOnRydWV9.TJVA95OrM7E2cBab30RMHrHD" +
			"cEfxjoYZgeFONFh7HgQ")
	})
	Required("jwt")
})

// JWTAuth is the JWTSecurity DSL function for adding JWT support in the API
var JWTAuth = JWTSecurity("jwt", func() {
	Description(`Secures endpoint by requiring a valid
JWT token retrieved via the signin endpoint. Supports
scopes "api:read" and "api:write".`)
	Scope("api:read", "Read-only access")
	Scope("api:write", "Read and write access")
})

// BasicAuth is the BasicAuth DSL function for
// adding basic auth support in the API
var BasicAuth = BasicAuthSecurity("basic", func() {
	Description("Basic authentication used to " +
		"authenticate security principal during signin")
	Scope("api:read", "Read-only access")
})

// Signin Service is the service used to authenticate users and assign JWT tokens for their sessions
var _ = Service("signin", func() {
	Description("The Signin service authenticates users and validate tokens")
	Error("unauthorized", String, "Credentials are invalid")
	HTTP(func() {
		Response("unauthorized", StatusUnauthorized)
	})
	Method("authenticate", func() {
		Description("Creates a valid JWT")
		Security(BasicAuth)
		Payload(func() {
			Description("Credentials used to authenticate to retrieve JWT token")
			UsernameField(1, "username",
				String, "Username used to perform signin", func() {
					Example("user")
				})
			PasswordField(2, "password",
				String, "Password used to perform signin", func() {
					Example("password")
				})
			Required("username", "password")
		})
		Result(Creds)
		HTTP(func() {
			POST("/signin/authenticate")
			Response(StatusOK)
		})
	})
})
