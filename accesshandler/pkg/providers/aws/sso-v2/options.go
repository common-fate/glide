package ssov2

import (
	"context"
	"errors"
	"strings"
	"sync"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/organizations"
	organizationTypes "github.com/aws/aws-sdk-go-v2/service/organizations/types"
	"github.com/aws/aws-sdk-go-v2/service/resourcegroupstaggingapi"
	resourcegroupstaggingapitypes "github.com/aws/aws-sdk-go-v2/service/resourcegroupstaggingapi/types"
	"github.com/aws/aws-sdk-go-v2/service/ssoadmin"
	ssoadmintypes "github.com/aws/aws-sdk-go-v2/service/ssoadmin/types"
	"github.com/common-fate/granted-approvals/accesshandler/pkg/providers"
	"github.com/common-fate/granted-approvals/accesshandler/pkg/types"
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
					var label string
					if po.PermissionSet.Name != nil {
						label = *po.PermissionSet.Name
					}
					if po.PermissionSet.Description != nil {
						label = label + ": " + *po.PermissionSet.Description
					}
					opts.Options = append(opts.Options, types.Option{Label: label, Value: ARNCopy})

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
		accounts, err := p.listAccountsForOrganization(ctx)
		if err != nil {
			return nil, err
		}
		for _, acct := range accounts {
			opts.Options = append(opts.Options, types.Option{Label: aws.ToString(acct.Name), Value: aws.ToString(acct.Id)})
		}
		log.Info("getting aws organization unit id set options")
		ous, err := p.generateOuGroupOptions(ctx)
		if err != nil {
			return nil, err
		}
		orgUnitGroup := types.Group{
			Title:   "Organizational Unit",
			Id:      "organizationalUnit",
			Options: ous,
		}

		tags, err := p.generateTagGroupOptionsForAccounts(ctx, accounts)
		if err != nil {
			return nil, err
		}
		tagGroup := types.Group{
			Title:   "Tag Name",
			Id:      "tag",
			Options: tags,
		}

		opts.Groups = &types.Groups{
			AdditionalProperties: map[string]types.Group{
				"organizationalUnit": orgUnitGroup,
				"tags":               tagGroup,
			},
		}

		return &opts, nil
	}

	return nil, &providers.InvalidArgumentError{Arg: arg}

}

func (p *Provider) listAccountsForOrganization(ctx context.Context) (accounts []organizationTypes.Account, err error) {
	hasMore := true
	var nextToken *string
	for hasMore {
		o, err := p.orgClient.ListAccounts(ctx, &organizations.ListAccountsInput{
			NextToken: nextToken,
		})
		if err != nil {
			return nil, err
		}
		nextToken = o.NextToken
		hasMore = nextToken != nil
		accounts = append(accounts, o.Accounts...)
	}
	return
}
func (p *Provider) listChildOusForParent(ctx context.Context, parentID string) (ous []organizationTypes.OrganizationalUnit, err error) {
	hasMore := true
	var nextToken *string
	for hasMore {
		ou, err := p.orgClient.ListOrganizationalUnitsForParent(ctx, &organizations.ListOrganizationalUnitsForParentInput{
			ParentId: aws.String(parentID),
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

func (p *Provider) listAccountsWithTag(ctx context.Context, tags []resourcegroupstaggingapitypes.TagFilter) (accounts []string, err error) {
	hasMore := true
	var nextToken *string
	for hasMore {
		resources, err := p.resourcesClient.GetResources(ctx, &resourcegroupstaggingapi.GetResourcesInput{
			TagFilters:          tags,
			ResourceTypeFilters: []string{"AWS::Organizations::Account"},
			PaginationToken:     nextToken,
		})
		if err != nil {
			return nil, err
		}
		nextToken = resources.PaginationToken
		hasMore = nextToken != nil

		// Split the account id from the arn
		for _, resource := range resources.ResourceTagMappingList {
			s := strings.Split(aws.ToString(resource.ResourceARN), "/")
			accounts = append(accounts, s[len(s)-1])
		}
	}
	return
}

func (p *Provider) listTagsForAccount(ctx context.Context, accountID string) (tags []organizationTypes.Tag, err error) {
	hasMore := true
	var nextToken *string
	for hasMore {
		acctTags, err := p.orgClient.ListTagsForResource(ctx, &organizations.ListTagsForResourceInput{
			ResourceId: &accountID,
			NextToken:  nextToken,
		})
		if err != nil {
			return nil, err
		}
		nextToken = acctTags.NextToken
		hasMore = nextToken != nil
		tags = append(tags, acctTags.Tags...)
	}
	return
}

func (p *Provider) generateTagGroupOptionsForAccounts(ctx context.Context, accounts []organizationTypes.Account) (groupOptions []types.GroupOption, err error) {
	tagAccountMap := make(map[string][]string)
	var mu sync.Mutex
	g, gctx := errgroup.WithContext(ctx)
	g.SetLimit(5) // set a limit here to avoid hitting API rate limits in cases where accounts have many permission sets
	for _, acct := range accounts {
		g.Go(func() error {
			tags, err := p.listTagsForAccount(gctx, aws.ToString(acct.Id))
			if err != nil {
				return err
			}
			mu.Lock()
			for _, tag := range tags {
				kv := aws.ToString(tag.Key) + ":" + aws.ToString(tag.Value)
				tagAccounts := tagAccountMap[kv]
				tagAccounts = append(tagAccounts, aws.ToString(acct.Id))
				tagAccountMap[kv] = tagAccounts
			}
			mu.Unlock()
			return nil
		})
		err = g.Wait()
		if err != nil {
			return nil, err
		}
	}
	for k, v := range tagAccountMap {
		groupOptions = append(groupOptions, types.GroupOption{
			Children: v,
			Label:    k,
			Value:    k,
		})
	}
	return
}

func (p *Provider) generateOuGroupOptions(ctx context.Context) (groupOptions []types.GroupOption, err error) {
	roots, err := p.orgClient.ListRoots(ctx, &organizations.ListRootsInput{})
	if err != nil {
		return nil, err
	}
	// @TODO this is only 1 level of OUs, it will not return any nested OUs
	childOus, err := p.listChildOusForParent(ctx, aws.ToString(roots.Roots[0].Id))
	if err != nil {
		return nil, err
	}
	for _, orgUnit := range childOus {
		option := types.GroupOption{Label: aws.ToString(orgUnit.Name), Value: aws.ToString(orgUnit.Id)}
		// @TODO this is only 1 level of accounts, it will not return accounts for nested OUs
		childAccounts, err := p.listChildAccountsForParent(ctx, aws.ToString(orgUnit.Id))
		if err != nil {
			return nil, err
		}
		for _, a := range childAccounts {
			option.Children = append(option.Children, aws.ToString(a.Id))
		}
		groupOptions = append(groupOptions, option)
	}
	return
}
