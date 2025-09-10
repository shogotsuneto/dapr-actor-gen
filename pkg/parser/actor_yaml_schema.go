package parser

// ActorYAMLSchema represents the root structure of an actor definition YAML file
type ActorYAMLSchema struct {
	Actors map[string]ActorDefinition `yaml:"actors"`
	Types  map[string]TypeDefinition  `yaml:"types"`
}

// ActorDefinition represents a single actor's definition
type ActorDefinition struct {
	Description string                    `yaml:"description,omitempty"`
	Methods     map[string]MethodDefinition `yaml:"methods"`
}

// MethodDefinition represents a single actor method's definition
type MethodDefinition struct {
	Description string `yaml:"description,omitempty"`
	Request     string `yaml:"request,omitempty"`     // Reference to type name for request
	Returns     string `yaml:"returns,omitempty"`     // Reference to type name for response
}

// TypeDefinition represents a type definition that mirrors OpenAPI schema structure
type TypeDefinition struct {
	Type                 string                     `yaml:"type"`
	Description          string                     `yaml:"description,omitempty"`
	Properties           map[string]PropertyDefinition `yaml:"properties,omitempty"`
	Required             []string                   `yaml:"required,omitempty"`
	Enum                 []interface{}              `yaml:"enum,omitempty"`
	Items                *TypeDefinition            `yaml:"items,omitempty"`     // For arrays
	AdditionalProperties *bool                      `yaml:"additionalProperties,omitempty"`
	Format               string                     `yaml:"format,omitempty"`    // e.g., "date-time", "int32", "double"
	Minimum              *float64                   `yaml:"minimum,omitempty"`
	Maximum              *float64                   `yaml:"maximum,omitempty"`
	MinLength            *int                       `yaml:"minLength,omitempty"`
	MaxLength            *int                       `yaml:"maxLength,omitempty"`
	Pattern              string                     `yaml:"pattern,omitempty"`
	Example              interface{}                `yaml:"example,omitempty"`
}

// PropertyDefinition represents a property within a type definition
type PropertyDefinition struct {
	Type        string                     `yaml:"type"`
	Description string                     `yaml:"description,omitempty"`
	Properties  map[string]PropertyDefinition `yaml:"properties,omitempty"` // For nested objects
	Items       *PropertyDefinition        `yaml:"items,omitempty"`       // For arrays
	Enum        []interface{}              `yaml:"enum,omitempty"`
	Format      string                     `yaml:"format,omitempty"`
	Minimum     *float64                   `yaml:"minimum,omitempty"`
	Maximum     *float64                   `yaml:"maximum,omitempty"`
	MinLength   *int                       `yaml:"minLength,omitempty"`
	MaxLength   *int                       `yaml:"maxLength,omitempty"`
	Pattern     string                     `yaml:"pattern,omitempty"`
	Example     interface{}                `yaml:"example,omitempty"`
	Ref         string                     `yaml:"$ref,omitempty"`        // Reference to another type
}