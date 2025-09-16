package parser

import (
	"fmt"
	"sort"
	"strings"

	"github.com/shogotsuneto/dapr-actor-gen/pkg/generator"
	"gopkg.in/yaml.v3"
)

// ActorSchemaParser handles conversion from Actor Schema specification to intermediate model
type ActorSchemaParser struct {
	schema *ActorSchemaSchema
}

// NewActorSchemaParser creates a new Actor Schema parser from YAML content
func NewActorSchemaParser(yamlContent []byte) (*ActorSchemaParser, error) {
	var schema ActorSchemaSchema
	if err := yaml.Unmarshal(yamlContent, &schema); err != nil {
		return nil, fmt.Errorf("failed to parse actor YAML: %v", err)
	}

	return &ActorSchemaParser{schema: &schema}, nil
}

// Parse converts the Actor Schema specification to an intermediate generator.GenerationModel
func (p *ActorSchemaParser) Parse() (*generator.GenerationModel, error) {
	model := &generator.GenerationModel{}

	// Parse actors and their methods
	if err := p.parseActors(model); err != nil {
		return nil, fmt.Errorf("failed to parse actors: %v", err)
	}

	// Parse types and assign them to actors that use them
	if err := p.parseAndCategorizeTypes(model); err != nil {
		return nil, fmt.Errorf("failed to parse and categorize types: %v", err)
	}

	return model, nil
}

// parseActors converts actors from YAML to the intermediate model
func (p *ActorSchemaParser) parseActors(model *generator.GenerationModel) error {
	var actors []generator.ActorInterface

	// Convert each actor definition
	for actorType, actorDef := range p.schema.Actors {
		actor := generator.ActorInterface{
			ActorType:     actorType,
			InterfaceName: actorType + "API",
			InterfaceDesc: actorDef.Description,
		}

		// Convert methods
		for methodName, methodDef := range actorDef.Methods {
			method := generator.Method{
				Name:    methodName,
				Comment: methodDef.Description,
			}

			// Set request type if specified
			if methodDef.Request != "" {
				method.HasRequest = true
				method.RequestType = methodDef.Request
			}

			// Set return type if specified
			if methodDef.Returns != "" {
				method.ReturnType = methodDef.Returns
			} else {
				method.ReturnType = "error"
			}

			actor.Methods = append(actor.Methods, method)
		}

		// Sort methods for consistent ordering
		sort.Slice(actor.Methods, func(i, j int) bool {
			return actor.Methods[i].Name < actor.Methods[j].Name
		})

		actors = append(actors, actor)
	}

	// Sort actors for consistent ordering
	sort.Slice(actors, func(i, j int) bool {
		return actors[i].ActorType < actors[j].ActorType
	})

	model.Actors = actors
	return nil
}

// parseAndCategorizeTypes converts types from YAML and assigns them to actors
func (p *ActorSchemaParser) parseAndCategorizeTypes(model *generator.GenerationModel) error {
	// Parse all types from the YAML schema
	allTypes, err := p.parseTypes()
	if err != nil {
		return err
	}

	// Sort all types for consistent ordering
	p.sortTypes(&allTypes)

	// Assign types to actors that use them
	return p.categorizeTypesIntoActors(model, allTypes)
}

// parseTypes converts type definitions from YAML to intermediate model
func (p *ActorSchemaParser) parseTypes() (generator.TypeDefinitions, error) {
	var types generator.TypeDefinitions

	for typeName, typeDef := range p.schema.Types {
		switch typeDef.Type {
		case "object":
			structType, err := p.convertToStruct(typeName, typeDef)
			if err != nil {
				return types, fmt.Errorf("failed to convert type %s: %v", typeName, err)
			}
			types.Structs = append(types.Structs, structType)

		case "string":
			if len(typeDef.Enum) > 0 {
				enumType := p.convertToEnum(typeName, typeDef)
				types.Enums = append(types.Enums, enumType)
			} else {
				aliasType := p.convertToAlias(typeName, typeDef)
				types.Aliases = append(types.Aliases, aliasType)
			}

		case "integer", "number", "boolean":
			if len(typeDef.Enum) > 0 {
				enumType := p.convertToEnum(typeName, typeDef)
				types.Enums = append(types.Enums, enumType)
			} else {
				aliasType := p.convertToAlias(typeName, typeDef)
				types.Aliases = append(types.Aliases, aliasType)
			}

		case "array":
			// Arrays become type aliases to Go slices
			aliasType := p.convertToAlias(typeName, typeDef)
			types.Aliases = append(types.Aliases, aliasType)

		default:
			return types, fmt.Errorf("unsupported type: %s for %s", typeDef.Type, typeName)
		}
	}

	return types, nil
}

// convertToStruct converts a YAML type definition to a struct type
func (p *ActorSchemaParser) convertToStruct(typeName string, typeDef TypeDefinition) (generator.StructType, error) {
	structType := generator.StructType{
		Name:        typeName,
		Description: typeDef.Description,
	}

	// Convert properties
	for propName, propDef := range typeDef.Properties {
		field := generator.Field{
			Name:    strings.Title(propName),
			JSONTag: propName,
			Comment: propDef.Description,
		}

		// Convert property type to Go type
		goType, err := p.convertPropertyToGoType(propDef)
		if err != nil {
			return structType, fmt.Errorf("failed to convert property %s: %v", propName, err)
		}

		// Check if field is required
		isRequired := false
		for _, required := range typeDef.Required {
			if required == propName {
				isRequired = true
				break
			}
		}

		// Make optional fields pointers for objects
		if !isRequired && (propDef.Type == "object" || propDef.Ref != "") {
			if !strings.HasPrefix(goType, "*") {
				goType = "*" + goType
			}
		}

		field.Type = goType
		structType.Fields = append(structType.Fields, field)
	}

	// Sort fields for consistent ordering
	sort.Slice(structType.Fields, func(i, j int) bool {
		return structType.Fields[i].Name < structType.Fields[j].Name
	})

	return structType, nil
}

// convertToEnum converts a YAML type definition to an enum type
func (p *ActorSchemaParser) convertToEnum(typeName string, typeDef TypeDefinition) generator.EnumType {
	enumType := generator.EnumType{
		Name:        typeName,
		Description: typeDef.Description,
		BaseType:    p.getGoTypeForYAMLType(typeDef.Type, typeDef.Format),
	}

	// Convert enum values to strings
	for _, value := range typeDef.Enum {
		enumType.Values = append(enumType.Values, fmt.Sprintf("%v", value))
	}

	return enumType
}

// convertToAlias converts a YAML type definition to a type alias
func (p *ActorSchemaParser) convertToAlias(typeName string, typeDef TypeDefinition) generator.TypeAlias {
	aliasType := generator.TypeAlias{
		Name:        typeName,
		Description: typeDef.Description,
	}

	if typeDef.Type == "array" && typeDef.Items != nil {
		// Array type
		itemType, _ := p.convertPropertyToGoType(PropertyDefinition{
			Type:   typeDef.Items.Type,
			Format: typeDef.Items.Format,
			Ref:    "", // TODO: Handle refs in arrays if needed
		})
		aliasType.AliasTarget = "[]" + itemType
	} else {
		// Simple type alias
		aliasType.AliasTarget = p.getGoTypeForYAMLType(typeDef.Type, typeDef.Format)
	}

	return aliasType
}

// convertPropertyToGoType converts a property definition to Go type string
func (p *ActorSchemaParser) convertPropertyToGoType(propDef PropertyDefinition) (string, error) {
	if propDef.Ref != "" {
		// Reference to another type - extract type name from $ref
		if strings.HasPrefix(propDef.Ref, "#/types/") {
			return strings.TrimPrefix(propDef.Ref, "#/types/"), nil
		}
		return strings.TrimPrefix(propDef.Ref, "#/"), nil
	}

	if propDef.Type == "array" && propDef.Items != nil {
		itemType, err := p.convertPropertyToGoType(*propDef.Items)
		if err != nil {
			return "", err
		}
		return "[]" + itemType, nil
	}

	if propDef.Type == "object" && len(propDef.Properties) > 0 {
		// Inline object - not supported, should use $ref
		return "interface{}", nil
	}

	return p.getGoTypeForYAMLType(propDef.Type, propDef.Format), nil
}

// getGoTypeForYAMLType maps YAML types to Go types
func (p *ActorSchemaParser) getGoTypeForYAMLType(yamlType, format string) string {
	switch yamlType {
	case "string":
		return "string"
	case "integer":
		switch format {
		case "int32":
			return "int32"
		case "int64":
			return "int64"
		default:
			return "int"
		}
	case "number":
		switch format {
		case "float":
			return "float32"
		case "double":
			return "float64"
		default:
			return "float64"
		}
	case "boolean":
		return "bool"
	default:
		return "interface{}"
	}
}

// sortTypes sorts all type collections for consistent ordering
func (p *ActorSchemaParser) sortTypes(types *generator.TypeDefinitions) {
	sort.Slice(types.Structs, func(i, j int) bool {
		return types.Structs[i].Name < types.Structs[j].Name
	})
	sort.Slice(types.Aliases, func(i, j int) bool {
		return types.Aliases[i].Name < types.Aliases[j].Name
	})
	sort.Slice(types.Enums, func(i, j int) bool {
		return types.Enums[i].Name < types.Enums[j].Name
	})
}

// categorizeTypesIntoActors analyzes types and assigns them to actors that use them
func (p *ActorSchemaParser) categorizeTypesIntoActors(model *generator.GenerationModel, allTypes generator.TypeDefinitions) error {
	// Create a map to track which types are used by which actors
	typeUsage := make(map[string]map[string]bool) // type -> actor -> used

	// Initialize usage map for all types
	for _, structType := range allTypes.Structs {
		typeUsage[structType.Name] = make(map[string]bool)
	}
	for _, aliasType := range allTypes.Aliases {
		typeUsage[aliasType.Name] = make(map[string]bool)
	}
	for _, enumType := range allTypes.Enums {
		typeUsage[enumType.Name] = make(map[string]bool)
	}

	// Analyze which actors use which types by examining request/response types
	for _, actor := range model.Actors {
		for _, method := range actor.Methods {
			// Track request types
			if method.HasRequest && method.RequestType != "" {
				p.markTypeAsUsed(typeUsage, method.RequestType, actor.ActorType, allTypes)
			}

			// Track return types (remove pointer prefix if present)
			returnType := strings.TrimPrefix(method.ReturnType, "*")
			if returnType != "error" && returnType != "" {
				p.markTypeAsUsed(typeUsage, returnType, actor.ActorType, allTypes)
			}
		}
	}

	// Assign types to actors based on usage
	for i := range model.Actors {
		actor := &model.Actors[i]

		// Assign structs
		for _, structType := range allTypes.Structs {
			if typeUsage[structType.Name][actor.ActorType] {
				actor.Types.Structs = append(actor.Types.Structs, structType)
			}
		}

		// Assign aliases
		for _, aliasType := range allTypes.Aliases {
			if typeUsage[aliasType.Name][actor.ActorType] {
				actor.Types.Aliases = append(actor.Types.Aliases, aliasType)
			}
		}

		// Assign enums
		for _, enumType := range allTypes.Enums {
			if typeUsage[enumType.Name][actor.ActorType] {
				actor.Types.Enums = append(actor.Types.Enums, enumType)
			}
		}
	}

	return nil
}

// markTypeAsUsed marks a type as used by an actor and recursively marks dependencies
func (p *ActorSchemaParser) markTypeAsUsed(typeUsage map[string]map[string]bool, typeName, actorType string, allTypes generator.TypeDefinitions) {
	if typeUsage[typeName] == nil {
		return // Type doesn't exist
	}

	if typeUsage[typeName][actorType] {
		return // Already marked
	}

	typeUsage[typeName][actorType] = true

	// Find the type definition and mark dependencies
	for _, structType := range allTypes.Structs {
		if structType.Name == typeName {
			for _, field := range structType.Fields {
				// Extract type name from field type (remove pointers, slices, etc.)
				fieldTypeName := p.extractTypeName(field.Type)
				if fieldTypeName != "" && typeUsage[fieldTypeName] != nil {
					p.markTypeAsUsed(typeUsage, fieldTypeName, actorType, allTypes)
				}
			}
		}
	}
}

// extractTypeName extracts the core type name from a Go type string
func (p *ActorSchemaParser) extractTypeName(goType string) string {
	// Remove pointer prefix
	goType = strings.TrimPrefix(goType, "*")

	// Remove slice prefix
	goType = strings.TrimPrefix(goType, "[]")

	// Remove built-in types
	builtinTypes := map[string]bool{
		"string": true, "int": true, "int32": true, "int64": true,
		"float32": true, "float64": true, "bool": true, "interface{}": true,
		"error": true,
	}

	if builtinTypes[goType] {
		return ""
	}

	return goType
}
