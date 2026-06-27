package swag

import (
	"testing"

	"github.com/stretchr/testify/assert"
	openapi "github.com/sv-tools/openapi/spec"
)

func TestAsyncOperation_ParseComment(t *testing.T) {
	t.Parallel()

	parser := New()
	ao := NewAsyncOperation(parser)

	t.Run("parse @asyncapi marker", func(t *testing.T) {
		err := ao.ParseComment("// @asyncapi", nil)
		assert.NoError(t, err)
	})

	t.Run("parse @asyncAction", func(t *testing.T) {
		err := ao.ParseComment("// @asyncAction send", nil)
		assert.NoError(t, err)
		assert.Equal(t, "send", ao.Operation.Action)
	})

	t.Run("parse @asyncChannel.name", func(t *testing.T) {
		err := ao.ParseComment("// @asyncChannel.name userEvents", nil)
		assert.NoError(t, err)
		assert.Equal(t, "userEvents", ao.ChannelName)
	})

	t.Run("parse @asyncChannel.address", func(t *testing.T) {
		err := ao.ParseComment("// @asyncChannel.address user.events", nil)
		assert.NoError(t, err)
		assert.Equal(t, "user.events", ao.Channel.Address)
	})

	t.Run("parse @asyncChannel.description", func(t *testing.T) {
		err := ao.ParseComment("// @asyncChannel.description Events about users", nil)
		assert.NoError(t, err)
		assert.Equal(t, "Events about users", ao.Channel.Description)
	})

	t.Run("parse @asyncOperation.name", func(t *testing.T) {
		err := ao.ParseComment("// @asyncOperation.name sendUserEvent", nil)
		assert.NoError(t, err)
		assert.Equal(t, "sendUserEvent", ao.OperationName)
	})

	t.Run("parse @asyncSummary", func(t *testing.T) {
		err := ao.ParseComment("// @asyncSummary Send a user event", nil)
		assert.NoError(t, err)
		assert.Equal(t, "Send a user event", ao.Operation.Summary)
	})

	t.Run("parse @asyncDescription", func(t *testing.T) {
		err := ao.ParseComment("// @asyncDescription First line", nil)
		assert.NoError(t, err)
		err = ao.ParseComment("// @asyncDescription Second line", nil)
		assert.NoError(t, err)
		assert.Equal(t, "First line\nSecond line", ao.Operation.Description)
	})

	t.Run("parse @asyncTag", func(t *testing.T) {
		err := ao.ParseComment("// @asyncTag user", nil)
		assert.NoError(t, err)
		assert.Len(t, ao.Tags, 1)
		assert.Equal(t, "user", ao.Tags[0].Name)
	})

	t.Run("parse @asyncServer.name", func(t *testing.T) {
		err := ao.ParseComment("// @asyncServer.name development", nil)
		assert.NoError(t, err)
		assert.Equal(t, "development", ao.ServerName)
		assert.NotNil(t, ao.Server)
	})

	t.Run("parse @asyncServer.host", func(t *testing.T) {
		err := ao.ParseComment("// @asyncServer.host localhost:9092", nil)
		assert.NoError(t, err)
		assert.Equal(t, "localhost:9092", ao.Server.Host)
	})

	t.Run("parse @asyncServer.protocol", func(t *testing.T) {
		err := ao.ParseComment("// @asyncServer.protocol kafka", nil)
		assert.NoError(t, err)
		assert.Equal(t, "kafka", ao.Server.Protocol)
	})

	t.Run("parse @asyncExternalDocs", func(t *testing.T) {
		err := ao.ParseComment("// @asyncExternalDocs.description Read more", nil)
		assert.NoError(t, err)
		err = ao.ParseComment("// @asyncExternalDocs.url https://example.com", nil)
		assert.NoError(t, err)
		assert.Equal(t, "Read more", ao.Operation.ExternalDocs.Description)
		assert.Equal(t, "https://example.com", ao.Operation.ExternalDocs.URL)
	})
}

func TestAsyncOperation_ParseMessageComment(t *testing.T) {
	t.Parallel()

	parser := New(SetParseDependency(ParseModels))
	parser.packages.uniqueDefinitions["UserEvent"] = &TypeSpecDef{}
	parser.parsedSchemasV3[parser.packages.uniqueDefinitions["UserEvent"]] = &SchemaV3{
		PkgPath: "",
		Name:    "UserEvent",
		Schema:  &openapi.Schema{},
	}

	ao := NewAsyncOperation(parser)

	t.Run("parse simple message", func(t *testing.T) {
		err := ao.ParseComment("// @asyncMessage {object} UserEvent \"A user event\"", nil)
		assert.NoError(t, err)
		assert.Len(t, ao.Messages, 1)
		assert.Equal(t, "UserEvent", ao.Messages[0].Name)
		assert.Equal(t, "A user event", ao.Messages[0].Description)
	})
}

func TestGenerateAsyncOperationID(t *testing.T) {
	t.Parallel()

	assert.Equal(t, "sendOperation", generateAsyncOperationID("", "send"))
	assert.Equal(t, "sendUserEvents", generateAsyncOperationID("userEvents", "send"))
}
