package deploymentsvc

import (
	"context"
)

// Deletes am existing provider stack
// Checks that the provider type matches one in our registry.
func (s *Service) Delete(ctx context.Context, stackID string) error {

	return DeleteProviderStack(ctx, stackID)
}
