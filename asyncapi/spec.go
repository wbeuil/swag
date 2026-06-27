package asyncapi

import (
	"encoding/json"

	openapi "github.com/sv-tools/openapi/spec"
)

// AsyncAPI is the root document object for the AsyncAPI specification.
type AsyncAPI struct {
	Asyncapi           string                `json:"asyncapi"`
	Id                 string                `json:"id,omitempty"`
	Info               *Info                 `json:"info"`
	Servers            map[string]*Server    `json:"servers,omitempty"`
	DefaultContentType string                `json:"defaultContentType,omitempty"`
	Channels           map[string]*Channel   `json:"channels"`
	Operations         map[string]*Operation `json:"operations"`
	Components         *Components           `json:"components,omitempty"`
	Tags               []*Tag                `json:"tags,omitempty"`
	ExternalDocs       *ExternalDocs         `json:"externalDocs,omitempty"`
}

// Info provides metadata about the API.
type Info struct {
	Title          string   `json:"title,omitempty"`
	Version        string   `json:"version,omitempty"`
	Description    string   `json:"description,omitempty"`
	TermsOfService string   `json:"termsOfService,omitempty"`
	Contact        *Contact `json:"contact,omitempty"`
	License        *License `json:"license,omitempty"`
}

// Contact information for the exposed API.
type Contact struct {
	Name  string `json:"name,omitempty"`
	URL   string `json:"url,omitempty"`
	Email string `json:"email,omitempty"`
}

// License information for the exposed API.
type License struct {
	Name string `json:"name,omitempty"`
	URL  string `json:"url,omitempty"`
}

// Server represents a message broker, a web server, or any other kind of computer program capable of sending and/or receiving data.
type Server struct {
	Host            string                     `json:"host"`
	Protocol        string                     `json:"protocol"`
	ProtocolVersion string                     `json:"protocolVersion,omitempty"`
	PathName        string                     `json:"pathname,omitempty"`
	Description     string                     `json:"description,omitempty"`
	Title           string                     `json:"title,omitempty"`
	Variables       map[string]*ServerVariable `json:"variables,omitempty"`
	Security        []SecurityRequirement      `json:"security,omitempty"`
	Tags            []*Tag                     `json:"tags,omitempty"`
	Bindings        map[string]interface{}     `json:"bindings,omitempty"`
}

// ServerVariable represents a Server Variable for server URL template substitution.
type ServerVariable struct {
	Enum        []string `json:"enum,omitempty"`
	Default     string   `json:"default,omitempty"`
	Description string   `json:"description,omitempty"`
	Examples    []string `json:"examples,omitempty"`
}

// SecurityRequirement lists the required security schemes to execute a specific operation.
type SecurityRequirement map[string][]string

// Channel represents a message broker channel.
type Channel struct {
	Address     string                 `json:"address,omitempty"`
	Messages    map[string]*Message    `json:"messages,omitempty"`
	Title       string                 `json:"title,omitempty"`
	Description string                 `json:"description,omitempty"`
	Parameters  map[string]*Parameter  `json:"parameters,omitempty"`
	Tags        []*Tag                 `json:"tags,omitempty"`
	Bindings    map[string]interface{} `json:"bindings,omitempty"`
}

// Parameter describes a Channel parameter.
type Parameter struct {
	Description string   `json:"description,omitempty"`
	Location    string   `json:"location,omitempty"`
	Enum        []string `json:"enum,omitempty"`
	Default     string   `json:"default,omitempty"`
	Examples    []string `json:"examples,omitempty"`
}

// Operation describes a send or receive operation.
type Operation struct {
	Action       string                 `json:"action,omitempty"` // send or receive
	Channel      *Reference             `json:"channel,omitempty"`
	Title        string                 `json:"title,omitempty"`
	Summary      string                 `json:"summary,omitempty"`
	Description  string                 `json:"description,omitempty"`
	Security     []SecurityRequirement  `json:"security,omitempty"`
	Tags         []*Tag                 `json:"tags,omitempty"`
	Traits       []*OperationTrait      `json:"traits,omitempty"`
	Messages     []*Message             `json:"messages,omitempty"`
	Reply        *OperationReply        `json:"reply,omitempty"`
	Bindings     map[string]interface{} `json:"bindings,omitempty"`
	ExternalDocs *ExternalDocs          `json:"externalDocs,omitempty"`
}

// OperationTrait is a reusable trait for an Operation.
type OperationTrait struct {
	Title       string                 `json:"title,omitempty"`
	Summary     string                 `json:"summary,omitempty"`
	Description string                 `json:"description,omitempty"`
	Security    []SecurityRequirement  `json:"security,omitempty"`
	Tags        []*Tag                 `json:"tags,omitempty"`
	Bindings    map[string]interface{} `json:"bindings,omitempty"`
}

// OperationReply describes the reply part of an operation.
type OperationReply struct {
	Channel  *Reference             `json:"channel,omitempty"`
	Messages []*Message             `json:"messages,omitempty"`
	Address  *OperationReplyAddress `json:"address,omitempty"`
}

// OperationReplyAddress specifies the address for the reply.
type OperationReplyAddress struct {
	Description string `json:"description,omitempty"`
	Location    string `json:"location,omitempty"`
}

// Message describes a message.
type Message struct {
	Name          string                 `json:"name,omitempty"`
	Title         string                 `json:"title,omitempty"`
	Summary       string                 `json:"summary,omitempty"`
	Description   string                 `json:"description,omitempty"`
	ContentType   string                 `json:"contentType,omitempty"`
	SchemaFormat  string                 `json:"schemaFormat,omitempty"`
	CorrelationID *Reference             `json:"correlationId,omitempty"`
	Headers       *MultiFormatSchema     `json:"headers,omitempty"`
	Payload       *MultiFormatSchema     `json:"payload,omitempty"`
	Tags          []*Tag                 `json:"tags,omitempty"`
	Examples      []*MessageExample      `json:"examples,omitempty"`
	Traits        []*MessageTrait        `json:"traits,omitempty"`
	Bindings      map[string]interface{} `json:"bindings,omitempty"`
	Ref           string                 `json:"$ref,omitempty"`
}

// MessageTrait is a reusable trait for a Message.
type MessageTrait struct {
	Name          string                 `json:"name,omitempty"`
	Title         string                 `json:"title,omitempty"`
	Summary       string                 `json:"summary,omitempty"`
	Description   string                 `json:"description,omitempty"`
	ContentType   string                 `json:"contentType,omitempty"`
	SchemaFormat  string                 `json:"schemaFormat,omitempty"`
	CorrelationID *Reference             `json:"correlationId,omitempty"`
	Headers       *MultiFormatSchema     `json:"headers,omitempty"`
	Tags          []*Tag                 `json:"tags,omitempty"`
	Bindings      map[string]interface{} `json:"bindings,omitempty"`
}

// MessageExample represents an example of a message.
type MessageExample struct {
	Name    string      `json:"name,omitempty"`
	Summary string      `json:"summary,omitempty"`
	Headers interface{} `json:"headers,omitempty"`
	Payload interface{} `json:"payload,omitempty"`
}

// Components holds a set of reusable objects for different aspects of the AsyncAPI Specification.
type Components struct {
	Schemas           map[string]*openapi.Schema        `json:"schemas,omitempty"`
	Servers           map[string]*Server                `json:"servers,omitempty"`
	Channels          map[string]*Channel               `json:"channels,omitempty"`
	Messages          map[string]*Message               `json:"messages,omitempty"`
	SecuritySchemes   map[string]*SecurityScheme        `json:"securitySchemes,omitempty"`
	ServerVariables   map[string]*ServerVariable        `json:"serverVariables,omitempty"`
	Parameters        map[string]*Parameter             `json:"parameters,omitempty"`
	CorrelationIds    map[string]*CorrelationID         `json:"correlationIds,omitempty"`
	Replies           map[string]*OperationReply        `json:"replies,omitempty"`
	ReplyAddresses    map[string]*OperationReplyAddress `json:"replyAddresses,omitempty"`
	Operations        map[string]*Operation             `json:"operations,omitempty"`
	OperationTraits   map[string]*OperationTrait        `json:"operationTraits,omitempty"`
	MessageTraits     map[string]*MessageTrait          `json:"messageTraits,omitempty"`
	ServerBindings    map[string]interface{}            `json:"serverBindings,omitempty"`
	ChannelBindings   map[string]interface{}            `json:"channelBindings,omitempty"`
	OperationBindings map[string]interface{}            `json:"operationBindings,omitempty"`
	MessageBindings   map[string]interface{}            `json:"messageBindings,omitempty"`
}

// Tag represents metadata about an API.
type Tag struct {
	Name         string        `json:"name,omitempty"`
	Description  string        `json:"description,omitempty"`
	ExternalDocs *ExternalDocs `json:"externalDocs,omitempty"`
}

// ExternalDocs allows referencing an external resource for extended documentation.
type ExternalDocs struct {
	Description string `json:"description,omitempty"`
	URL         string `json:"url,omitempty"`
}

// Reference is a simple reference object using $ref.
type Reference struct {
	Ref string `json:"$ref,omitempty"`
}

// SecurityScheme defines a security scheme that can be used by the operations.
type SecurityScheme struct {
	Type             string      `json:"type,omitempty"`
	Description      string      `json:"description,omitempty"`
	Name             string      `json:"name,omitempty"`
	In               string      `json:"in,omitempty"`
	Scheme           string      `json:"scheme,omitempty"`
	BearerFormat     string      `json:"bearerFormat,omitempty"`
	Flows            *OAuthFlows `json:"flows,omitempty"`
	OpenIDConnectURL string      `json:"openIdConnectUrl,omitempty"`
}

// OAuthFlows allows configuration of the supported OAuth Flows.
type OAuthFlows struct {
	Implicit          *OAuthFlow `json:"implicit,omitempty"`
	Password          *OAuthFlow `json:"password,omitempty"`
	ClientCredentials *OAuthFlow `json:"clientCredentials,omitempty"`
	AuthorizationCode *OAuthFlow `json:"authorizationCode,omitempty"`
}

// OAuthFlow provides configuration details for a supported OAuth Flow.
type OAuthFlow struct {
	AuthorizationURL string            `json:"authorizationUrl,omitempty"`
	TokenURL         string            `json:"tokenUrl,omitempty"`
	RefreshURL       string            `json:"refreshUrl,omitempty"`
	Scopes           map[string]string `json:"scopes,omitempty"`
}

// CorrelationID specifies an identifier to correlate messages.
type CorrelationID struct {
	Description string `json:"description,omitempty"`
	Location    string `json:"location,omitempty"`
}

// MultiFormatSchema allows the use of different schema formats.
type MultiFormatSchema struct {
	Schema       *openapi.Schema `json:"schema,omitempty"`
	SchemaFormat string          `json:"schemaFormat,omitempty"`
}

// NewAsyncAPI creates a new AsyncAPI document with defaults.
func NewAsyncAPI() *AsyncAPI {
	return &AsyncAPI{
		Asyncapi:   "3.1.0",
		Info:       &Info{},
		Servers:    make(map[string]*Server),
		Channels:   make(map[string]*Channel),
		Operations: make(map[string]*Operation),
		Components: &Components{
			Schemas:           make(map[string]*openapi.Schema),
			Servers:           make(map[string]*Server),
			Channels:          make(map[string]*Channel),
			Messages:          make(map[string]*Message),
			SecuritySchemes:   make(map[string]*SecurityScheme),
			ServerVariables:   make(map[string]*ServerVariable),
			Parameters:        make(map[string]*Parameter),
			CorrelationIds:    make(map[string]*CorrelationID),
			Replies:           make(map[string]*OperationReply),
			ReplyAddresses:    make(map[string]*OperationReplyAddress),
			Operations:        make(map[string]*Operation),
			OperationTraits:   make(map[string]*OperationTrait),
			MessageTraits:     make(map[string]*MessageTrait),
			ServerBindings:    make(map[string]interface{}),
			ChannelBindings:   make(map[string]interface{}),
			OperationBindings: make(map[string]interface{}),
			MessageBindings:   make(map[string]interface{}),
		},
	}
}

// MarshalJSON implements custom JSON marshalling.
func (a *AsyncAPI) MarshalJSON() ([]byte, error) {
	type alias AsyncAPI
	return json.Marshal((*alias)(a))
}
