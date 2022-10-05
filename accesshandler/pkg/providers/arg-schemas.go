package providers

import "github.com/invopop/jsonschema"

// GetArgumentTitleFromSchema will lookup an argument in a provider arg schema and return the title if it exists else and empty string
func GetArgumentTitleFromSchema(s *jsonschema.Schema, argumentID string) string {
	for _, def := range s.Definitions {
		if a, ok := def.Properties.Get(argumentID); ok {
			if as, ok := a.(*jsonschema.Schema); ok {
				return as.Title
			}
		}
	}
	return ""
}

// GetArgumentDescriptionFromSchema will lookup an argument in a provider arg schema and return the description if it exists else an empty string
func GetArgumentDescriptionFromSchema(s *jsonschema.Schema, argumentID string) string {
	for _, def := range s.Definitions {
		if a, ok := def.Properties.Get(argumentID); ok {
			if as, ok := a.(*jsonschema.Schema); ok {
				return as.Description
			}
		}
	}
	return ""
}
