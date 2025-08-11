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
func (g *Generator) GenerateActorPackages(model *GenerationModel, baseOutputDir string, options GenerationOptions) error {
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

		// Optionally generate partial implementation
		if options.GenerateImpl {
			err = g.generatePartialImplementation(&actorModel, outputDir)
			if err != nil {
				return fmt.Errorf("failed to generate partial implementation for %s: %v", actor.ActorType, err)
			}
		}

		fmt.Printf("Generated actor package: %s\n", outputDir)
		fmt.Printf("  %s/types.go\n", outputDir)
		fmt.Printf("  %s/api.go\n", outputDir)
		fmt.Printf("  %s/factory.go\n", outputDir)
		if options.GenerateImpl {
			fmt.Printf("  %s/impl.go\n", outputDir)
		}
	}

	// Optionally generate example application
	if options.GenerateExample {
		err := g.generateExampleApplication(model, baseOutputDir)
		if err != nil {
			return fmt.Errorf("failed to generate example application: %v", err)
		}
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
		Enums:   make([]EnumType, len(actorModel.Types.Enums)),
	}
	copy(processedTypes.Structs, actorModel.Types.Structs)
	copy(processedTypes.Aliases, actorModel.Types.Aliases)
	copy(processedTypes.Enums, actorModel.Types.Enums)

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

func (g *Generator) generatePartialImplementation(actorModel *ActorModel, outputDir string) error {
	// Load template from embedded filesystem
	tmpl, err := getEmbeddedTemplate("partial_impl.tmpl")
	if err != nil {
		return fmt.Errorf("failed to parse partial implementation template: %v", err)
	}

	// Generate implementation file for this actor
	data := SingleActorTemplateData{
		PackageName: actorModel.PackageName,
		Actor:       actorModel.ActorInterface,
	}

	implFile, err := os.Create(filepath.Join(outputDir, "impl.go"))
	if err != nil {
		return fmt.Errorf("failed to create implementation file: %v", err)
	}
	defer implFile.Close()

	err = tmpl.Execute(implFile, data)
	if err != nil {
		return fmt.Errorf("failed to execute partial implementation template: %v", err)
	}

	return nil
}

func (g *Generator) generateExampleApplication(model *GenerationModel, baseOutputDir string) error {
	// Generate main.go
	err := g.generateExampleMain(model, baseOutputDir)
	if err != nil {
		return fmt.Errorf("failed to generate example main.go: %v", err)
	}

	// Generate go.mod
	err = g.generateExampleGoMod(model, baseOutputDir)
	if err != nil {
		return fmt.Errorf("failed to generate example go.mod: %v", err)
	}

	fmt.Printf("Generated example application files:\n")
	fmt.Printf("  %s/main.go\n", baseOutputDir)
	fmt.Printf("  %s/go.mod\n", baseOutputDir)

	return nil
}

func (g *Generator) generateExampleMain(model *GenerationModel, baseOutputDir string) error {
	// Load template from embedded filesystem
	tmpl, err := getEmbeddedTemplate("example_main.tmpl")
	if err != nil {
		return fmt.Errorf("failed to parse example main template: %v", err)
	}

	// Generate main.go file
	data := struct {
		Actors     []ActorInterface
		ModuleName string
	}{
		Actors:     model.Actors,
		ModuleName: "example-dapr-actors",
	}

	mainFile, err := os.Create(filepath.Join(baseOutputDir, "main.go"))
	if err != nil {
		return fmt.Errorf("failed to create main.go file: %v", err)
	}
	defer mainFile.Close()

	err = tmpl.Execute(mainFile, data)
	if err != nil {
		return fmt.Errorf("failed to execute example main template: %v", err)
	}

	return nil
}

func (g *Generator) generateExampleGoMod(model *GenerationModel, baseOutputDir string) error {
	// Load template from embedded filesystem
	tmpl, err := getEmbeddedTemplate("example_gomod.tmpl")
	if err != nil {
		return fmt.Errorf("failed to parse example go.mod template: %v", err)
	}

	// Generate go.mod file
	data := struct {
		ModuleName string
	}{
		ModuleName: "example-dapr-actors",
	}

	goModFile, err := os.Create(filepath.Join(baseOutputDir, "go.mod"))
	if err != nil {
		return fmt.Errorf("failed to create go.mod file: %v", err)
	}
	defer goModFile.Close()

	err = tmpl.Execute(goModFile, data)
	if err != nil {
		return fmt.Errorf("failed to execute example go.mod template: %v", err)
	}

	return nil
}

// Utility functions
