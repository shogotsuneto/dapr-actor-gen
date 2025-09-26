package parser

import (
	"fmt"
	"io/ioutil"

	"github.com/getkin/kin-openapi/openapi3"
)

// NewParser creates the appropriate parser based on the specified format
func NewParser(format, filePath string) (Parser, error) {
	switch format {
	case "openapi":
		return newOpenAPIParser(filePath)
	case "actor-schema":
		return newActorSchemaParser(filePath)
	default:
		return nil, fmt.Errorf("unsupported format: %s. Supported formats are 'openapi' and 'actor-schema'", format)
	}
}

// newOpenAPIParser creates an OpenAPI parser from file
func newOpenAPIParser(filePath string) (Parser, error) {
	loader := openapi3.NewLoader()
	doc, err := loader.LoadFromFile(filePath)
	if err != nil {
		return nil, fmt.Errorf("failed to load as OpenAPI spec: %v", err)
	}
	return NewOpenAPIParser(doc), nil
}

// newActorSchemaParser creates an Actor Schema parser from file
func newActorSchemaParser(filePath string) (Parser, error) {
	content, err := ioutil.ReadFile(filePath)
	if err != nil {
		return nil, fmt.Errorf("failed to read file %s: %v", filePath, err)
	}
	return NewActorSchemaParser(content)
}
