package integration

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/getkin/kin-openapi/openapi3"
	"github.com/shogotsuneto/dapr-actor-gen/pkg/generator"
	"github.com/shogotsuneto/dapr-actor-gen/pkg/parser"
)

func TestEnumGeneration(t *testing.T) {
	// Load the type alias test OpenAPI spec which contains UserStatus enum
	loader := openapi3.NewLoader()
	doc, err := loader.LoadFromFile("testdata/type-alias.yaml")
	if err != nil {
		t.Fatalf("Failed to load type alias OpenAPI spec: %v", err)
	}

	// Parse the spec to intermediate model
	p := parser.NewOpenAPIParser(doc)
	model, err := p.Parse()
	if err != nil {
		t.Fatalf("Failed to parse OpenAPI spec: %v", err)
	}

	// Generate code and verify enum generation
	outputDir := "test-output/enum-generation"
	gen := &generator.Generator{}
	err = gen.GenerateActorPackages(model, outputDir, generator.GenerationOptions{
		GenerateImpl: true,
	})
	if err != nil {
		t.Fatalf("Failed to generate actor packages: %v", err)
	}

	// Read the generated types.go file
	typesFile := filepath.Join(outputDir, "user", "types.go")
	content, err := os.ReadFile(typesFile)
	if err != nil {
		t.Fatalf("Failed to read generated types file: %v", err)
	}

	contentStr := string(content)

	// Verify enum type definition is generated
	if !strings.Contains(contentStr, "type UserStatus string") {
		t.Error("Expected enum type definition 'type UserStatus string' not found")
	}

	// Verify enum constants are generated
	expectedConstants := []string{
		"UserStatusActive UserStatus = \"active\"",
		"UserStatusInactive UserStatus = \"inactive\"",
		"UserStatusSuspended UserStatus = \"suspended\"",
		"UserStatusPending UserStatus = \"pending\"",
	}

	for _, constant := range expectedConstants {
		if !strings.Contains(contentStr, constant) {
			t.Errorf("Expected enum constant '%s' not found", constant)
		}
	}

	// Verify enum methods are generated
	expectedMethods := []string{
		"func (e UserStatus) IsValid() bool",
		"func (e UserStatus) String() string",
		"func AllUserStatusValues() []UserStatus",
	}

	for _, method := range expectedMethods {
		if !strings.Contains(contentStr, method) {
			t.Errorf("Expected enum method '%s' not found", method)
		}
	}

	// Verify that switch cases are generated for validation
	expectedSwitchCases := []string{
		"case UserStatusActive:",
		"case UserStatusInactive:",
		"case UserStatusSuspended:",
		"case UserStatusPending:",
	}

	for _, switchCase := range expectedSwitchCases {
		if !strings.Contains(contentStr, switchCase) {
			t.Errorf("Expected switch case '%s' not found", switchCase)
		}
	}

	// Verify that the AllValues function returns all constants
	expectedInAllValues := []string{
		"UserStatusActive,",
		"UserStatusInactive,",
		"UserStatusSuspended,",
		"UserStatusPending,",
	}

	for _, valueInList := range expectedInAllValues {
		if !strings.Contains(contentStr, valueInList) {
			t.Errorf("Expected value in AllValues function '%s' not found", valueInList)
		}
	}

	// Clean up
	os.RemoveAll(outputDir)
}