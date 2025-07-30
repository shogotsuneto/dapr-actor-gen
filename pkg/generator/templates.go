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
		"ToLower": strings.ToLower,
	})
	
	return tmpl.ParseFS(templatesFS, "templates/"+templateName)
}