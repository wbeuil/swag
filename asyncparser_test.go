package swag

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	openapi "github.com/wbeuil/openapi"

	"github.com/swaggo/swag/v2/asyncapi"
)

func TestParser_ParseAsyncAPIInfo(t *testing.T) {
	t.Parallel()

	searchDir := "testdata/asyncapi/simple"
	p := New(SetParseDependency(ParseModels))

	err := p.ParseAPIMultiSearchDir([]string{searchDir}, "main.go", defaultParseDepth)
	require.NoError(t, err)

	asyncAPI := p.GetAsyncAPI()
	require.NotNil(t, asyncAPI)
	assert.Equal(t, "3.1.0", asyncAPI.Asyncapi)

	// Info
	require.NotNil(t, asyncAPI.Info)
	assert.Equal(t, "AsyncAPI Test API", asyncAPI.Info.Title)
	assert.Equal(t, "1.0.0", asyncAPI.Info.Version)
	assert.Equal(t, "A test API for AsyncAPI generation", asyncAPI.Info.Description)
	assert.Equal(t, "urn:com:example:asyncapi", asyncAPI.Id)
	require.NotNil(t, asyncAPI.Info.Contact)
	assert.Equal(t, "Example Support", asyncAPI.Info.Contact.Name)
	assert.Equal(t, "support@example.com", asyncAPI.Info.Contact.Email)
	assert.Equal(t, "https://example.com", asyncAPI.Info.Contact.URL)
	require.NotNil(t, asyncAPI.Info.ExternalDocs)
	assert.Equal(t, "AsyncAPI specification", asyncAPI.Info.ExternalDocs.Description)
	assert.Equal(t, "https://www.asyncapi.com", asyncAPI.Info.ExternalDocs.URL)
	require.Len(t, asyncAPI.Info.Tags, 1)
	assert.Equal(t, "Users", asyncAPI.Info.Tags[0].Name)
	assert.Equal(t, "User-related events", asyncAPI.Info.Tags[0].Description)

	// Servers
	require.Contains(t, asyncAPI.Servers, "development")
	devServer := asyncAPI.Servers["development"]
	assert.Equal(t, "localhost:9092", devServer.Host)
	assert.Equal(t, "kafka", devServer.Protocol)

	// Channels
	require.Contains(t, asyncAPI.Channels, "userEvents")
	channel := asyncAPI.Channels["userEvents"]
	assert.Equal(t, "user.events", channel.Address)
	assert.Equal(t, "User Events", channel.Title)
	assert.Equal(t, "Stream of user lifecycle events", channel.Summary)
	assert.Equal(t, "Events related to users", channel.Description)

	// Operations
	require.Contains(t, asyncAPI.Operations, "sendUserEvent")
	op := asyncAPI.Operations["sendUserEvent"]
	assert.Equal(t, "send", op.Action)
	assert.Equal(t, "Send User Event", op.Title)
	assert.Equal(t, "Send a user event", op.Summary)
	require.NotNil(t, op.Channel)
	assert.Equal(t, "#/channels/userEvents", op.Channel.Ref)

	// Messages
	require.Contains(t, asyncAPI.Components.Messages, "UserEvent")
	msg := asyncAPI.Components.Messages["UserEvent"]
	assert.Equal(t, "UserEvent", msg.Name)
	assert.Equal(t, "User Event", msg.Title)
	assert.Equal(t, "A user event message", msg.Description)
	require.NotNil(t, msg.Payload)
	assert.Equal(t, "#/components/schemas/UserEvent", msg.Payload.Ref)

	// Schemas
	require.Contains(t, asyncAPI.Components.Schemas, "UserEvent")
}

func TestParser_ParseAsyncAPIComplete(t *testing.T) {
	t.Parallel()

	searchDir := "testdata/asyncapi/complete"
	p := New(SetParseDependency(ParseModels), GenerateOpenAPI3Doc(true))

	err := p.ParseAPIMultiSearchDir([]string{searchDir}, "main.go", defaultParseDepth)
	require.NoError(t, err)

	asyncAPI := p.GetAsyncAPI()
	require.NotNil(t, asyncAPI)

	require.Contains(t, asyncAPI.Operations, "receiveEmail")
	require.Contains(t, asyncAPI.Operations, "sendNotification")

	email := asyncAPI.Channels["email"]
	require.NotNil(t, email)
	assert.Equal(t, "email", email.Address)
	require.Contains(t, email.Bindings, "sqs")
	sqs, ok := email.Bindings["sqs"].(map[string]interface{})
	require.True(t, ok)
	assert.Equal(t, "0.3.0", sqs["bindingVersion"])

	dev := asyncAPI.Servers["development"]
	require.NotNil(t, dev)
	require.Contains(t, dev.Bindings, "sqs")

	require.Contains(t, asyncAPI.Components.Messages, "EmailJob")
	require.Contains(t, asyncAPI.Components.Messages, "NotificationPayload")
	assert.Equal(t, "#/components/schemas/EmailJob", asyncAPI.Components.Messages["EmailJob"].Payload.Ref)
	assert.Equal(t, "#/components/schemas/NotificationPayload", asyncAPI.Components.Messages["NotificationPayload"].Payload.Ref)
	require.Contains(t, asyncAPI.Components.Schemas, "EmailJob")
	require.Contains(t, asyncAPI.Components.Schemas, "NotificationPayload")
	require.Contains(t, asyncAPI.Components.Schemas, "main.Recipient")
	require.NotContains(t, asyncAPI.Components.Schemas, "Pet")
	require.NotContains(t, asyncAPI.Components.Schemas, "main.Pet")

	openAPI := p.GetOpenAPI()
	require.Contains(t, openAPI.Components.Spec.Schemas, "main.Pet")
	require.NotContains(t, openAPI.Components.Spec.Schemas, "EmailJob")
	require.NotContains(t, openAPI.Components.Spec.Schemas, "main.Recipient")
	require.Contains(t, asyncAPI.Components.Schemas, "EmailJob")
	require.Contains(t, asyncAPI.Components.Schemas, "main.Recipient")

	assert.Equal(t, map[string]interface{}{"team": "platform"}, asyncAPI.Operations["receiveEmail"].Extensions["x-owner"])
}

func TestCopyOpenAPISchemasToAsyncAPI_OnlyReferenced(t *testing.T) {
	t.Parallel()

	p := New(GenerateOpenAPI3Doc(true))
	p.openAPI.Components.Spec.Schemas = map[string]*openapi.RefOrSpec[openapi.Schema]{
		"Pet": {Spec: &openapi.Schema{Title: "Pet"}},
		"Recipient": {Spec: &openapi.Schema{
			Title: "Recipient",
			Properties: map[string]*openapi.RefOrSpec[openapi.Schema]{
				"address": {Ref: &openapi.Ref{Ref: "#/components/schemas/Address"}},
			},
		}},
		"Address": {Spec: &openapi.Schema{Title: "Address"}},
		"EmailJob": {Spec: &openapi.Schema{
			Title: "EmailJob",
			Properties: map[string]*openapi.RefOrSpec[openapi.Schema]{
				"recipient": {Ref: &openapi.Ref{Ref: "#/components/schemas/Recipient"}},
			},
		}},
	}
	emailJob := *p.openAPI.Components.Spec.Schemas["EmailJob"].Spec
	p.asyncAPI.Components = &asyncapi.Components{
		Schemas: map[string]*openapi.Schema{
			"EmailJob": &emailJob,
		},
	}

	p.copyOpenAPISchemasToAsyncAPI()

	require.Contains(t, p.asyncAPI.Components.Schemas, "EmailJob")
	require.Contains(t, p.asyncAPI.Components.Schemas, "Recipient")
	require.Contains(t, p.asyncAPI.Components.Schemas, "Address")
	require.NotContains(t, p.asyncAPI.Components.Schemas, "Pet")
}

func TestPruneUnreferencedOpenAPISchemas(t *testing.T) {
	t.Parallel()

	p := New(GenerateOpenAPI3Doc(true))
	op := openapi.NewExtendable(&openapi.Operation{})
	op.AddExt("x-ref", "#/components/schemas/Pet")
	p.openAPI.Paths = openapi.NewPaths()
	p.openAPI.Paths.Spec.Add("/pet", openapi.NewRefOrExtSpec[openapi.PathItem](&openapi.PathItem{Get: op}))
	p.openAPI.Components.Spec.Schemas = map[string]*openapi.RefOrSpec[openapi.Schema]{
		"Pet": {Spec: &openapi.Schema{Title: "Pet"}},
		"EmailJob": {Spec: &openapi.Schema{
			Title: "EmailJob",
			Properties: map[string]*openapi.RefOrSpec[openapi.Schema]{
				"recipient": {Ref: &openapi.Ref{Ref: "#/components/schemas/Recipient"}},
			},
		}},
		"Recipient": {Spec: &openapi.Schema{Title: "Recipient"}},
	}

	p.GetOpenAPI()

	require.Contains(t, p.openAPI.Components.Spec.Schemas, "Pet")
	require.NotContains(t, p.openAPI.Components.Spec.Schemas, "EmailJob")
	require.NotContains(t, p.openAPI.Components.Spec.Schemas, "Recipient")
}
