package main

import (
	"log"
	"os"

	"github.com/getkin/kin-openapi/openapi3"
	generator "github.com/shogotsuneto/dapr-actor-gen/pkg"
)

func main() {
	if len(os.Args) < 3 {
		log.Fatal("Usage: generator <openapi-file> <base-output-dir>")
	}

	schemaFile := os.Args[1]
	baseOutputDir := os.Args[2]

	// Load OpenAPI spec
	loader := openapi3.NewLoader()
	doc, err := loader.LoadFromFile(schemaFile)
	if err != nil {
		log.Fatalf("Failed to load OpenAPI spec: %v", err)
	}

	// Parse OpenAPI to intermediate model
	parser := generator.NewOpenAPIParser(doc)
	model, err := parser.Parse()
	if err != nil {
		log.Fatalf("Failed to parse OpenAPI spec: %v", err)
	}

	// Generate actor-specific packages using the intermediate model
	gen := &generator.Generator{}
	err = gen.GenerateActorPackages(model, baseOutputDir)
	if err != nil {
		log.Fatalf("Failed to generate actor packages: %v", err)
	}
}