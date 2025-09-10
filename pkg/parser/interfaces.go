package parser

import "github.com/shogotsuneto/dapr-actor-gen/pkg/generator"

// Parser is the interface for all schema parsers that convert input formats to GenerationModel
type Parser interface {
	Parse() (*generator.GenerationModel, error)
}