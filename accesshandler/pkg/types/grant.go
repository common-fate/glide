package types

// NewGrant creates a pending Grant from a validated
// CreateGrant payload.
func NewGrant(vcg ValidCreateGrant) Grant {
	return Grant{
		ID:       vcg.Id,
		Provider: vcg.Provider,
		End:      vcg.End,
		Start:    vcg.Start,
		Status:   PENDING,
		Subject:  vcg.Subject,
		With: Grant_With{
			AdditionalProperties: vcg.With.AdditionalProperties,
		},
	}
}
