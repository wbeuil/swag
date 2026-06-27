package swag

import (
	"encoding/json"
	"fmt"
	"go/ast"
	"regexp"
	"strings"

	"github.com/swaggo/swag/v2/asyncapi"
)

// AsyncOperation represents an AsyncAPI operation being built from comments.
type AsyncOperation struct {
	parser        *Parser
	Operation     *asyncapi.Operation
	OperationName string
	Channel       *asyncapi.Channel
	ChannelName   string
	Messages      []*asyncapi.Message
	Tags          []*asyncapi.Tag
	Server        *asyncapi.Server
	ServerName    string
}

// NewAsyncOperation creates a new AsyncOperation.
func NewAsyncOperation(parser *Parser) *AsyncOperation {
	return &AsyncOperation{
		parser:    parser,
		Operation: &asyncapi.Operation{},
		Channel:   &asyncapi.Channel{},
		Messages:  make([]*asyncapi.Message, 0),
		Tags:      make([]*asyncapi.Tag, 0),
	}
}

// ParseComment parses an AsyncAPI comment.
func (ao *AsyncOperation) ParseComment(comment string, astFile *ast.File) error {
	commentLine := strings.TrimSpace(strings.TrimLeft(comment, "/"))
	if len(commentLine) == 0 {
		return nil
	}

	fields := FieldsByAnySpace(commentLine, 2)
	attribute := fields[0]
	lowerAttribute := strings.ToLower(attribute)
	var lineRemainder string
	if len(fields) > 1 {
		lineRemainder = fields[1]
	}

	switch lowerAttribute {
	case asyncapiAttr:
		// no-op, just marks as asyncapi
	case asyncOperationAttr:
		ao.OperationName = lineRemainder
	case "@asyncoperation.name":
		ao.OperationName = lineRemainder
	case asyncActionAttr:
		ao.Operation.Action = lineRemainder
	case asyncChannelAttr:
		ao.ChannelName = lineRemainder
	case "@asyncchannel.name":
		ao.ChannelName = lineRemainder
	case "@asyncchannel.address":
		ao.Channel.Address = lineRemainder
	case "@asyncchannel.description":
		ao.Channel.Description = lineRemainder
	case asyncMessageAttr:
		return ao.ParseMessageComment(lineRemainder, astFile)
	case "@asyncmessage.name":
		if len(ao.Messages) > 0 {
			ao.Messages[len(ao.Messages)-1].Name = lineRemainder
		}
	case "@asyncmessage.contenttype":
		if len(ao.Messages) > 0 {
			ao.Messages[len(ao.Messages)-1].ContentType = lineRemainder
		}
	case asyncSummaryAttr:
		ao.Operation.Summary = lineRemainder
	case asyncDescriptionAttr:
		if ao.Operation.Description == "" {
			ao.Operation.Description = lineRemainder
		} else {
			ao.Operation.Description += "\n" + lineRemainder
		}
	case asyncTagAttr:
		return ao.ParseTagComment(lineRemainder)
	case asyncBindingAttr:
		return ao.ParseBindingComment(lineRemainder)
	case asyncServerAttr:
		if ao.Server == nil {
			ao.Server = &asyncapi.Server{}
		}
		ao.ServerName = lineRemainder
	case "@asyncserver.name":
		if ao.Server == nil {
			ao.Server = &asyncapi.Server{}
		}
		ao.ServerName = lineRemainder
	case "@asyncserver.host":
		if ao.Server != nil {
			ao.Server.Host = lineRemainder
		}
	case "@asyncserver.protocol":
		if ao.Server != nil {
			ao.Server.Protocol = lineRemainder
		}
	case "@asyncserver.description":
		if ao.Server != nil {
			ao.Server.Description = lineRemainder
		}
	case asyncExternalDocsDescAttr:
		if ao.Operation.ExternalDocs == nil {
			ao.Operation.ExternalDocs = &asyncapi.ExternalDocs{}
		}
		ao.Operation.ExternalDocs.Description = lineRemainder
	case asyncExternalDocsURLAttr:
		if ao.Operation.ExternalDocs == nil {
			ao.Operation.ExternalDocs = &asyncapi.ExternalDocs{}
		}
		ao.Operation.ExternalDocs.URL = lineRemainder
	default:
		if strings.HasPrefix(lowerAttribute, "@asyncx-") {
			return ao.ParseMetadataExtension(attribute, lowerAttribute, lineRemainder)
		}
	}

	return nil
}

// ParseMessageComment parses a message declaration comment.
func (ao *AsyncOperation) ParseMessageComment(lineRemainder string, astFile *ast.File) error {
	re := regexp.MustCompile(`^([\w{}]+)\s+([\w\-.\\{}=,\[\s\]]+)\s*("[^"]*(?:"|$))?\s*([\w/+-]+)?$`)
	matches := re.FindStringSubmatch(lineRemainder)
	if len(matches) < 3 {
		return fmt.Errorf("can not parse async message comment %s", lineRemainder)
	}

	typeName := strings.TrimSpace(matches[2])
	description := strings.Trim(matches[3], `"`)
	contentType := matches[4]

	message := &asyncapi.Message{
		Name:        typeName,
		Description: description,
		ContentType: contentType,
	}

	schema, err := ao.parser.getTypeSchemaV3(typeName, astFile, true)
	if err != nil {
		return err
	}

	message.Payload = &asyncapi.MultiFormatSchema{}

	if schema.Spec != nil {
		message.Payload.Schema = schema.Spec
	} else if schema.Ref != nil {
		// Complex type returned as a ref; look up the actual schema from OpenAPI components
		refName := strings.TrimPrefix(schema.Ref.Ref, "#/components/schemas/")
		if openAPISchema, ok := ao.parser.openAPI.Components.Spec.Schemas[refName]; ok && openAPISchema.Spec != nil {
			copiedSchema := *openAPISchema.Spec
			message.Payload.Schema = &copiedSchema
		}
	}

	ao.Messages = append(ao.Messages, message)
	return nil
}

// ParseTagComment parses a tag comment.
func (ao *AsyncOperation) ParseTagComment(lineRemainder string) error {
	fields := FieldsByAnySpace(lineRemainder, 2)
	tag := &asyncapi.Tag{Name: fields[0]}
	if len(fields) > 1 {
		tag.Description = fields[1]
	}
	ao.Tags = append(ao.Tags, tag)
	return nil
}

// ParseBindingComment parses a binding comment.
func (ao *AsyncOperation) ParseBindingComment(lineRemainder string) error {
	// Format: @asyncBinding.kafka.groupId my-group
	// or @asyncBinding.mqtt.qos 1
	if ao.Operation.Bindings == nil {
		ao.Operation.Bindings = make(map[string]interface{})
	}
	// Store raw for now; refinement in future
	return nil
}

// ParseMetadataExtension parses async x- extensions.
func (ao *AsyncOperation) ParseMetadataExtension(attribute, lowerAttribute, lineRemainder string) error {
	if len(lineRemainder) == 0 {
		return fmt.Errorf("annotation %s needs a value", attribute)
	}
	var valueJSON any
	if err := json.Unmarshal([]byte(lineRemainder), &valueJSON); err != nil {
		return fmt.Errorf("annotation %s can not parse value: %v", attribute, err)
	}
	if ao.Operation.Bindings == nil {
		ao.Operation.Bindings = make(map[string]interface{})
	}
	return nil
}
