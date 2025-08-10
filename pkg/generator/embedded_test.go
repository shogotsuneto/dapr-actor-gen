package generator

import (
	"strings"
	"testing"
)

func TestEmbeddedTemplates(t *testing.T) {
	// Test that all required templates can be loaded from embedded filesystem
	templateNames := []string{
		"actor_types.tmpl",
		"interface.tmpl",
		"factory.tmpl",
	}

	for _, templateName := range templateNames {
		t.Run(templateName, func(t *testing.T) {
			tmpl, err := getEmbeddedTemplate(templateName)
			if err != nil {
				t.Fatalf("Failed to load embedded template %s: %v", templateName, err)
			}

			if tmpl == nil {
				t.Fatalf("Template %s is nil", templateName)
			}

			// Verify the template has some content by checking it's not empty
			if len(tmpl.Templates()) == 0 {
				t.Fatalf("Template %s appears to be empty", templateName)
			}
		})
	}
}

func TestEmbeddedTemplateContent(t *testing.T) {
	// Test that the embedded templates contain expected content
	tests := []struct {
		templateName    string
		expectedContent string
	}{
		{"interface.tmpl", "actor.ServerContext"},
		{"factory.tmpl", "NewActorFactory"},
		{"actor_types.tmpl", "package test"},
	}

	for _, test := range tests {
		t.Run(test.templateName, func(t *testing.T) {
			tmpl, err := getEmbeddedTemplate(test.templateName)
			if err != nil {
				t.Fatalf("Failed to load template: %v", err)
			}

			// Get the template definition to check content
			templateDef := tmpl.Lookup(test.templateName)
			if templateDef == nil {
				t.Fatalf("Template definition not found for %s", test.templateName)
			}

			// Execute template with dummy data to verify it contains expected content
			var output strings.Builder
			dummyData := map[string]interface{}{
				"PackageName": "test",
				"Actor": map[string]interface{}{
					"ActorType":     "TestActor",
					"InterfaceName": "TestAPI",
					"Methods":       []interface{}{},
				},
				"Types": map[string]interface{}{
					"Structs": []interface{}{},
					"Aliases": []interface{}{},
				},
			}

			err = tmpl.Execute(&output, dummyData)
			if err != nil {
				t.Fatalf("Failed to execute template: %v", err)
			}

			result := output.String()
			if !strings.Contains(result, test.expectedContent) {
				t.Fatalf("Template %s does not contain expected content '%s'. Got:\n%s", test.templateName, test.expectedContent, result)
			}
		})
	}
}
