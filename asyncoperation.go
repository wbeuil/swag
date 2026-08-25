package swag

import (
	"encoding/json"
	"fmt"
	"go/ast"
	"regexp"
	"strconv"
	"strings"

	"github.com/swaggo/swag/v2/asyncapi"
)

var bindingTargets = map[string]struct{}{
	"server":    {},
	"channel":   {},
	"operation": {},
	"message":   {},
}

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
		return ao.ParseBindingComment(attribute, lineRemainder)
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
		if strings.HasPrefix(lowerAttribute, asyncBindingAttr+".") {
			return ao.ParseBindingComment(attribute, lineRemainder)
		}
		if strings.HasPrefix(lowerAttribute, "@asyncx-") {
			return ao.ParseMetadataExtension(attribute, lineRemainder)
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
		Description: description,
		ContentType: contentType,
	}
	if !IsGolangPrimitiveType(typeName) {
		message.Name = typeName
	}

	schema, err := ao.parser.getTypeSchemaV3(typeName, astFile, true)
	if err != nil {
		return err
	}

	message.Payload = &asyncapi.MultiFormatSchema{}
	if schema.Ref != nil && schema.Ref.Ref != "" {
		message.Payload.Ref = schema.Ref.Ref
	} else if schema.Spec != nil {
		message.Payload.Schema = schema.Spec
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

// ParseBindingComment parses a protocol binding comment.
//
// Supported forms:
//
//	@asyncBinding.<protocol> <json>
//	@asyncBinding.<protocol>.<field.path> <value>
//	@asyncBinding.<target>.<protocol> <json>
//	@asyncBinding.<target>.<protocol>.<field.path> <value>
//
// target is one of server, channel, operation, or message. Omitting it
// applies the binding to the current operation.
func (ao *AsyncOperation) ParseBindingComment(attribute, lineRemainder string) error {
	rest, ok := trimAttrPrefix(attribute, asyncBindingAttr)
	if !ok {
		return fmt.Errorf("can not parse async binding comment %s", attribute)
	}

	target, protocol, fieldPath, err := splitBindingPath(rest)
	if err != nil {
		return err
	}

	value, err := parseBindingValue(lineRemainder)
	if err != nil {
		return fmt.Errorf("annotation %s: %w", attribute, err)
	}

	bindings, err := ao.bindingMap(target)
	if err != nil {
		return err
	}

	return assignBinding(bindings, protocol, fieldPath, value)
}

func (ao *AsyncOperation) bindingMap(target string) (map[string]interface{}, error) {
	switch target {
	case "server":
		if ao.Server == nil {
			ao.Server = &asyncapi.Server{}
		}
		if ao.Server.Bindings == nil {
			ao.Server.Bindings = make(map[string]interface{})
		}

		return ao.Server.Bindings, nil
	case "channel":
		if ao.Channel.Bindings == nil {
			ao.Channel.Bindings = make(map[string]interface{})
		}

		return ao.Channel.Bindings, nil
	case "message":
		if len(ao.Messages) == 0 {
			return nil, fmt.Errorf("@asyncBinding.message requires a preceding @asyncMessage")
		}
		msg := ao.Messages[len(ao.Messages)-1]
		if msg.Bindings == nil {
			msg.Bindings = make(map[string]interface{})
		}

		return msg.Bindings, nil
	default:
		if ao.Operation.Bindings == nil {
			ao.Operation.Bindings = make(map[string]interface{})
		}

		return ao.Operation.Bindings, nil
	}
}

func splitBindingPath(rest string) (target, protocol string, fieldPath []string, err error) {
	if rest == "" {
		return "", "", nil, fmt.Errorf("binding protocol is required")
	}

	parts := strings.Split(rest, ".")
	target = "operation"
	if _, ok := bindingTargets[strings.ToLower(parts[0])]; ok {
		target = strings.ToLower(parts[0])
		parts = parts[1:]
	}
	if len(parts) == 0 || parts[0] == "" {
		return "", "", nil, fmt.Errorf("binding protocol is required")
	}

	return target, parts[0], parts[1:], nil
}

func parseBindingValue(raw string) (interface{}, error) {
	raw = strings.TrimSpace(raw)
	if raw == "" {
		return nil, fmt.Errorf("binding value is empty")
	}
	if raw[0] == '{' || raw[0] == '[' {
		var value interface{}
		if err := json.Unmarshal([]byte(raw), &value); err != nil {
			return nil, err
		}

		return value, nil
	}
	switch raw {
	case "true":
		return true, nil
	case "false":
		return false, nil
	}
	if i, err := strconv.ParseInt(raw, 10, 64); err == nil {
		return i, nil
	}
	if f, err := strconv.ParseFloat(raw, 64); err == nil {
		return f, nil
	}

	return raw, nil
}

func assignBinding(bindings map[string]interface{}, protocol string, fieldPath []string, value interface{}) error {
	if protocol == "" {
		return fmt.Errorf("binding protocol is required")
	}
	if len(fieldPath) == 0 {
		if obj, ok := value.(map[string]interface{}); ok {
			existing, _ := bindings[protocol].(map[string]interface{})
			if existing == nil {
				bindings[protocol] = obj
				return nil
			}
			for key, nested := range obj {
				existing[key] = nested
			}

			return nil
		}
		bindings[protocol] = value

		return nil
	}

	existing, _ := bindings[protocol].(map[string]interface{})
	if existing == nil {
		existing = make(map[string]interface{})
		bindings[protocol] = existing
	}
	setNestedMap(existing, fieldPath, value)

	return nil
}

func setNestedMap(current map[string]interface{}, path []string, value interface{}) {
	for i, key := range path {
		if i == len(path)-1 {
			current[key] = value
			return
		}
		next, ok := current[key].(map[string]interface{})
		if !ok {
			next = make(map[string]interface{})
			current[key] = next
		}
		current = next
	}
}

func trimAttrPrefix(attribute, prefix string) (string, bool) {
	if len(attribute) < len(prefix) || !strings.EqualFold(attribute[:len(prefix)], prefix) {
		return "", false
	}

	return strings.TrimPrefix(attribute[len(prefix):], "."), true
}

// ParseMetadataExtension parses async x- extensions.
func (ao *AsyncOperation) ParseMetadataExtension(attribute, lineRemainder string) error {
	if len(lineRemainder) == 0 {
		return fmt.Errorf("annotation %s needs a value", attribute)
	}
	var valueJSON any
	if err := json.Unmarshal([]byte(lineRemainder), &valueJSON); err != nil {
		return fmt.Errorf("annotation %s can not parse value: %v", attribute, err)
	}
	if ao.Operation.Extensions == nil {
		ao.Operation.Extensions = make(map[string]interface{})
	}

	key := "x-" + attribute[len("@asyncx-"):]
	ao.Operation.Extensions[key] = valueJSON

	return nil
}
