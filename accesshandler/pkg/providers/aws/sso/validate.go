package sso

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/identitystore"
	"github.com/aws/aws-sdk-go-v2/service/identitystore/types"
	"github.com/aws/aws-sdk-go-v2/service/organizations"
	orgtypes "github.com/aws/aws-sdk-go-v2/service/organizations/types"
	"github.com/aws/aws-sdk-go-v2/service/ssoadmin"
	"golang.org/x/sync/errgroup"
)

// Validate the access against AWS SSO without actually granting it.
// This provider requires that the user name matches the user's email address.
func (p *Provider) Validate(ctx context.Context, subject string, args []byte) error {
	var a Args
	err := json.Unmarshal(args, &a)
	if err != nil {
		return err
	}

	// run the validations concurrently, as we need to wait for the API to respond.
	g := new(errgroup.Group)

	// the user should exist in AWS SSO.
	g.Go(func() error {
		res, err := p.idStoreClient.ListUsers(ctx, &identitystore.ListUsersInput{
			IdentityStoreId: aws.String(p.identityStoreID.Get()),
			Filters: []types.Filter{{
				AttributePath:  aws.String("UserName"),
				AttributeValue: aws.String(subject),
			}},
		})
		if err != nil {
			return err
		}
		if len(res.Users) == 0 {
			return fmt.Errorf("could not find user %s in AWS SSO", subject)
		}
		if len(res.Users) > 1 {
			// this should never happen, but check it anyway.
			return fmt.Errorf("expected 1 user but found %v", len(res.Users))
		}
		return nil
	})

	// the permission set should exist. and be tagged as managed-by-granted
	g.Go(func() error {
		_, err = p.client.DescribePermissionSet(ctx, &ssoadmin.DescribePermissionSetInput{
			InstanceArn:      aws.String(p.instanceARN.Get()),
			PermissionSetArn: &a.PermissionSetARN,
		})
		if err != nil {
			return &PermissionSetNotFoundErr{PermissionSet: a.PermissionSetARN, AWSErr: err}
		}
		isManaged, err := p.checkPermissionSetIsManagedByGranted(ctx, a.PermissionSetARN)
		if err != nil || !isManaged {
			return &PermissionSetNotManagedByGrantedError{PermissionSet: a.PermissionSetARN, AWSErr: err}
		}
		return nil
	})

	// the account should exist.
	g.Go(func() error {
		return p.ensureAccountExists(ctx, a.AccountID)
	})

	return g.Wait()
}

func (p *Provider) ensureAccountExists(ctx context.Context, accountID string) error {
	_, err := p.orgClient.DescribeAccount(ctx, &organizations.DescribeAccountInput{
		AccountId: &accountID,
	})
	var anf *orgtypes.AccountNotFoundException
	if errors.As(err, &anf) {
		return &AccountNotFoundError{AccountID: accountID}
	}

	return err
}
