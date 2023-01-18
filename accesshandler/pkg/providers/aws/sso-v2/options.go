package ssov2

import (
	"context"
	"errors"
	"sync"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/organizations"
	organizationTypes "github.com/aws/aws-sdk-go-v2/service/organizations/types"
	"github.com/aws/aws-sdk-go-v2/service/ssoadmin"
	ssoadmintypes "github.com/aws/aws-sdk-go-v2/service/ssoadmin/types"

	"github.com/common-fate/common-fate/accesshandler/pkg/providers"
	"github.com/common-fate/common-fate/accesshandler/pkg/types"
	"go.uber.org/zap"
	"golang.org/x/sync/errgroup"
)

// List options for arg
func (p *Provider) Options(ctx context.Context, arg string) (*types.ArgOptionsResponse, error) {
	switch arg {
	case "permissionSetArn":
		log := zap.S().With("arg", arg)
		log.Info("getting sso permission set options")

		var opts types.ArgOptionsResponse
		// prevent concurrent writes to `opts` in goroutines
		var mu sync.Mutex

		g, gctx := errgroup.WithContext(ctx)
		g.SetLimit(5) // set a limit here to avoid hitting API rate limits in cases where accounts have many permission sets

		hasMore := true
		var nextToken *string

		for hasMore {
			o, err := p.client.ListPermissionSets(ctx, &ssoadmin.ListPermissionSetsInput{
				InstanceArn: aws.String(p.instanceARN.Get()),
				NextToken:   nextToken,
			})
			if err != nil {
				// ensure we don't have stale goroutines hanging around - just send the error into the errgroup
				// and then call Wait() to wrap up goroutines.
				g.Go(func() error { return err })
				_ = g.Wait()
				return nil, err
			}

			for _, ARN := range o.PermissionSets {
				ARNCopy := ARN

				g.Go(func() error {
					po, err := p.client.DescribePermissionSet(gctx, &ssoadmin.DescribePermissionSetInput{
						InstanceArn: aws.String(p.instanceARN.Get()), PermissionSetArn: aws.String(ARNCopy),
					})

					var ade *ssoadmintypes.AccessDeniedException
					if errors.As(err, &ade) {
						// we don't have access to this permission set, so don't include it in the options.
						log.Debug("access denied when attempting to describe permission set, not including in options", "permissionset.arn", ARNCopy, zap.Error(err))
						return nil
					}

					if err != nil {
						return err
					}

					mu.Lock()
					defer mu.Unlock()
					opts.Options = append(opts.Options, types.Option{Label: aws.ToString(po.PermissionSet.Name), Value: ARNCopy, Description: po.PermissionSet.Description})

					return nil
				})
			}

			nextToken = o.NextToken
			hasMore = nextToken != nil
		}

		err := g.Wait()
		if err != nil {
			return nil, err
		}

		return &opts, nil
	case "accountId":
		log := zap.S().With("arg", arg)
		log.Info("getting sso permission set options")
		var opts types.ArgOptionsResponse

		graph, err := p.buildOrganizationGraph(ctx)
		if err != nil {
			return nil, err
		}

		for _, acct := range graph.Root.DescendantAccounts() {
			opts.Options = append(opts.Options, types.Option{Label: aws.ToString(acct.Account.Name), Value: aws.ToString(acct.Account.Id)})
		}
		log.Info("getting aws organization unit id set options")

		ous, err := graph.generateOuGroupOptions(ctx)
		if err != nil {
			return nil, err
		}

		// tags, err := p.generateTagGroupOptionsForAccounts(ctx, graph.Root.DescendantOrganisationTypeAccounts())
		// if err != nil {
		// 	return nil, err
		// }
		// tagGroup := types.Group{
		// 	Title:   "Tags",
		// 	Id:      "tag",
		// 	Options: tags,
		// }

		opts.Groups = &types.Groups{"organizationalUnit": ous}

		return &opts, nil
	}

	return nil, &providers.InvalidArgumentError{Arg: arg}

}

func (g *OrganizationGraph) generateOuGroupOptions(ctx context.Context) ([]types.GroupOption, error) {
	// first add the organization root
	groupOptions := []types.GroupOption{
		{Label: aws.ToString(g.Root.Root.Name), Value: aws.ToString(g.Root.Root.Id), Children: g.Root.DescendantAccountIDs()},
	}
	for _, orgUnit := range g.Root.DescendantOrganizationalUnits() {
		labelPrefix := orgUnit.Ancestors.Path() + "/"
		option := types.GroupOption{Label: aws.ToString(orgUnit.OrganizationalUnit.Name), Value: aws.ToString(orgUnit.OrganizationalUnit.Id), Children: orgUnit.DescendantAccountIDs(), LabelPrefix: &labelPrefix}
		groupOptions = append(groupOptions, option)
	}
	return groupOptions, nil
}

func (p *Provider) listChildOusForParent(ctx context.Context, parentID string) (ous []organizationTypes.OrganizationalUnit, err error) {
	hasMore := true
	var nextToken *string
	for hasMore {
		ou, err := p.orgClient.ListOrganizationalUnitsForParent(ctx, &organizations.ListOrganizationalUnitsForParentInput{
			ParentId:  aws.String(parentID),
			NextToken: nextToken,
		})
		if err != nil {
			return nil, err
		}
		nextToken = ou.NextToken
		hasMore = nextToken != nil
		ous = append(ous, ou.OrganizationalUnits...)
	}
	return
}

func (p *Provider) listChildAccountsForParent(ctx context.Context, parentID string) (accounts []organizationTypes.Account, err error) {
	hasMore := true
	var nextToken *string
	for hasMore {
		ou, err := p.orgClient.ListAccountsForParent(ctx, &organizations.ListAccountsForParentInput{
			ParentId:  aws.String(parentID),
			NextToken: nextToken,
		})
		if err != nil {
			return nil, err
		}
		nextToken = ou.NextToken
		hasMore = nextToken != nil
		accounts = append(accounts, ou.Accounts...)
	}
	return
}
func (p *Provider) ArgOptionGroupValues(ctx context.Context, argId string, groupID string, groupValues []string) ([]string, error) {
	graph, err := p.buildOrganizationGraph(ctx)
	if err != nil {
		return nil, err
	}
	switch argId {
	case "accountId":
		switch groupID {
		case "organizationalUnit":
			accountIDs := make(map[string]string)
			for _, groupValue := range groupValues {
				if node, ok := graph.idMap[groupValue]; ok {
					if node.IsOrganizationalUnit() || node.IsRoot() {
						for _, accountID := range node.DescendantAccountIDs() {
							accountIDs[accountID] = accountID
						}
					}
				}
			}
			keys := make([]string, 0, len(accountIDs))
			for k := range accountIDs {
				keys = append(keys, k)
			}

			return keys, nil
		// case "tag":
		// 	var tags []resourcegroupstaggingapitypes.TagFilter
		// 	for _, gv := range groupValues {
		// 		kv := strings.SplitN(gv, ":", 2)
		// 		if len(kv) != 2 {
		// 			return nil, &providers.InvalidGroupValueError{GroupID: groupID, GroupValue: gv}
		// 		}
		// 		tags = append(tags, resourcegroupstaggingapitypes.TagFilter{
		// 			Key:    aws.String(kv[0]),
		// 			Values: []string{kv[1]},
		// 		})
		// 	}
		// 	return p.listAccountsWithTag(ctx, tags)
		default:
			return nil, &providers.InvalidGroupIDError{GroupID: groupID}
		}
	default:
		return nil, &providers.InvalidArgumentError{Arg: argId}
	}
}

// func (p *Provider) listAccountsWithTag(ctx context.Context, tags []resourcegroupstaggingapitypes.TagFilter) (accounts []string, err error) {
// 	hasMore := true
// 	var nextToken *string
// 	for hasMore {
// 		resources, err := p.resourcesClient.GetResources(ctx, &resourcegroupstaggingapi.GetResourcesInput{
// 			// TagFilters:          tags,
// 			ResourceTypeFilters: []string{"organizations:account"},
// 			PaginationToken:     nextToken,
// 		})
// 		if err != nil {
// 			return nil, err
// 		}
// 		nextToken = resources.PaginationToken
// 		if nextToken != nil && *nextToken == "" {
// 			nextToken = nil
// 		}
// 		hasMore = nextToken != nil

// 		// Split the account id from the arn
// 		for _, resource := range resources.ResourceTagMappingList {
// 			s := strings.Split(aws.ToString(resource.ResourceARN), "/")
// 			accounts = append(accounts, s[len(s)-1])
// 		}
// 	}
// 	return
// }

// func (p *Provider) listTagsForAccount(ctx context.Context, accountID string) (tags []organizationTypes.Tag, err error) {
// 	hasMore := true
// 	var nextToken *string
// 	for hasMore {
// 		acctTags, err := p.orgClient.ListTagsForResource(ctx, &organizations.ListTagsForResourceInput{
// 			ResourceId: &accountID,
// 			NextToken:  nextToken,
// 		})
// 		if err != nil {
// 			return nil, err
// 		}
// 		nextToken = acctTags.NextToken
// 		hasMore = nextToken != nil
// 		tags = append(tags, acctTags.Tags...)
// 	}
// 	return
// }

// func (p *Provider) generateTagGroupOptionsForAccounts(ctx context.Context, accounts []organizationTypes.Account) ([]types.GroupOption, error) {
// 	groupOptions := []types.GroupOption{}
// 	tagAccountMap := make(map[string][]string)
// 	var mu sync.Mutex
// 	// commented out all the go routines because it was causing a context cancelled error
// 	// g, gctx := errgroup.WithContext(ctx)
// 	// g.SetLimit(1) // set a limit here to avoid hitting API rate limits in cases where accounts have many permission sets
// 	for _, acct := range accounts {
// 		// g.Go(func() error {
// 		tags, err := p.listTagsForAccount(ctx, aws.ToString(acct.Id))
// 		if err != nil {
// 			return nil, err
// 		}
// 		mu.Lock()
// 		for _, tag := range tags {
// 			// Note: tags are key value pairs and we need both to look them up, we join them with a :
// 			// TODO:consider adding native key:value pair support for option values?
// 			// If required make the value some opaque encoded value if its difficult to store
// 			kv := aws.ToString(tag.Key) + ":" + aws.ToString(tag.Value)
// 			tagAccounts := tagAccountMap[kv]
// 			tagAccounts = append(tagAccounts, aws.ToString(acct.Id))
// 			tagAccountMap[kv] = tagAccounts
// 		}
// 		mu.Unlock()
// 		// return nil
// 		// })
// 		// err = g.Wait()
// 		if err != nil {
// 			return nil, err
// 		}
// 	}
// 	for k, v := range tagAccountMap {
// 		groupOptions = append(groupOptions, types.GroupOption{
// 			Children: v,
// 			Label:    k,
// 			Value:    k,
// 		})
// 	}
// 	return groupOptions, nil
// }
