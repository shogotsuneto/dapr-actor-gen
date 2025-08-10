package generator

import (
	"embed"
	"strings"
	"text/template"
	"unicode"
)

//go:embed templates/*.tmpl
var templatesFS embed.FS

// capitalizeFirst capitalizes the first letter of a string for enum constants
func capitalizeFirst(s string) string {
	if len(s) == 0 {
		return s
	}
	runes := []rune(s)
	runes[0] = unicode.ToUpper(runes[0])
	return string(runes)
}

// getEmbeddedTemplate loads a template from the embedded filesystem
func getEmbeddedTemplate(templateName string) (*template.Template, error) {
	// Create template with helper functions
	tmpl := template.New(templateName).Funcs(template.FuncMap{
		"ToLower": strings.ToLower,
		"title":   capitalizeFirst,
	})

	return tmpl.ParseFS(templatesFS, "templates/"+templateName)
}
