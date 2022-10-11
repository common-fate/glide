package ssov2

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/organizations"
	"github.com/common-fate/granted-approvals/accesshandler/pkg/providers"
	"github.com/common-fate/granted-approvals/accesshandler/pkg/types"
	"go.uber.org/zap"
)

// List filter values for arg
func (p *Provider) Filters(ctx context.Context, filterId string) ([]types.Option, error) {
	switch filterId {
	case "organizationalUnit":
		log := zap.S().With("filterId", filterId)
		log.Info("getting aws organization unit id set options")

		opts := []types.Option{}
		hasMore := true
		var nextToken *string

		roots, err := p.orgClient.ListRoots(ctx, &organizations.ListRootsInput{})
		if err != nil {
			return nil, err
		}

		for hasMore {
			ou, err := p.orgClient.ListOrganizationalUnitsForParent(ctx, &organizations.ListOrganizationalUnitsForParentInput{
				ParentId: roots.Roots[0].Id,
			})
			if err != nil {
				return nil, err
			}

			nextToken = ou.NextToken
			hasMore = nextToken != nil
			for _, orgUnit := range ou.OrganizationalUnits {
				opts = append(opts, types.Option{Label: aws.ToString(orgUnit.Name), Value: aws.ToString(orgUnit.Id)})
			}
		}
		return opts, nil
	}

	return nil, &providers.InvalidFilterIdError{FilterId: filterId}

}
