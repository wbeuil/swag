package swag

import (
	"strings"

	"github.com/swaggo/swag/v2/asyncapi"
)

// parseGeneralAsyncAPIInfo parses AsyncAPI general API info comments.
func (parser *Parser) parseGeneralAsyncAPIInfo(comments []string) error {
	if isAsyncOperationCommentBlock(comments) {
		return nil
	}

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
		case asyncIdAttr:
			parser.asyncAPI.Id = value
		case asyncDescriptionAttr:
			if previousAttribute == attribute {
				parser.asyncAPI.Info.Description += "\n" + value
				continue
			}
			parser.asyncAPI.Info.Description = value
		case asyncDefaultContentTypeAttr:
			parser.asyncAPI.DefaultContentType = value
		case asyncContactNameAttr:
			parser.ensureAsyncContact().Name = value
		case asyncContactURLAttr:
			parser.ensureAsyncContact().URL = value
		case asyncContactEmailAttr:
			parser.ensureAsyncContact().Email = value
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
		case asyncTagAttr:
			tagFields := FieldsByAnySpace(value, 2)
			tag := &asyncapi.Tag{Name: tagFields[0]}
			if len(tagFields) > 1 {
				tag.Description = tagFields[1]
			}
			parser.asyncAPI.Info.Tags = append(parser.asyncAPI.Info.Tags, tag)
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
			parser.ensureAsyncInfoExternalDocs().Description = value
		case asyncExternalDocsURLAttr:
			parser.ensureAsyncInfoExternalDocs().URL = value
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

func (parser *Parser) ensureAsyncContact() *asyncapi.Contact {
	if parser.asyncAPI.Info.Contact == nil {
		parser.asyncAPI.Info.Contact = &asyncapi.Contact{}
	}

	return parser.asyncAPI.Info.Contact
}

func (parser *Parser) ensureAsyncInfoExternalDocs() *asyncapi.ExternalDocs {
	if parser.asyncAPI.Info.ExternalDocs == nil {
		parser.asyncAPI.Info.ExternalDocs = &asyncapi.ExternalDocs{}
	}

	return parser.asyncAPI.Info.ExternalDocs
}

func isAsyncOperationCommentBlock(comments []string) bool {
	for _, comment := range comments {
		comment = strings.TrimSpace(comment)
		if comment == "" {
			continue
		}
		attr := strings.ToLower(FieldsByAnySpace(comment, 2)[0])
		switch attr {
		case asyncOperationAttr, "@asyncoperation.name", asyncActionAttr, asyncChannelAttr, "@asyncchannel.name", asyncMessageAttr:
			return true
		}
	}

	return false
}

func (parser *Parser) serverBindingMap(name string) map[string]interface{} {
	server := parser.asyncAPI.Servers[name]
	if server.Bindings == nil {
		server.Bindings = make(map[string]interface{})
	}

	return server.Bindings
}
