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
			for _, acct := range o.Accounts {
				opts.Options = append(opts.Options, types.Option{Label: aws.ToString(acct.Name), Value: aws.ToString(acct.Id)})
			}
		}
		log.Info("getting aws organization unit id set options")

		roots, err := p.orgClient.ListRoots(ctx, &organizations.ListRootsInput{})
		if err != nil {
			return nil, err
		}

		orgUnitGroup := types.Group{
			Title: "Organizational Unit",
			Id:    "organizationalUnit",
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
			orgUnitGroup.Options = append(orgUnitGroup.Options, option)
		}
		opts.Groups = &types.Groups{
			AdditionalProperties: map[string]types.Group{
				"organizationalUnit": orgUnitGroup,
				"tags": {
					Title:   "Tag Name",
					Id:      "tag",
					Options: []types.GroupOption{{Label: "abc", Value: "123"}, {Label: "xyz", Value: "456"}},
				},
			},
		}

		return &opts, nil
	}

	return nil, &providers.InvalidArgumentError{Arg: arg}

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
