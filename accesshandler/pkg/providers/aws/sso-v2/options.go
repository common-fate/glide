package ssov2

import (
	"context"
	"errors"
	"sync"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/organizations"
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

		hasMore = true
		nextToken = nil
		roots, err := p.orgClient.ListRoots(ctx, &organizations.ListRootsInput{})
		if err != nil {
			return nil, err
		}
		orgUnitGroup := types.Group{
			Title: "Organizational Unit",
			Id:    "organizationalUnit",
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
				option := types.GroupOption{Label: aws.ToString(orgUnit.Name), Value: aws.ToString(orgUnit.Id)}
				ouHasMore := true
				var ouNextToken *string
				for ouHasMore {
					children, err := p.orgClient.ListAccountsForParent(ctx, &organizations.ListAccountsForParentInput{
						ParentId:  orgUnit.Id,
						NextToken: ou.NextToken,
					})
					if err != nil {
						return nil, err
					}
					for _, a := range children.Accounts {
						option.Children = append(option.Children, aws.ToString(a.Id))
					}
					ouNextToken = children.NextToken
					ouHasMore = ouNextToken != nil
				}
				orgUnitGroup.Options = append(orgUnitGroup.Options, option)
			}
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
