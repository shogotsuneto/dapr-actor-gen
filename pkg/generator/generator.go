package generator

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

// Generator handles code generation from the intermediate model
type Generator struct{}

// GenerateActorPackages generates actor-specific packages from the intermediate model
func (g *Generator) GenerateActorPackages(model *GenerationModel, baseOutputDir string) error {
	if len(model.Actors) == 0 {
		return fmt.Errorf("no actors found in the model")
	}

	// Generate package for each actor type
	for _, actor := range model.Actors {
		// Create actor-specific package name and directory using actorType as is
		packageName := strings.ToLower(actor.ActorType)

		outputDir := filepath.Join(baseOutputDir, packageName)

		// Create output directory
		err := os.MkdirAll(outputDir, 0755)
		if err != nil {
			return fmt.Errorf("failed to create output directory %s: %v", outputDir, err)
		}

		// Get actor-specific types directly from the actor
		actorSpecificTypes := actor.Types

		// Create actor model for this specific actor
		actorModel := ActorModel{
			ActorType:      actor.ActorType,
			PackageName:    packageName,
			Types:          actorSpecificTypes,
			ActorInterface: actor,
		}

		// Generate types for this actor
		err = g.generateActorTypes(&actorModel, outputDir)
		if err != nil {
			return fmt.Errorf("failed to generate types for %s: %v", actor.ActorType, err)
		}

		// Generate interface for this actor
		err = g.generateActorInterface(&actorModel, outputDir)
		if err != nil {
			return fmt.Errorf("failed to generate interface for %s: %v", actor.ActorType, err)
		}

		// Generate factory for this actor
		err = g.generateActorFactory(&actorModel, outputDir)
		if err != nil {
			return fmt.Errorf("failed to generate factory for %s: %v", actor.ActorType, err)
		}

		fmt.Printf("Generated actor package: %s\n", outputDir)
		fmt.Printf("  %s/types.go\n", outputDir)
		fmt.Printf("  %s/api.go\n", outputDir)
		fmt.Printf("  %s/factory.go\n", outputDir)
	}

	return nil
}

func (g *Generator) generateActorTypes(actorModel *ActorModel, outputDir string) error {
	// Load template from embedded filesystem
	tmpl, err := getEmbeddedTemplate("actor_types.tmpl")
	if err != nil {
		return fmt.Errorf("failed to parse actor types template: %v", err)
	}

	// Process types directly from the actor model
	processedTypes := TypeDefinitions{
		Structs: make([]StructType, len(actorModel.Types.Structs)),
		Aliases: make([]TypeAlias, len(actorModel.Types.Aliases)),
	}
	copy(processedTypes.Structs, actorModel.Types.Structs)
	copy(processedTypes.Aliases, actorModel.Types.Aliases)

	// Generate types file
	data := struct {
		PackageName string
		Types       TypeDefinitions
	}{
		PackageName: actorModel.PackageName,
		Types:       processedTypes,
	}

	typesFile, err := os.Create(fmt.Sprintf("%s/types.go", outputDir))
	if err != nil {
		return fmt.Errorf("failed to create types file: %v", err)
	}
	defer typesFile.Close()

	err = tmpl.Execute(typesFile, data)
	if err != nil {
		return fmt.Errorf("failed to execute actor types template: %v", err)
	}

	return nil
}

func (g *Generator) generateActorInterface(actorModel *ActorModel, outputDir string) error {
	// Load template from embedded filesystem
	tmpl, err := getEmbeddedTemplate("interface.tmpl")
	if err != nil {
		return fmt.Errorf("failed to parse interface template: %v", err)
	}

	// Generate interface file for this actor
	data := SingleActorTemplateData{
		PackageName: actorModel.PackageName,
		Actor:       actorModel.ActorInterface,
	}

	// Use api.go as filename instead of generated.go for better clarity
	interfaceFile, err := os.Create(filepath.Join(outputDir, "api.go"))
	if err != nil {
		return fmt.Errorf("failed to create interface file: %v", err)
	}
	defer interfaceFile.Close()

	err = tmpl.Execute(interfaceFile, data)
	if err != nil {
		return fmt.Errorf("failed to execute interface template: %v", err)
	}

	return nil
}

func (g *Generator) generateActorFactory(actorModel *ActorModel, outputDir string) error {
	// Load template from embedded filesystem
	tmpl, err := getEmbeddedTemplate("factory.tmpl")
	if err != nil {
		return fmt.Errorf("failed to parse factory template: %v", err)
	}

	// Generate factory file for this actor
	data := SingleActorTemplateData{
		PackageName: actorModel.PackageName,
		Actor:       actorModel.ActorInterface,
	}

	factoryFile, err := os.Create(filepath.Join(outputDir, "factory.go"))
	if err != nil {
		return fmt.Errorf("failed to create factory file: %v", err)
	}
	defer factoryFile.Close()

	err = tmpl.Execute(factoryFile, data)
	if err != nil {
		return fmt.Errorf("failed to execute factory template: %v", err)
	}

	return nil
}

// Utility functions