package ssov2

import (
	"context"
	"fmt"
	"log"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/organizations"
	"github.com/common-fate/granted-approvals/accesshandler/pkg/providers"
)

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
