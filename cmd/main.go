package main

import (
	"flag"
	"log"

	"github.com/shogotsuneto/dapr-actor-gen/pkg/generator"
	"github.com/shogotsuneto/dapr-actor-gen/pkg/parser"
)

func main() {
	var generateImpl = flag.Bool("generate-impl", false, "Generate partial implementation stubs with not-implemented errors")
	var generateExample = flag.Bool("generate-example", false, "Generate example main.go, go.mod and other files for a complete app")
	flag.Parse()

	args := flag.Args()
	if len(args) < 2 {
		log.Fatal("Usage: generator [flags] <openapi-file> <base-output-dir>\n" +
			"Flags:\n" +
			"  -generate-impl    Generate partial implementation stubs with not-implemented errors\n" +
			"  -generate-example Generate example main.go, go.mod and other files for a complete app")
	}

	schemaFile := args[0]
	baseOutputDir := args[1]

	// Create parser based on file format
	p, err := parser.NewParserFromFile(schemaFile)
	if err != nil {
		log.Fatalf("Failed to create parser: %v", err)
	}

	// Parse schema to intermediate model
	model, err := p.Parse()
	if err != nil {
		log.Fatalf("Failed to parse schema: %v", err)
	}

	// Create generation options
	options := generator.GenerationOptions{
		GenerateImpl:    *generateImpl,
		GenerateExample: *generateExample,
	}

	// Generate actor-specific packages using the intermediate model
	gen := &generator.Generator{}
	err = gen.GenerateActorPackages(model, baseOutputDir, options)
	if err != nil {
		log.Fatalf("Failed to generate actor packages: %v", err)
	}
}
