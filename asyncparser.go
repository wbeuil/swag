package swag

import (
	"strings"

	"github.com/swaggo/swag/v2/asyncapi"
)

// parseGeneralAsyncAPIInfo parses AsyncAPI general API info comments.
func (parser *Parser) parseGeneralAsyncAPIInfo(comments []string) error {
	previousAttribute := ""
	var currentServerName string

	for line := 0; line < len(comments); line++ {
		commentLine := comments[line]
		commentLine = strings.TrimSpace(commentLine)
		if len(commentLine) == 0 {
			continue
		}

		fields := FieldsByAnySpace(commentLine, 2)
		attribute := fields[0]
		var value string
		if len(fields) > 1 {
			value = fields[1]
		}

		switch attr := strings.ToLower(attribute); attr {
		case asyncTitleAttr:
			parser.asyncAPI.Info.Title = value
		case asyncVersionAttr:
			parser.asyncAPI.Info.Version = value
		case asyncDescriptionAttr:
			if previousAttribute == attribute {
				parser.asyncAPI.Info.Description += "\n" + value
				continue
			}
			parser.asyncAPI.Info.Description = value
		case asyncDefaultContentTypeAttr:
			parser.asyncAPI.DefaultContentType = value
		case asyncLicenseNameAttr:
			if parser.asyncAPI.Info.License == nil {
				parser.asyncAPI.Info.License = &asyncapi.License{}
			}
			parser.asyncAPI.Info.License.Name = value
		case asyncLicenseURLAttr:
			if parser.asyncAPI.Info.License == nil {
				parser.asyncAPI.Info.License = &asyncapi.License{}
			}
			parser.asyncAPI.Info.License.URL = value
		case "@asyncserver.name":
			currentServerName = value
			if parser.asyncAPI.Servers == nil {
				parser.asyncAPI.Servers = make(map[string]*asyncapi.Server)
			}
			parser.asyncAPI.Servers[currentServerName] = &asyncapi.Server{}
		case "@asyncserver.host":
			if currentServerName != "" {
				parser.asyncAPI.Servers[currentServerName].Host = value
			}
		case "@asyncserver.protocol":
			if currentServerName != "" {
				parser.asyncAPI.Servers[currentServerName].Protocol = value
			}
		case "@asyncserver.description":
			if currentServerName != "" {
				parser.asyncAPI.Servers[currentServerName].Description = value
			}
		case "@asyncserver.protocolversion":
			if currentServerName != "" {
				parser.asyncAPI.Servers[currentServerName].ProtocolVersion = value
			}
		case "@asyncserver.pathname":
			if currentServerName != "" {
				parser.asyncAPI.Servers[currentServerName].PathName = value
			}
		case asyncExternalDocsDescAttr:
			if parser.asyncAPI.ExternalDocs == nil {
				parser.asyncAPI.ExternalDocs = &asyncapi.ExternalDocs{}
			}
			parser.asyncAPI.ExternalDocs.Description = value
		case asyncExternalDocsURLAttr:
			if parser.asyncAPI.ExternalDocs == nil {
				parser.asyncAPI.ExternalDocs = &asyncapi.ExternalDocs{}
			}
			parser.asyncAPI.ExternalDocs.URL = value
		default:
			if currentServerName != "" {
				if rest, ok := trimAttrPrefix(attribute, "@asyncserver.binding"); ok {
					protocol := rest
					fieldPath := []string(nil)
					if idx := strings.Index(rest, "."); idx >= 0 {
						protocol = rest[:idx]
						if rest[idx+1:] != "" {
							fieldPath = strings.Split(rest[idx+1:], ".")
						}
					}
					bindingValue, err := parseBindingValue(value)
					if err != nil {
						return err
					}
					if err := assignBinding(parser.serverBindingMap(currentServerName), protocol, fieldPath, bindingValue); err != nil {
						return err
					}
				}
			}
		}

		previousAttribute = attribute
	}

	return nil
}

func (parser *Parser) serverBindingMap(name string) map[string]interface{} {
	server := parser.asyncAPI.Servers[name]
	if server.Bindings == nil {
		server.Bindings = make(map[string]interface{})
	}

	return server.Bindings
}
