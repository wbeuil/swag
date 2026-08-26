package main

//	@asyncapi
//	@asyncTitle					AsyncAPI Complete Example
//	@asyncVersion				1.0.0
//	@asyncDescription			Two operations, bindings, and extensions in one file
//	@asyncServer.name			development
//	@asyncServer.host			localhost:4566
//	@asyncServer.protocol		sqs
//	@asyncServer.binding.sqs	{"bindingVersion":"0.3.0"}
//	@asyncDefaultContentType	application/json

// @asyncapi
// @asyncChannel				email
// @asyncChannel.address		email
// @asyncChannel.description	Outbound email jobs
// @asyncOperation				receiveEmail
// @asyncAction				receive
// @asyncSummary				Receive an email job
// @asyncMessage				{object}	EmailJob	"Email job payload"
// @asyncMessage.contentType	application/json
// @asyncBinding.channel.sqs	{"queue":{"name":"email","fifoQueue":false},"deadLetterQueue":{"name":"email-dlq"},"bindingVersion":"0.3.0"}
// @asyncx-owner				{"team":"platform"}
func consumeEmail() {}

// @asyncapi
// @asyncChannel		notifications
// @asyncChannel.address	notifications
// @asyncOperation		sendNotification
// @asyncAction		send
// @asyncSummary		Send a notification
// @asyncMessage		{object}	string	"Raw notification payload"
// @asyncMessage.name	NotificationPayload
func publishNotification() {}

// GetPet is an HTTP-only handler whose model must not leak into AsyncAPI.
//
//	@Summary	Get a pet
//	@Success	200	{object}	Pet
//	@Router		/pet [get]
func getPet() {}

// Pet is an OpenAPI-only model.
type Pet struct {
	Name string `json:"name"`
}

// EmailJob is the SQS email payload.
type EmailJob struct {
	ToAddress    string            `json:"toAddress"`
	TemplateType string            `json:"templateType"`
	Data         map[string]string `json:"data"`
	Recipient    Recipient         `json:"recipient"`
}

// Recipient is nested from EmailJob and should be copied into AsyncAPI.
type Recipient struct {
	Email string `json:"email"`
}
