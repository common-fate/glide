package eksrolessso

import (
	"context"
	"encoding/json"
	"fmt"

	v1 "k8s.io/api/rbac/v1"
	v1meta "k8s.io/apimachinery/pkg/apis/meta/v1"
)

type Args struct {
	Role string `json:"role" jsonschema:"title=Role"`
}

func (p *Provider) Grant(ctx context.Context, subject string, args []byte, grantID string) error {
	var a Args
	err := json.Unmarshal(args, &a)
	if err != nil {
		return err
	}
	// how will we track all the items that we have created?
	// We only ever need to create a single copy of the iam policy using a create or update operation
	// The permission set is created per user and needs to be deleted afterwards, so how do we track this state, naming convention?
	//

	// create iam policy with eks permissions
	// create permission set with policy
	// create a kubernetes role-binding for subject to the kubernetes role
	// create a role map entry for the iam role of the permission set to the kubernetes user in the aws-auth config map
	// assign user to permission set
	rb := v1.RoleBinding{
		TypeMeta:   v1meta.TypeMeta{Kind: "RoleBinding", APIVersion: "rbac.authorization.k8s.io/v1"},
		ObjectMeta: v1meta.ObjectMeta{Name: fmt.Sprintf("granted-approvals-%s", subject, a.Role)},
	}
	p.kubeClient.RbacV1().RoleBindings(p.clusterName.Get()).Create(ctx, &rb, v1meta.CreateOptions{})
	// p.kubeClient.RbacV1().RoleBindings(p.clusterName.Get()).Delete(ctx, "", v1meta.DeleteOptions{})
	return nil
}

func (p *Provider) Revoke(ctx context.Context, subject string, args []byte, grantID string) error {
	var a Args
	err := json.Unmarshal(args, &a)
	if err != nil {
		return err
	}

	// reverse the process from grant step

	return err
}

func (p *Provider) IsActive(ctx context.Context, subject string, args []byte, grantID string) (bool, error) {
	var a Args
	err := json.Unmarshal(args, &a)
	if err != nil {
		return false, err
	}

	// we didn't find the user, so return false.
	return false, nil
}

// func (p *Provider) Instructions(ctx context.Context, subject string, args []byte) (string, error) {
// 	return "", nil
// }
