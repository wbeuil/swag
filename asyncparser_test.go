package swag

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
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

	// Servers
	require.Contains(t, asyncAPI.Servers, "development")
	devServer := asyncAPI.Servers["development"]
	assert.Equal(t, "localhost:9092", devServer.Host)
	assert.Equal(t, "kafka", devServer.Protocol)

	// Channels
	require.Contains(t, asyncAPI.Channels, "userEvents")
	channel := asyncAPI.Channels["userEvents"]
	assert.Equal(t, "user.events", channel.Address)
	assert.Equal(t, "Events related to users", channel.Description)

	// Operations
	require.Contains(t, asyncAPI.Operations, "sendUserEvent")
	op := asyncAPI.Operations["sendUserEvent"]
	assert.Equal(t, "send", op.Action)
	assert.Equal(t, "Send a user event", op.Summary)
	require.NotNil(t, op.Channel)
	assert.Equal(t, "#/channels/userEvents", op.Channel.Ref)

	// Messages
	require.Contains(t, asyncAPI.Components.Messages, "UserEvent")
	msg := asyncAPI.Components.Messages["UserEvent"]
	assert.Equal(t, "UserEvent", msg.Name)
	assert.Equal(t, "A user event message", msg.Description)
	require.NotNil(t, msg.Payload)

	// Schemas
	require.Contains(t, asyncAPI.Components.Schemas, "UserEvent")
}
