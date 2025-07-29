package generator

import (
	"embed"
	"text/template"
)

//go:embed templates/*.tmpl
var templatesFS embed.FS

// getEmbeddedTemplate loads a template from the embedded filesystem
func getEmbeddedTemplate(templateName string) (*template.Template, error) {
	return template.ParseFS(templatesFS, "templates/"+templateName)
}