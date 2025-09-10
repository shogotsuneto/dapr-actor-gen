package parser

import (
	"fmt"
	"io/ioutil"
	"path/filepath"
	"strings"

	"github.com/getkin/kin-openapi/openapi3"
	"gopkg.in/yaml.v3"
)

// NewParserFromFile creates the appropriate parser based on file content and extension
func NewParserFromFile(filePath string) (Parser, error) {
	content, err := ioutil.ReadFile(filePath)
	if err != nil {
		return nil, fmt.Errorf("failed to read file %s: %v", filePath, err)
	}

	// Try to determine format from content and extension
	if isActorYAMLFormat(content, filePath) {
		return NewActorYAMLParser(content)
	}

	// Default to OpenAPI format
	loader := openapi3.NewLoader()
	doc, err := loader.LoadFromFile(filePath)
	if err != nil {
		return nil, fmt.Errorf("failed to load as OpenAPI spec: %v", err)
	}

	return NewOpenAPIParser(doc), nil
}

// isActorYAMLFormat attempts to determine if the YAML file is in actor format
func isActorYAMLFormat(content []byte, filePath string) bool {
	// Parse YAML and check for actor-specific structure
	var data map[string]interface{}
	if err := yaml.Unmarshal(content, &data); err != nil {
		return false
	}

	// Check for OpenAPI-specific keys first (if present, it's definitely OpenAPI)
	if _, hasOpenAPI := data["openapi"]; hasOpenAPI {
		return false
	}
	if _, hasSwagger := data["swagger"]; hasSwagger {
		return false
	}
	
	// Check for actor-specific keys (if present, it's definitely Actor YAML)
	if _, hasActors := data["actors"]; hasActors {
		return true
	}

	// Check if paths follow OpenAPI actor pattern (if so, it's OpenAPI)
	if paths, hasPaths := data["paths"]; hasPaths {
		if pathsMap, ok := paths.(map[string]interface{}); ok {
			for path := range pathsMap {
				if strings.Contains(path, "/{actorId}/method/") {
					return false // This is OpenAPI actor format
				}
			}
		}
		// Has paths but not actor pattern - likely OpenAPI
		return false
	}

	// Check file extension/name patterns for hints as fallback
	filename := strings.ToLower(filepath.Base(filePath))
	if strings.Contains(filename, "actor") && !strings.Contains(filename, "openapi") {
		return true
	}

	// If no clear indicators, default to OpenAPI (safer for backward compatibility)
	return false
}