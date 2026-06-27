package main

// @asyncapi
// @asyncTitle AsyncAPI Test API
// @asyncVersion 1.0.0
// @asyncDescription A test API for AsyncAPI generation
// @asyncServer.name development
// @asyncServer.host localhost:9092
// @asyncServer.protocol kafka
// @asyncDefaultContentType application/json

// @asyncapi
// @asyncChannel userEvents
// @asyncChannel.address user.events
// @asyncChannel.description Events related to users
// @asyncOperation sendUserEvent
// @asyncAction send
// @asyncSummary Send a user event
// @asyncMessage {object} UserEvent "A user event message"
// @asyncTag user
func publishUserEvent() {}

// UserEvent represents a user-related event.
type UserEvent struct {
	ID        string `json:"id"`
	Name      string `json:"name"`
	Email     string `json:"email"`
	CreatedAt string `json:"created_at"`
}
