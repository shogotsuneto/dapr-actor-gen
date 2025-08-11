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
	// Load the multi-actors OpenAPI spec (which includes enums)
	loader := openapi3.NewLoader()
	doc, err := loader.LoadFromFile("../../examples/multi-actors/openapi.yaml")
	if err != nil {
		t.Fatalf("Failed to load multi-actors OpenAPI spec: %v", err)
	}

	// Parse the spec to intermediate model
	p := parser.NewOpenAPIParser(doc)
	model, err := p.Parse()
	if err != nil {
		t.Fatalf("Failed to parse OpenAPI spec: %v", err)
	}

	// Find the BankAccount actor (which should have AccountEventEventType enum)
	var bankActor *generator.ActorInterface
	for i := range model.Actors {
		if model.Actors[i].ActorType == "BankAccount" {
			bankActor = &model.Actors[i]
			break
		}
	}
	if bankActor == nil {
		t.Fatal("BankAccount actor not found")
	}

	// Verify that enum types are generated for BankAccount
	if len(bankActor.Types.Enums) == 0 {
		t.Error("Expected enum types to be generated for BankAccount actor, but found none")
	}

	// Look for AccountEventEventType enum
	var eventTypeEnum *generator.EnumType
	for i := range bankActor.Types.Enums {
		if bankActor.Types.Enums[i].Name == "AccountEventEventType" {
			eventTypeEnum = &bankActor.Types.Enums[i]
			break
		}
	}
	if eventTypeEnum == nil {
		t.Error("Expected AccountEventEventType enum not found in BankAccount actor")
	} else {
		// Verify enum has correct values
		expectedValues := []string{"AccountCreated", "MoneyDeposited", "MoneyWithdrawn"}
		if len(eventTypeEnum.Values) != len(expectedValues) {
			t.Errorf("Expected AccountEventEventType to have %d values, got %d", len(expectedValues), len(eventTypeEnum.Values))
		}
		for i, expected := range expectedValues {
			if i >= len(eventTypeEnum.Values) || eventTypeEnum.Values[i] != expected {
				t.Errorf("Expected AccountEventEventType value[%d] to be '%s', got '%s'", i, expected, eventTypeEnum.Values[i])
			}
		}
	}

	// Find the Counter actor (which should have CounterStatus and CounterOperation enums)
	var counterActor *generator.ActorInterface
	for i := range model.Actors {
		if model.Actors[i].ActorType == "Counter" {
			counterActor = &model.Actors[i]
			break
		}
	}
	if counterActor == nil {
		t.Fatal("Counter actor not found")
	}

	// Verify that enum types are generated for Counter
	if len(counterActor.Types.Enums) == 0 {
		t.Error("Expected enum types to be generated for Counter actor, but found none")
	}

	// Look for CounterStatus and CounterOperation enums
	enumNames := make(map[string]bool)
	for _, enum := range counterActor.Types.Enums {
		enumNames[enum.Name] = true
	}

	expectedEnums := []string{"CounterStatus", "CounterOperation"}
	for _, expected := range expectedEnums {
		if !enumNames[expected] {
			t.Errorf("Expected enum type '%s' not found in Counter actor", expected)
		}
	}

	// Verify CounterStatus enum values
	for _, enum := range counterActor.Types.Enums {
		if enum.Name == "CounterStatus" {
			expectedValues := []string{"active", "paused", "error", "reset"}
			if len(enum.Values) != len(expectedValues) {
				t.Errorf("Expected CounterStatus to have %d values, got %d", len(expectedValues), len(enum.Values))
			}
			for i, expected := range expectedValues {
				if i >= len(enum.Values) || enum.Values[i] != expected {
					t.Errorf("Expected CounterStatus value[%d] to be '%s', got '%s'", i, expected, enum.Values[i])
				}
			}
		}
	}
}
