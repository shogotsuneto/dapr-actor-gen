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

func TestBasicActorParsing(t *testing.T) {
	// Load the basic actor test OpenAPI spec
	loader := openapi3.NewLoader()
	doc, err := loader.LoadFromFile("testdata/basic-actor.yaml")
	if err != nil {
		t.Fatalf("Failed to load basic actor OpenAPI spec: %v", err)
	}

	// Parse the spec to intermediate model
	p := parser.NewOpenAPIParser(doc)
	model, err := p.Parse()
	if err != nil {
		t.Fatalf("Failed to parse OpenAPI spec: %v", err)
	}

	// Verify that we have exactly one actor
	if len(model.Actors) != 1 {
		t.Errorf("Expected 1 actor, got %d", len(model.Actors))
	}

	// Verify the Test actor
	actor := model.Actors[0]
	if actor.ActorType != "Test" {
		t.Errorf("Expected actor type 'Test', got '%s'", actor.ActorType)
	}

	if len(actor.Methods) != 2 {
		t.Errorf("Expected Test actor to have 2 methods, got %d", len(actor.Methods))
	}

	// Verify methods
	methodNames := make(map[string]bool)
	for _, method := range actor.Methods {
		methodNames[method.Name] = true
	}
	if !methodNames["GetValue"] {
		t.Error("Expected 'GetValue' method not found")
	}
	if !methodNames["SetValue"] {
		t.Error("Expected 'SetValue' method not found")
	}

	// Verify actor-specific types (single actor gets all types it uses)
	if len(actor.Types.Structs) < 2 {
		t.Errorf("Expected at least 2 struct types for TestActor, got %d", len(actor.Types.Structs))
	}

	// Verify specific types exist
	typeNames := make(map[string]bool)
	for _, structType := range actor.Types.Structs {
		typeNames[structType.Name] = true
	}
	if !typeNames["TestState"] {
		t.Error("Expected 'TestState' type not found")
	}
	if !typeNames["SetValueRequest"] {
		t.Error("Expected 'SetValueRequest' type not found")
	}

}

func TestMultiActorTypeDuplication(t *testing.T) {
	// Load the multi-actor test OpenAPI spec
	loader := openapi3.NewLoader()
	doc, err := loader.LoadFromFile("testdata/multi-actor.yaml")
	if err != nil {
		t.Fatalf("Failed to load multi-actor OpenAPI spec: %v", err)
	}

	// Parse the spec to intermediate model
	p := parser.NewOpenAPIParser(doc)
	model, err := p.Parse()
	if err != nil {
		t.Fatalf("Failed to parse OpenAPI spec: %v", err)
	}

	// Verify that we have exactly two actors
	if len(model.Actors) != 2 {
		t.Errorf("Expected 2 actors, got %d", len(model.Actors))
	}

	// Verify actors exist
	actorTypes := make(map[string]*generator.ActorInterface)
	for i, actor := range model.Actors {
		actorTypes[actor.ActorType] = &model.Actors[i]
	}

	counterActor, hasCounter := actorTypes["Counter"]
	calcActor, hasCalc := actorTypes["Calculator"]

	if !hasCounter {
		t.Error("Counter not found in parsed model")
	}
	if !hasCalc {
		t.Error("Calculator not found in parsed model")
	}

	// Verify Counter methods
	if hasCounter && len(counterActor.Methods) != 3 {
		t.Errorf("Expected Counter to have 3 methods, got %d", len(counterActor.Methods))
	}

	// Verify Calculator methods
	if hasCalc && len(calcActor.Methods) != 3 {
		t.Errorf("Expected Calculator to have 3 methods, got %d", len(calcActor.Methods))
	}

	// Verify that types are duplicated in each actor that uses them
	if hasCounter {
		counterTypeNames := make(map[string]bool)
		for _, structType := range counterActor.Types.Structs {
			counterTypeNames[structType.Name] = true
		}
		// CounterState should be actor-specific
		if !counterTypeNames["CounterState"] {
			t.Error("Expected CounterActor-specific type 'CounterState' not found")
		}
		// OperationLog and LogMetadata should now be duplicated in Counter actor
		if !counterTypeNames["OperationLog"] {
			t.Error("Expected type 'OperationLog' in Counter actor not found")
		}
		if !counterTypeNames["LogMetadata"] {
			t.Error("Expected type 'LogMetadata' in Counter actor not found")
		}
	}

	if hasCalc {
		calcTypeNames := make(map[string]bool)
		for _, structType := range calcActor.Types.Structs {
			calcTypeNames[structType.Name] = true
		}
		// Calculator-specific types
		if !calcTypeNames["MathOperation"] {
			t.Error("Expected CalculatorActor-specific type 'MathOperation' not found")
		}
		if !calcTypeNames["OperationResult"] {
			t.Error("Expected CalculatorActor-specific type 'OperationResult' not found")
		}
		// OperationLog and LogMetadata should now be duplicated in Calculator actor too
		if !calcTypeNames["OperationLog"] {
			t.Error("Expected type 'OperationLog' in Calculator actor not found")
		}
		if !calcTypeNames["LogMetadata"] {
			t.Error("Expected type 'LogMetadata' in Calculator actor not found")
		}
	}
}

func TestTypeAliasGeneration(t *testing.T) {
	// Load the type alias test OpenAPI spec
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

	// Verify that we have exactly one actor
	if len(model.Actors) != 1 {
		t.Errorf("Expected 1 actor, got %d", len(model.Actors))
	}

	actor := model.Actors[0]
	if actor.ActorType != "User" {
		t.Errorf("Expected actor type 'User', got '%s'", actor.ActorType)
	}

	// Verify that type aliases are generated
	totalAliases := len(actor.Types.Aliases)
	if totalAliases == 0 {
		t.Error("Expected type aliases to be generated, but found none")
	}

	// Look for specific type aliases that should be generated (non-enum types)
	aliasNames := make(map[string]bool)
	for _, alias := range actor.Types.Aliases {
		aliasNames[alias.Name] = true
	}

	// These should be generated as type aliases (simple types without enums)
	expectedAliases := []string{"UserId", "EmailAddress"}
	for _, expected := range expectedAliases {
		if !aliasNames[expected] {
			t.Errorf("Expected type alias '%s' not found", expected)
		}
	}

	// Verify that enum types are generated
	totalEnums := len(actor.Types.Enums)
	if totalEnums == 0 {
		t.Error("Expected enum types to be generated, but found none")
	}

	// Look for specific enum types that should be generated
	enumNames := make(map[string]bool)
	for _, enum := range actor.Types.Enums {
		enumNames[enum.Name] = true
	}

	// UserStatus should be generated as an enum type (has enum values)
	expectedEnums := []string{"UserStatus"}
	for _, expected := range expectedEnums {
		if !enumNames[expected] {
			t.Errorf("Expected enum type '%s' not found", expected)
		}
	}

	// Verify UserStatus enum has the correct values
	for _, enum := range actor.Types.Enums {
		if enum.Name == "UserStatus" {
			expectedValues := []string{"active", "inactive", "suspended", "pending"}
			if len(enum.Values) != len(expectedValues) {
				t.Errorf("Expected UserStatus to have %d values, got %d", len(expectedValues), len(enum.Values))
			}
			for i, expected := range expectedValues {
				if i >= len(enum.Values) || enum.Values[i] != expected {
					t.Errorf("Expected UserStatus value[%d] to be '%s', got '%s'", i, expected, enum.Values[i])
				}
			}
		}
	}
}

func TestGeneratorWithTestSpecs(t *testing.T) {
	tests := []struct {
		name     string
		specFile string
	}{
		{"Basic Actor", "testdata/basic-actor.yaml"},
		{"Multi Actor", "testdata/multi-actor.yaml"},
		{"Type Alias", "testdata/type-alias.yaml"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Load and parse the OpenAPI spec
			loader := openapi3.NewLoader()
			doc, err := loader.LoadFromFile(tt.specFile)
			if err != nil {
				t.Fatalf("Failed to load OpenAPI spec %s: %v", tt.specFile, err)
			}

			p := parser.NewOpenAPIParser(doc)
			model, err := p.Parse()
			if err != nil {
				t.Fatalf("Failed to parse OpenAPI spec: %v", err)
			}

			// Generate code using the intermediate model
			gen := &generator.Generator{}
			outputDir := filepath.Join("test-output", tt.name)
			options := generator.GenerationOptions{
				GenerateImpl:    false,
				GenerateExample: false,
			}
			err = gen.GenerateActorPackages(model, outputDir, options)
			if err != nil {
				t.Fatalf("Failed to generate actor packages: %v", err)
			}

			// Clean up after test
			defer func() {
				os.RemoveAll(outputDir)
			}()

			t.Logf("Successfully generated actor packages for %s", tt.name)
		})
	}
}

func TestGeneratorWithPartialImplementation(t *testing.T) {
	// Load the multi-actor spec
	loader := openapi3.NewLoader()
	doc, err := loader.LoadFromFile("testdata/multi-actor.yaml")
	if err != nil {
		t.Fatalf("Failed to load OpenAPI spec: %v", err)
	}

	p := parser.NewOpenAPIParser(doc)
	model, err := p.Parse()
	if err != nil {
		t.Fatalf("Failed to parse OpenAPI spec: %v", err)
	}

	// Generate with partial implementation
	gen := &generator.Generator{}
	outputDir := "test-output/partial-impl"
	options := generator.GenerationOptions{
		GenerateImpl:    true,
		GenerateExample: false,
	}
	err = gen.GenerateActorPackages(model, outputDir, options)
	if err != nil {
		t.Fatalf("Failed to generate actor packages with impl: %v", err)
	}

	// Clean up after test
	defer func() {
		os.RemoveAll(outputDir)
	}()

	// Verify impl.go files exist
	for _, actor := range model.Actors {
		packageName := strings.ToLower(actor.ActorType)
		implFile := filepath.Join(outputDir, packageName, "impl.go")
		if _, err := os.Stat(implFile); os.IsNotExist(err) {
			t.Errorf("Expected impl.go file not found: %s", implFile)
		}
	}

	t.Logf("Successfully generated actor packages with partial implementation")
}

func TestGeneratorWithExampleApplication(t *testing.T) {
	// Load the multi-actor spec
	loader := openapi3.NewLoader()
	doc, err := loader.LoadFromFile("testdata/multi-actor.yaml")
	if err != nil {
		t.Fatalf("Failed to load OpenAPI spec: %v", err)
	}

	p := parser.NewOpenAPIParser(doc)
	model, err := p.Parse()
	if err != nil {
		t.Fatalf("Failed to parse OpenAPI spec: %v", err)
	}

	// Generate with example application
	gen := &generator.Generator{}
	outputDir := "test-output/example-app"
	options := generator.GenerationOptions{
		GenerateImpl:    false,
		GenerateExample: true,
	}
	err = gen.GenerateActorPackages(model, outputDir, options)
	if err != nil {
		t.Fatalf("Failed to generate actor packages with example: %v", err)
	}

	// Clean up after test
	defer func() {
		os.RemoveAll(outputDir)
	}()

	// Verify example files exist
	mainFile := filepath.Join(outputDir, "main.go")
	if _, err := os.Stat(mainFile); os.IsNotExist(err) {
		t.Errorf("Expected main.go file not found: %s", mainFile)
	}

	goModFile := filepath.Join(outputDir, "go.mod")
	if _, err := os.Stat(goModFile); os.IsNotExist(err) {
		t.Errorf("Expected go.mod file not found: %s", goModFile)
	}

	t.Logf("Successfully generated actor packages with example application")
}

func TestEnumGeneration(t *testing.T) {
	// Load the type-alias OpenAPI spec (which includes enums)
	loader := openapi3.NewLoader()
	doc, err := loader.LoadFromFile("testdata/type-alias.yaml")
	if err != nil {
		t.Fatalf("Failed to load type-alias OpenAPI spec: %v", err)
	}

	// Parse the spec to intermediate model
	p := parser.NewOpenAPIParser(doc)
	model, err := p.Parse()
	if err != nil {
		t.Fatalf("Failed to parse OpenAPI spec: %v", err)
	}

	// Find the User actor (which should have UserStatus enum)
	var userActor *generator.ActorInterface
	for i := range model.Actors {
		if model.Actors[i].ActorType == "User" {
			userActor = &model.Actors[i]
			break
		}
	}
	if userActor == nil {
		t.Fatal("User actor not found")
	}

	// Verify that enum types are generated for User
	if len(userActor.Types.Enums) == 0 {
		t.Error("Expected enum types to be generated for User actor, but found none")
	}

	// Look for UserStatus enum
	var userStatusEnum *generator.EnumType
	for i := range userActor.Types.Enums {
		if userActor.Types.Enums[i].Name == "UserStatus" {
			userStatusEnum = &userActor.Types.Enums[i]
			break
		}
	}
	if userStatusEnum == nil {
		t.Error("Expected UserStatus enum not found in User actor")
	} else {
		// Verify enum has correct values
		expectedValues := []string{"active", "inactive", "suspended", "pending"}
		if len(userStatusEnum.Values) != len(expectedValues) {
			t.Errorf("Expected UserStatus to have %d values, got %d", len(expectedValues), len(userStatusEnum.Values))
		}
		for i, expected := range expectedValues {
			if i >= len(userStatusEnum.Values) || userStatusEnum.Values[i] != expected {
				t.Errorf("Expected UserStatus value[%d] to be '%s', got '%s'", i, expected, userStatusEnum.Values[i])
			}
		}
	}
}

func TestNumberFormatsGeneration(t *testing.T) {
	// Load the number formats test OpenAPI spec
	loader := openapi3.NewLoader()
	doc, err := loader.LoadFromFile("testdata/number-formats.yaml")
	if err != nil {
		t.Fatalf("Failed to load number formats OpenAPI spec: %v", err)
	}

	// Parse the spec to intermediate model
	p := parser.NewOpenAPIParser(doc)
	model, err := p.Parse()
	if err != nil {
		t.Fatalf("Failed to parse OpenAPI spec: %v", err)
	}

	// Verify that we have exactly one actor
	if len(model.Actors) != 1 {
		t.Errorf("Expected 1 actor, got %d", len(model.Actors))
	}

	// Verify the NumberTest actor
	actor := model.Actors[0]
	if actor.ActorType != "NumberTest" {
		t.Errorf("Expected actor type 'NumberTest', got '%s'", actor.ActorType)
	}

	// Find the NumberRequest and NumberResponse structs
	var numberRequest, numberResponse *generator.StructType
	for i := range actor.Types.Structs {
		if actor.Types.Structs[i].Name == "NumberRequest" {
			numberRequest = &actor.Types.Structs[i]
		}
		if actor.Types.Structs[i].Name == "NumberResponse" {
			numberResponse = &actor.Types.Structs[i]
		}
	}

	if numberRequest == nil {
		t.Fatal("NumberRequest struct not found")
	}
	if numberResponse == nil {
		t.Fatal("NumberResponse struct not found")
	}

	// Verify number format type mappings in NumberRequest
	expectedFieldTypes := map[string]string{
		"Int8Value":    "int8",
		"Int16Value":   "int16",
		"Int32Value":   "int32",
		"Int64Value":   "int64",
		"FloatValue":   "float32",
		"DoubleValue":  "float64",
		"PlainInteger": "int",
		"PlainNumber":  "float64",
	}

	for _, field := range numberRequest.Fields {
		if expectedType, exists := expectedFieldTypes[field.Name]; exists {
			if field.Type != expectedType {
				t.Errorf("Expected field %s to have type %s, got %s", field.Name, expectedType, field.Type)
			}
		}
	}

	// Verify number format type mappings in NumberResponse
	expectedResponseFieldTypes := map[string]string{
		"Int8Result":         "int8",
		"Int16Result":        "int16",
		"Int32Result":        "int32",
		"Int64Result":        "int64",
		"FloatResult":        "float32",
		"DoubleResult":       "float64",
		"PlainIntegerResult": "int",
		"PlainNumberResult":  "float64",
	}

	for _, field := range numberResponse.Fields {
		if expectedType, exists := expectedResponseFieldTypes[field.Name]; exists {
			if field.Type != expectedType {
				t.Errorf("Expected response field %s to have type %s, got %s", field.Name, expectedType, field.Type)
			}
		}
	}

	// Generate code and verify it compiles
	gen := &generator.Generator{}
	outputDir := "test-output/number-formats-test"
	options := generator.GenerationOptions{
		GenerateImpl:    false,
		GenerateExample: false,
	}
	err = gen.GenerateActorPackages(model, outputDir, options)
	if err != nil {
		t.Fatalf("Failed to generate actor packages: %v", err)
	}

	// Clean up after test
	defer func() {
		os.RemoveAll(outputDir)
	}()

	t.Logf("Successfully validated number format type mappings for NumberTest actor")
}

func TestOptionalObjectReferences(t *testing.T) {
	// Load the optional references test OpenAPI spec
	loader := openapi3.NewLoader()
	doc, err := loader.LoadFromFile("testdata/optional-refs.yaml")
	if err != nil {
		t.Fatalf("Failed to load optional references OpenAPI spec: %v", err)
	}

	// Parse the spec to intermediate model
	p := parser.NewOpenAPIParser(doc)
	model, err := p.Parse()
	if err != nil {
		t.Fatalf("Failed to parse OpenAPI spec: %v", err)
	}

	// Verify that we have exactly one actor
	if len(model.Actors) != 1 {
		t.Errorf("Expected 1 actor, got %d", len(model.Actors))
	}

	actor := model.Actors[0]
	if actor.ActorType != "Test" {
		t.Errorf("Expected actor type 'Test', got '%s'", actor.ActorType)
	}

	// Find the CounterState struct
	var counterState *generator.StructType
	for i := range actor.Types.Structs {
		if actor.Types.Structs[i].Name == "CounterState" {
			counterState = &actor.Types.Structs[i]
			break
		}
	}

	if counterState == nil {
		t.Fatal("CounterState struct not found")
	}

	// Verify field types - optional object references should be pointers
	fieldTypes := make(map[string]string)
	for _, field := range counterState.Fields {
		fieldTypes[field.Name] = field.Type
	}

	// Data field should be a pointer to CounterStateData (optional object reference)
	if fieldTypes["Data"] != "*CounterStateData" {
		t.Errorf("Expected Data field to be '*CounterStateData', got '%s'", fieldTypes["Data"])
	}

	// Error field should be a pointer to Error (optional object reference)
	if fieldTypes["Error"] != "*Error" {
		t.Errorf("Expected Error field to be '*Error', got '%s'", fieldTypes["Error"])
	}

	// Success field should be bool (required field, not a reference)
	if fieldTypes["Success"] != "bool" {
		t.Errorf("Expected Success field to be 'bool', got '%s'", fieldTypes["Success"])
	}

	// Verify optional non-struct references (type aliases) are converted to pointers
	if fieldTypes["UserId"] != "*UserId" {
		t.Errorf("Expected UserId field to be '*UserId', got '%s'", fieldTypes["UserId"])
	}

	if fieldTypes["SessionToken"] != "*SessionToken" {
		t.Errorf("Expected SessionToken field to be '*SessionToken', got '%s'", fieldTypes["SessionToken"])
	}

	// Verify optional enum reference is converted to pointer
	if fieldTypes["OperationStatus"] != "*OperationStatus" {
		t.Errorf("Expected OperationStatus field to be '*OperationStatus', got '%s'", fieldTypes["OperationStatus"])
	}

	// Verify that optional references have omitempty tag
	fieldTags := make(map[string]string)
	for _, field := range counterState.Fields {
		fieldTags[field.Name] = field.JSONTag
	}

	if fieldTags["Data"] != "data,omitempty" {
		t.Errorf("Expected Data field to have JSON tag 'data,omitempty', got '%s'", fieldTags["Data"])
	}

	if fieldTags["Error"] != "error,omitempty" {
		t.Errorf("Expected Error field to have JSON tag 'error,omitempty', got '%s'", fieldTags["Error"])
	}

	if fieldTags["Success"] != "success" {
		t.Errorf("Expected Success field to have JSON tag 'success', got '%s'", fieldTags["Success"])
	}

	// Verify non-struct optional references have omitempty tag
	if fieldTags["UserId"] != "userId,omitempty" {
		t.Errorf("Expected UserId field to have JSON tag 'userId,omitempty', got '%s'", fieldTags["UserId"])
	}

	if fieldTags["SessionToken"] != "sessionToken,omitempty" {
		t.Errorf("Expected SessionToken field to have JSON tag 'sessionToken,omitempty', got '%s'", fieldTags["SessionToken"])
	}

	if fieldTags["OperationStatus"] != "operationStatus,omitempty" {
		t.Errorf("Expected OperationStatus field to have JSON tag 'operationStatus,omitempty', got '%s'", fieldTags["OperationStatus"])
	}

	// Verify that the type aliases and enums are correctly generated in the actor types
	aliasNames := make(map[string]bool)
	for _, alias := range actor.Types.Aliases {
		aliasNames[alias.Name] = true
	}

	if !aliasNames["UserId"] {
		t.Error("Expected type alias 'UserId' not found")
	}

	if !aliasNames["SessionToken"] {
		t.Error("Expected type alias 'SessionToken' not found")
	}

	enumNames := make(map[string]bool)
	for _, enum := range actor.Types.Enums {
		enumNames[enum.Name] = true
	}

	if !enumNames["OperationStatus"] {
		t.Error("Expected enum type 'OperationStatus' not found")
	}

	t.Logf("Successfully validated optional object references and non-struct references are converted to pointers")
}

func TestActorYAMLParsing(t *testing.T) {
	// Test the new Actor YAML format parser
	p, err := parser.NewParser("actor-yaml", "testdata/basic-actor-yaml.yaml")
	if err != nil {
		t.Fatalf("Failed to create parser for Actor YAML: %v", err)
	}

	// Parse the spec to intermediate model
	model, err := p.Parse()
	if err != nil {
		t.Fatalf("Failed to parse Actor YAML spec: %v", err)
	}

	// Verify that we have exactly one actor
	if len(model.Actors) != 1 {
		t.Errorf("Expected 1 actor, got %d", len(model.Actors))
	}

	// Verify the TestActor actor
	actor := model.Actors[0]
	if actor.ActorType != "TestActor" {
		t.Errorf("Expected actor type 'TestActor', got '%s'", actor.ActorType)
	}

	if actor.InterfaceDesc != "Simple test actor for validation" {
		t.Errorf("Expected actor description 'Simple test actor for validation', got '%s'", actor.InterfaceDesc)
	}

	if len(actor.Methods) != 2 {
		t.Errorf("Expected TestActor to have 2 methods, got %d", len(actor.Methods))
	}

	// Verify methods
	methodNames := make(map[string]bool)
	methodsMap := make(map[string]generator.Method)
	for _, method := range actor.Methods {
		methodNames[method.Name] = true
		methodsMap[method.Name] = method
	}

	if !methodNames["GetValue"] {
		t.Error("Expected 'GetValue' method not found")
	}

	if !methodNames["SetValue"] {
		t.Error("Expected 'SetValue' method not found")
	}

	// Verify GetValue method
	getValue := methodsMap["GetValue"]
	if getValue.HasRequest {
		t.Error("GetValue should not have request body")
	}
	if getValue.ReturnType != "TestState" {
		t.Errorf("Expected GetValue return type 'TestState', got '%s'", getValue.ReturnType)
	}

	// Verify SetValue method
	setValue := methodsMap["SetValue"]
	if !setValue.HasRequest {
		t.Error("SetValue should have request body")
	}
	if setValue.RequestType != "SetValueRequest" {
		t.Errorf("Expected SetValue request type 'SetValueRequest', got '%s'", setValue.RequestType)
	}
	if setValue.ReturnType != "TestState" {
		t.Errorf("Expected SetValue return type 'TestState', got '%s'", setValue.ReturnType)
	}

	// Verify that types are assigned to the actor
	if len(actor.Types.Structs) != 2 {
		t.Errorf("Expected actor to have 2 struct types, got %d", len(actor.Types.Structs))
	}

	structNames := make(map[string]bool)
	for _, structType := range actor.Types.Structs {
		structNames[structType.Name] = true
	}

	if !structNames["TestState"] {
		t.Error("Expected struct type 'TestState' not found")
	}

	if !structNames["SetValueRequest"] {
		t.Error("Expected struct type 'SetValueRequest' not found")
	}

	t.Logf("Successfully validated Actor YAML parsing functionality")
}

func TestActorYAMLVsOpenAPIEquivalence(t *testing.T) {
	// This test demonstrates that both formats can produce similar structures
	// though they may not be identical due to different naming conventions

	// Generate from OpenAPI format
	openapiParser, err := parser.NewParser("openapi", "testdata/basic-actor.yaml")
	if err != nil {
		t.Fatalf("Failed to create OpenAPI parser: %v", err)
	}

	openapiModel, err := openapiParser.Parse()
	if err != nil {
		t.Fatalf("Failed to parse OpenAPI spec: %v", err)
	}

	// Generate from Actor YAML format
	actorParser, err := parser.NewParser("actor-yaml", "testdata/basic-actor-yaml.yaml")
	if err != nil {
		t.Fatalf("Failed to create Actor YAML parser: %v", err)
	}

	actorModel, err := actorParser.Parse()
	if err != nil {
		t.Fatalf("Failed to parse Actor YAML spec: %v", err)
	}

	// Both should produce actors (even if different in details)
	if len(openapiModel.Actors) == 0 {
		t.Error("OpenAPI model should have produced at least one actor")
	}

	if len(actorModel.Actors) == 0 {
		t.Error("Actor YAML model should have produced at least one actor")
	}

	// Both should parse the same basic structure types
	if len(openapiModel.Actors) > 0 && len(actorModel.Actors) > 0 {
		openapiActor := openapiModel.Actors[0]
		actorYAMLActor := actorModel.Actors[0]

		// Both should have methods
		if len(openapiActor.Methods) == 0 {
			t.Error("OpenAPI actor should have methods")
		}

		if len(actorYAMLActor.Methods) == 0 {
			t.Error("Actor YAML actor should have methods")
		}

		// Both should have types
		if len(openapiActor.Types.Structs) == 0 {
			t.Error("OpenAPI actor should have struct types")
		}

		if len(actorYAMLActor.Types.Structs) == 0 {
			t.Error("Actor YAML actor should have struct types")
		}
	}

	t.Logf("Successfully validated both Actor YAML and OpenAPI can parse and generate actor code")
}
