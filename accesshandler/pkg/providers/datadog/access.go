package datadog

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/DataDog/datadog-api-client-go/v2/api/datadogV1"
	"github.com/DataDog/datadog-api-client-go/v2/api/datadogV2"
	"github.com/pkg/errors"
	"go.uber.org/zap"
)

type Args struct {
	Dashboard string `json:"dashboard"`
}

// Grant the access.
func (p *Provider) Grant(ctx context.Context, subject string, args []byte, grantID string) error {
	var a Args
	err := json.Unmarshal(args, &a)
	if err != nil {
		return err
	}
	ctx = p.DDContext(ctx)

	log := zap.S().With("args", a)
	log.Info("getting dashboard")

	// create a role for the grant ID
	roleAPI := datadogV2.NewRolesApi(p.apiClient)

	createRole := datadogV2.RoleCreateRequest{
		Data: datadogV2.RoleCreateData{
			Type: datadogV2.ROLESTYPE_ROLES.Ptr(),
			Attributes: datadogV2.RoleCreateAttributes{
				Name: grantID,
			},
		},
	}

	role, _, err := roleAPI.CreateRole(ctx, createRole)
	if err != nil {
		return errors.Wrap(err, "creating role")
	}

	// get the user
	user, err := p.FindUserByEmail(ctx, subject)
	if err != nil {
		return errors.Wrap(err, "find user by email")
	}

	// add the user to the role
	_, _, err = roleAPI.AddUserToRole(ctx, *role.GetData().Id, datadogV2.RelationshipToUser{
		Data: datadogV2.RelationshipToUserData{
			Id:   user.GetId(),
			Type: datadogV2.USERSTYPE_USERS,
		},
	})
	if err != nil {
		return errors.Wrap(err, "add user to role")
	}

	// get the dashboard
	dapi := datadogV1.NewDashboardsApi(p.apiClient)
	dashboard, _, err := dapi.GetDashboard(ctx, a.Dashboard)
	if err != nil {
		return errors.Wrap(err, "get dashboard")
	}

	// add the role to the dashboard
	dashboard.RestrictedRoles = append(dashboard.RestrictedRoles, role.Data.GetId())

	time.Sleep(time.Second * 2)

	// update the dashboard
	_, _, err = dapi.UpdateDashboard(ctx, a.Dashboard, dashboard)
	if err != nil {
		return errors.Wrap(err, "update dashboard")
	}

	return nil
}

func (p *Provider) FindUserByEmail(ctx context.Context, email string) (*datadogV2.User, error) {
	api := datadogV2.NewUsersApi(p.apiClient)
	ctx = p.DDContext(ctx)
	resp, _, err := api.ListUsers(ctx, *datadogV2.NewListUsersOptionalParameters().WithFilter(email))
	if err != nil {
		return nil, err
	}

	users := resp.GetData()
	if len(users) == 0 {
		return nil, fmt.Errorf("no users found for email %s", email)
	}
	if len(users) > 1 {
		return nil, fmt.Errorf("multiple users found for email %s", email)
	}
	return &users[0], nil
}

// Revoke the access.
func (p *Provider) Revoke(ctx context.Context, subject string, args []byte, grantID string) error {
	var a Args
	err := json.Unmarshal(args, &a)
	if err != nil {
		return err
	}
	ctx = p.DDContext(ctx)

	// delete the role
	roleAPI := datadogV2.NewRolesApi(p.apiClient)

	role, err := p.FindRoleByName(ctx, grantID)
	if err != nil {
		return errors.Wrap(err, "finding role")
	}

	_, err = roleAPI.DeleteRole(ctx, role.GetId())
	if err != nil {
		return errors.Wrap(err, "deleting role")
	}

	return nil
}

func (p *Provider) FindRoleByName(ctx context.Context, name string) (*datadogV2.Role, error) {
	api := datadogV2.NewRolesApi(p.apiClient)
	ctx = p.DDContext(ctx)
	resp, _, err := api.ListRoles(p.DDContext(ctx), *datadogV2.NewListRolesOptionalParameters().WithFilter(name))
	if err != nil {
		return nil, err
	}

	roles := resp.GetData()
	if len(roles) == 0 {
		return nil, fmt.Errorf("no roles found for name %s", name)
	}
	if len(roles) > 1 {
		return nil, fmt.Errorf("multiple roles found for name %s", name)
	}
	return &roles[0], nil
}
