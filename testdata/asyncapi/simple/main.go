package main

//	@asyncapi
//	@asyncTitle					AsyncAPI Test API
//	@asyncVersion				1.0.0
//	@asyncId					urn:com:example:asyncapi
//	@asyncDescription			A test API for AsyncAPI generation
//	@asyncContact.name			Example Support
//	@asyncContact.email			support@example.com
//	@asyncContact.url			https://example.com
//	@asyncExternalDocs.description	AsyncAPI specification
//	@asyncExternalDocs.url		https://www.asyncapi.com
//	@asyncTag					Users	User-related events
//	@asyncServer.name			development
//	@asyncServer.host			localhost:9092
//	@asyncServer.protocol		kafka
//	@asyncDefaultContentType	application/json

//	@asyncapi
//	@asyncChannel				userEvents
//	@asyncChannel.address		user.events
//	@asyncChannel.title			User Events
//	@asyncChannel.summary		Stream of user lifecycle events
//	@asyncChannel.description	Events related to users
//	@asyncOperation				sendUserEvent
//	@asyncOperation.title		Send User Event
//	@asyncAction				send
//	@asyncSummary				Send a user event
//	@asyncMessage				{object}	UserEvent	"A user event message"
//	@asyncMessage.title			User Event
//	@asyncTag					user
func publishUserEvent() {}

// UserEvent represents a user-related event.
type UserEvent struct {
	ID        string `json:"id"`
	Name      string `json:"name"`
	Email     string `json:"email"`
	CreatedAt string `json:"created_at"`
}
