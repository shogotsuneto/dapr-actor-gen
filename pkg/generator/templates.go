package generator

import (
	"embed"
	"strings"
	"text/template"
)

//go:embed templates/*.tmpl
var templatesFS embed.FS

// getEmbeddedTemplate loads a template from the embedded filesystem
func getEmbeddedTemplate(templateName string) (*template.Template, error) {
	// Create template with helper functions
	tmpl := template.New(templateName).Funcs(template.FuncMap{
		"ToLower":      strings.ToLower,
		"ToPascalCase": toPascalCase,
	})

	return tmpl.ParseFS(templatesFS, "templates/"+templateName)
}

// toPascalCase converts a string to PascalCase
// Examples: "increment" -> "Increment", "account_created" -> "AccountCreated", "AccountCreated" -> "AccountCreated"
func toPascalCase(s string) string {
	if s == "" {
		return s
	}

	// If the string contains underscores, hyphens, or spaces, split on those
	if strings.ContainsAny(s, "_- ") {
		words := strings.FieldsFunc(s, func(c rune) bool {
			return c == '_' || c == '-' || c == ' '
		})

		var result strings.Builder
		for _, word := range words {
			if word == "" {
				continue
			}
			// Capitalize first letter and add rest as lowercase
			result.WriteString(strings.Title(strings.ToLower(word)))
		}
		return result.String()
	}

	// If the string is already likely PascalCase or camelCase (has uppercase letters),
	// just ensure it starts with uppercase
	if strings.ToUpper(s[:1]) == s[:1] {
		// Already starts with uppercase, likely already PascalCase
		return s
	} else {
		// Starts with lowercase, convert first letter to uppercase (camelCase -> PascalCase)
		return strings.ToUpper(s[:1]) + s[1:]
	}
}
