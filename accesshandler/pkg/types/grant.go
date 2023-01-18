package types

// NewGrant creates a pending Grant from a validated
// CreateGrant payload.
func NewGrant(vcg ValidCreateGrant) Grant {
	return Grant{
		ID:       vcg.Id,
		Provider: vcg.Provider,
		End:      vcg.End,
		Start:    vcg.Start,
		Status:   GrantStatusPENDING,
		Subject:  vcg.Subject,
		With:     vcg.With,
	}

}
