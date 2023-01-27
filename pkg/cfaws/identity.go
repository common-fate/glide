package cfaws

import (
	"context"

	// cloudproof is a repo from commonfate which contains helpers to validate aws identity
	"github.com/common-fate/cloudproof/aws"
)

func VerifyCallerIdentity(ctx context.Context, identityProof aws.IdentityProof) (*aws.Identity, error) {
	return identityProof.Verify(ctx)
}

func CreateIdentityProof(ctx context.Context) (*aws.IdentityProof, error) {
	return aws.NewIdentityProof(ctx)
}
