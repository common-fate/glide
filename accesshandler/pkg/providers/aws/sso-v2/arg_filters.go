package ssov2

import (
	"context"
	"fmt"
	"log"

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

	case "tag":
		return []types.Option{{Label: "abc", Value: "123"}, {Label: "xyz", Value: "456"}}, nil

	}

	return nil, &providers.InvalidFilterIdError{FilterId: filterId}

}

func (p *Provider) FetchArgValuesFromDynamicIds(ctx context.Context, argId string, groupName string, groupValues []string) ([]string, error) {
	switch argId {
	case "accountId":
		switch groupName {
		case "organizationalUnit":

			var values []string

			for _, orgUnit := range groupValues {
				hasMore := true
				var nextToken *string

				if hasMore {
					accounts, err := p.orgClient.ListAccountsForParent(ctx, &organizations.ListAccountsForParentInput{
						ParentId: aws.String(orgUnit),
					})
					fmt.Println("accounts", accounts)
					if err != nil {
						log.Fatal("the err is", err)
						return []string{}, err
					}

					nextToken = accounts.NextToken
					// FIXME: Not sure why gostatic check says this is unused.
					hasMore = nextToken != nil

					for _, account := range accounts.Accounts {
						values = append(values, aws.ToString(account.Id))
					}
				}

			}

			return values, nil

		case "tag":
			return []string{}, nil
		}

		return nil, &providers.InvalidFilterIdError{FilterId: groupName}
	}
	return nil, &providers.InvalidArgumentError{Arg: argId}
}
