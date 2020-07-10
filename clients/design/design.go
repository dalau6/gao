*
This is the design file. It contains the API specification, methods, inputs, and outputs using Goa DSL code. The objective is to use this as a single source of truth for the entire API source code.
*/
package design
 
import (
    . "goa.design/goa/v3/dsl"
)
 
// Main API declaration
var _ = API("clients", func() {
    Title("An api for clients")
    Description("This api manages clients with CRUD operations")
    Server("clients", func() {
        Host("localhost", func() {
            URI("http://localhost:8080/api/v1")
        })
    })
})
 
// Client Service declaration with two methods and Swagger API specification file
var _ = Service("client", func() {
    Description("The Client service allows access to client members")
    Method("add", func() {
        Payload(func() {
            Field(1, "ClientID", String, "Client ID")
            Field(2, "ClientName", String, "Client ID")
            Required("ClientID", "ClientName")
        })
        Result(Empty)
        Error("not_found", NotFound, "Client not found")
        HTTP(func() {
            POST("/api/v1/client/{ClientID}")
            Response(StatusCreated)
        })
    })
 
    Method("get", func() {
		Payload(func() {
			Field(1, "ClientID", String, "Client ID")
			Required("ClientID")
		})
		Result(ClientManagement)
		Error("not_found", NotFound, "Client not found")
		HTTP(func() {
			GET("/api/v1/client/{ClientID}")
			Response(StatusOK)
		})
	})
        	
	Method("show", func() {
		Result(CollectionOf(ClientManagement))
		HTTP(func() {
			GET("/api/v1/client")
			Response(StatusOK)
		})
	})
	Files("/openapi.json", "./gen/http/openapi.json")
})
 
// ClientManagement is a custom ResultType used to configure views for our custom type
var ClientManagement = ResultType("application/vnd.client", func() {
	Description("A ClientManagement type describes a Client of company.")
	Reference(Client)
	TypeName("ClientManagement")

	Attributes(func() {
		Attribute("ClientID", String, "ID is the unique id of the Client.", func() {
			Example("ABCDEF12356890")
		})
		Field(2, "ClientName")
	})

	View("default", func() {
		Attribute("ClientID")
		Attribute("ClientName")
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
	Required("ClientID", "ClientName")
})
 
// NotFound is a custom type where we add the queried field in the response
var NotFound = Type("NotFound", func() {
	Description("NotFound is the type returned when the requested data that does not exist.")
	Attribute("message", String, "Message of error", func() {
		Example("Client ABCDEF12356890 not found")
	})
	Field(2, "id", String, "ID of missing data")
	Required("message", "id")
})