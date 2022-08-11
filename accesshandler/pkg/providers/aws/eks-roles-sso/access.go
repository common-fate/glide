package eksrolessso

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/iam"
	iamtypes "github.com/aws/aws-sdk-go-v2/service/iam/types"
	"github.com/aws/aws-sdk-go-v2/service/identitystore"
	idtypes "github.com/aws/aws-sdk-go-v2/service/identitystore/types"
	"github.com/labstack/gommon/log"
	"go.uber.org/zap"

	"github.com/aws/aws-sdk-go-v2/service/ssoadmin"
	"github.com/aws/aws-sdk-go-v2/service/ssoadmin/types"
	"github.com/common-fate/granted-approvals/pkg/cfaws/policy"
	"github.com/pkg/errors"
	"github.com/sethvargo/go-retry"
	"gopkg.in/yaml.v2"
	v1 "k8s.io/api/rbac/v1"
	v1meta "k8s.io/apimachinery/pkg/apis/meta/v1"
)

type Args struct {
	Role string `json:"role" jsonschema:"title=Role"`
}

func (p *Provider) Grant(ctx context.Context, subject string, args []byte, grantID string) error {
	log := zap.S().With("args", args)
	log.Info("granting with eks provider")
	var a Args
	err := json.Unmarshal(args, &a)
	if err != nil {
		return err
	}
	permissionSetName := permissionSetNameFromGrantID(grantID)

	log.Info("adding user to permission set ", permissionSetName)
	// Create and assign user to permission set for this grant
	roleARN, err := p.createPermissionSetAndAssignment(ctx, subject, permissionSetName)
	if err != nil {
		return err
	}

	// Create a kubernetes role-binding for the object key as user to the kubernetes Role
	err = p.createKubernetesRoleBinding(ctx, objectKeyFromGrantID(grantID), a.Role)
	if err != nil {
		return err
	}

	// Assign the aws IAM role from the permission set assignment to the objectID user in kubernetes
	return p.createAWSAuthConfigMapRoleMapEntry(ctx, roleARN, objectKeyFromGrantID(grantID))
}

func (p *Provider) Revoke(ctx context.Context, subject string, args []byte, grantID string) error {
	var a Args
	err := json.Unmarshal(args, &a)
	if err != nil {
		return err
	}

	permissionSetName := permissionSetNameFromGrantID(grantID)
	roleARN, err := p.getSanitisedRoleARNForPermissionSetAssignment(ctx, permissionSetName)
	if err != nil {
		return err
	}
	// Remove the aws-auth config map entry
	err = p.removeAWSAuthConfigMapRoleMapEntry(ctx, roleARN)
	if err != nil {
		return err
	}
	// Remove the role binding
	err = p.kubeClient.RbacV1().RoleBindings(p.namespace.Get()).Delete(ctx, objectKeyFromGrantID(grantID), v1meta.DeleteOptions{})
	if err != nil {
		return err
	}

	return p.removePermissionSet(ctx, permissionSetName)
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
func (p *Provider) Instructions(ctx context.Context, subject string, args []byte) (string, error) {
	url := fmt.Sprintf("https://%s.awsapps.com/start", p.identityStoreID)
	instructions := fmt.Sprintf("You can access this role at your [AWS SSO URL](%s)\nOr use `assume --sso --sso-start-url %s --sso-region %s --account-id %s --role-name <Replace with your requestID>`", url, url, p.ssoRegion.Get(), p.awsAccountID)
	return instructions, nil
}
func objectKeyFromGrantID(grantID string) string {
	return fmt.Sprintf("granted-approvals-%s", grantID)
}

// Permission set names have a maximum length of 32, in normal use a KSUID will be the grant ID so this should never get truncated
// however if it is > 32 chars it will be truncated
func permissionSetNameFromGrantID(grantID string) string {
	permissionSetName := grantID
	if len(permissionSetName) > 32 {
		permissionSetName = permissionSetName[:32]
	}
	return permissionSetName
}

// createAWSAuthConfigMapRoleMapEntry appends an entry in the mapRoles section of the aws-auth config map
// by first fetching the current config and appending the new entry to the list, then updating the config map
func (p *Provider) createAWSAuthConfigMapRoleMapEntry(ctx context.Context, roleARN string, objectKey string) error {
	log.Info("get k8s config map: aws-auth")

	awsAuth, err := p.kubeClient.CoreV1().ConfigMaps("kube-system").Get(ctx, "aws-auth", v1meta.GetOptions{})
	if err != nil {
		return err
	}

	type MapRoleEntry struct {
		RoleARN  *string `yaml:"rolearn,omitempty"`
		Username *string `yaml:"username,omitempty"`
	}
	var dc []interface{}
	err = yaml.NewDecoder(bytes.NewBufferString(awsAuth.Data["mapRoles"])).Decode(&dc)
	if err != nil {
		return err
	}

	dc = append(dc, MapRoleEntry{RoleARN: aws.String(roleARN), Username: aws.String(objectKey)})
	var buf bytes.Buffer
	err = yaml.NewEncoder(&buf).Encode(dc)
	if err != nil {
		return err
	}
	awsAuth.Data["mapRoles"] = buf.String()
	log.Info("update k8s config map: aws-auth: ", awsAuth)

	_, err = p.kubeClient.CoreV1().ConfigMaps("kube-system").Update(ctx, awsAuth, v1meta.UpdateOptions{})
	return err
}

// removeAWSAuthConfigMapRoleMapEntry removes an entry from the config map if it exists by matching the roleARN
func (p *Provider) removeAWSAuthConfigMapRoleMapEntry(ctx context.Context, roleARN string) error {
	awsAuth, err := p.kubeClient.CoreV1().ConfigMaps("kube-system").Get(ctx, "aws-auth", v1meta.GetOptions{})
	if err != nil {
		return err
	}

	type mapRoleEntry struct {
		RoleARN  *string `yaml:"rolearn,omitempty"`
		Username *string `yaml:"username,omitempty"`
	}
	var dc []mapRoleEntry
	err = yaml.NewDecoder(bytes.NewBufferString(awsAuth.Data["mapRoles"])).Decode(&dc)
	if err != nil {
		return err
	}

	found := -1
	for i, entry := range dc {
		if entry.RoleARN != nil && *entry.RoleARN == roleARN {
			found = i
			break
		}
	}
	if found != -1 {
		dc = append(dc[:found], dc[found+1:]...)
	}

	var buf bytes.Buffer
	err = yaml.NewEncoder(&buf).Encode(dc)
	if err != nil {
		return err
	}
	awsAuth.Data["mapRoles"] = buf.String()
	_, err = p.kubeClient.CoreV1().ConfigMaps("kube-system").Update(ctx, awsAuth, v1meta.UpdateOptions{})
	return err
}

// createKubernetesRoleBinding uses the kubernetes API to create a role binding for use in the grant
func (p *Provider) createKubernetesRoleBinding(ctx context.Context, objectKey string, kubernetesRoleName string) error {
	rb := v1.RoleBinding{
		TypeMeta: v1meta.TypeMeta{Kind: "RoleBinding", APIVersion: "rbac.authorization.k8s.io/v1"},
		// use the key for the name
		ObjectMeta: v1meta.ObjectMeta{Name: objectKey},
		// use the key as the user too
		Subjects: []v1.Subject{{Kind: "User", APIGroup: "rbac.authorization.k8s.io", Name: objectKey, Namespace: p.namespace.Get()}},
		RoleRef:  v1.RoleRef{APIGroup: "rbac.authorization.k8s.io", Kind: "Role", Name: kubernetesRoleName},
	}
	log.Info("create kubernetes role binding ", rb)
	_, err := p.kubeClient.RbacV1().RoleBindings(p.namespace.Get()).Create(ctx, &rb, v1meta.CreateOptions{})
	return err
}
func (p *Provider) removePermissionSet(ctx context.Context, permissionSetName string) error {
	hasMore := true
	var nextToken *string
	var arnMatch *string
	for hasMore {
		o, err := p.ssoClient.ListPermissionSets(ctx, &ssoadmin.ListPermissionSetsInput{
			InstanceArn: aws.String(p.instanceARN.Get()),
			NextToken:   nextToken,
		})
		if err != nil {
			return err
		}
		nextToken = o.NextToken
		hasMore = nextToken != nil

		for _, arn := range o.PermissionSets {
			po, err := p.ssoClient.DescribePermissionSet(ctx, &ssoadmin.DescribePermissionSetInput{
				InstanceArn: aws.String(p.instanceARN.Get()), PermissionSetArn: aws.String(arn),
			})
			if err != nil {
				return err
			}
			if aws.ToString(po.PermissionSet.Name) == permissionSetName {
				arnMatch = po.PermissionSet.PermissionSetArn
				break
			}
		}
		if arnMatch != nil {
			break
		}
	}
	// Permission set does not exist, do nothing
	if arnMatch == nil {
		return nil
	}
	_, err := p.ssoClient.DeletePermissionSet(ctx, &ssoadmin.DeletePermissionSetInput{
		InstanceArn:      aws.String(p.instanceARN.Get()),
		PermissionSetArn: arnMatch,
	})
	return err
}

// createPermissionSetAndAssignment creates a permission set with a name = grantID
func (p *Provider) createPermissionSetAndAssignment(ctx context.Context, subject string, permissionSetName string) (roleARN string, err error) {
	eksPolicyDocument := policy.Policy{
		Version: "2012-10-17",
		Statements: []policy.Statement{
			{
				Effect: "Allow",
				Action: []string{
					"eks:AccessKubernetesApi",
				},
				Resource: []string{fmt.Sprintf("arn:aws:eks:%s:%s:cluster/%s", p.clusterRegion.Get(), p.awsAccountID, p.clusterName.Get())},
			},
		},
	}

	// find the user ID from the provided email address.
	user, err := p.getUser(ctx, subject)
	if err != nil {
		return "", err
	}
	// create permission set with policy
	permSet, err := p.ssoClient.CreatePermissionSet(ctx, &ssoadmin.CreatePermissionSetInput{
		InstanceArn: aws.String(p.instanceARN.Get()),
		Name:        aws.String(permissionSetName),
		Description: aws.String("This permission set was automatically generated by Granted Approvals"),
	})
	if err != nil {
		return "", err
	}
	// Assign eks policy to permission set
	_, err = p.ssoClient.PutInlinePolicyToPermissionSet(ctx, &ssoadmin.PutInlinePolicyToPermissionSetInput{
		InlinePolicy:     aws.String(eksPolicyDocument.String()),
		InstanceArn:      aws.String(p.instanceARN.Get()),
		PermissionSetArn: permSet.PermissionSet.PermissionSetArn,
	})
	if err != nil {
		return "", err
	}

	// assign user to permission set
	res, err := p.ssoClient.CreateAccountAssignment(ctx, &ssoadmin.CreateAccountAssignmentInput{
		InstanceArn:      aws.String(p.instanceARN.Get()),
		PermissionSetArn: permSet.PermissionSet.PermissionSetArn,
		PrincipalType:    types.PrincipalTypeUser,
		PrincipalId:      user.UserId,
		TargetId:         &p.awsAccountID,
		TargetType:       types.TargetTypeAwsAccount,
	})
	if err != nil {
		return "", err
	}

	if res.AccountAssignmentCreationStatus.FailureReason != nil {
		return "", fmt.Errorf("failed creating account assignment: %s", *res.AccountAssignmentCreationStatus.FailureReason)
	}
	return p.getSanitisedRoleARNForPermissionSetAssignment(ctx, permissionSetName)
}

// The role ARN for a permission set role includes the following substring which needs to be removed.
// when a user gets credentails and accesses the kubernetes API, this portion of the ARN is not present!
// so if it is left in, the role mapping will fail with an unhelpful error
func (p *Provider) getSanitisedRoleARNForPermissionSetAssignment(ctx context.Context, permissionSetName string) (string, error) {
	// fetch the new IAM role associated with the permission set assignment
	role, err := p.getIAMRoleForPermissionSetWithRetry(ctx, permissionSetName)
	if err != nil {
		return "", err
	}

	substringToRemove := fmt.Sprintf("aws-reserved/sso.amazonaws.com/%s/", p.ssoRegion.Get())
	return strings.Replace(*role.Arn, substringToRemove, "", 1), nil
}

// getIAMRoleForPermissionSetWithRetry uses a retry function to try to fetch the role that was created after assigning a user to a permission set
// the process takes around 30 seconds normally and the role ARN is partially autogenerated so we need to do a list and check for a name prefix.
func (p *Provider) getIAMRoleForPermissionSetWithRetry(ctx context.Context, permissionSetName string) (*iamtypes.Role, error) {
	var roleOutput *iamtypes.Role
	b := retry.NewFibonacci(time.Second)
	b = retry.WithMaxDuration(time.Second*60, b)
	err := retry.Do(ctx, b, func(ctx context.Context) (err error) {
		var marker *string
		hasMore := true

		// This is the path prefix assigned to all roles generated by SSO
		ssoPathPrefix := fmt.Sprintf("/aws-reserved/sso.amazonaws.com/%s/", p.ssoRegion.Get())
		roleNamePrefix := "AWSReservedSSO_" + permissionSetName
		for hasMore {
			listRolesResponse, err := p.iamClient.ListRoles(ctx, &iam.ListRolesInput{PathPrefix: aws.String(ssoPathPrefix), Marker: marker})
			if err != nil {
				return retry.RetryableError(err)
			}
			marker = listRolesResponse.Marker
			hasMore = listRolesResponse.IsTruncated
			for _, role := range listRolesResponse.Roles {
				if strings.HasPrefix(aws.ToString(role.RoleName), roleNamePrefix) {
					r := role
					roleOutput = &r
				}
			}
		}
		if roleOutput == nil {
			return retry.RetryableError(errors.New("role not yet available or does not exist"))
		}
		return nil
	})
	if err != nil {
		return nil, err
	}
	if roleOutput == nil {
		return nil, errors.New("role not found after assiging permission set")
	}

	return roleOutput, nil
}

// getUser retrieves the AWS SSO user from a provided email address.
func (p *Provider) getUser(ctx context.Context, email string) (*idtypes.User, error) {
	res, err := p.idStoreClient.ListUsers(ctx, &identitystore.ListUsersInput{
		IdentityStoreId: aws.String(p.identityStoreID.Get()),
		Filters: []idtypes.Filter{{
			AttributePath:  aws.String("UserName"),
			AttributeValue: aws.String(email),
		}},
	})
	if err != nil {
		return nil, err
	}
	if len(res.Users) == 0 {
		return nil, &UserNotFoundError{Email: email}
	}
	if len(res.Users) > 1 {
		// this should never happen, but check it anyway.
		return nil, fmt.Errorf("expected 1 user but found %v", len(res.Users))
	}

	return &res.Users[0], nil
}
