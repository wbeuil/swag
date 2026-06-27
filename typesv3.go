package swag

import "github.com/wbeuil/openapi"

// SchemaV3 parsed schema.
type SchemaV3 struct {
	*openapi.Schema        //
	PkgPath         string // package import path used to rename Name of a definition int case of conflict
	Name            string // Name in definitions
}
