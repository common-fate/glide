package eksrolessso

import (
	"context"
	"encoding/base64"
	"errors"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/eks"
	"github.com/aws/aws-sdk-go-v2/service/iam"
	"github.com/aws/aws-sdk-go-v2/service/identitystore"
	"github.com/aws/aws-sdk-go-v2/service/organizations"
	"github.com/aws/aws-sdk-go-v2/service/ssoadmin"
	"github.com/aws/aws-sdk-go-v2/service/sts"
	"github.com/common-fate/granted-approvals/accesshandler/pkg/providers"
	"github.com/common-fate/granted-approvals/accesshandler/pkg/types"
	"github.com/common-fate/granted-approvals/pkg/cfaws"
	"github.com/common-fate/granted-approvals/pkg/gconfig"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
	clientcmdapi "k8s.io/client-go/tools/clientcmd/api"
	"sigs.k8s.io/aws-iam-authenticator/pkg/token"
)

type Provider struct {
	kubeClient    *kubernetes.Clientset
	ssoClient     *ssoadmin.Client
	iamClient     *iam.Client
	idStoreClient *identitystore.Client
	orgClient     *organizations.Client
	// the account that the eks cluster runs in, this is fetched after assuming the cluster role
	eksClusterRoleAccountID string

	// configured by gconfig
	clusterAccessRoleARN gconfig.StringValue
	clusterName          gconfig.StringValue
	namespace            gconfig.StringValue
	clusterRegion        gconfig.StringValue
	// sso instance
	instanceARN gconfig.StringValue
	// The globally unique identifier for the identity store, such as d-1234567890.
	identityStoreID gconfig.StringValue
	// The aws region where the identity store runs
	ssoRegion gconfig.StringValue
	// a role which can be assumed and has required sso permissions
	ssoRoleARN gconfig.StringValue
}

func (p *Provider) Config() gconfig.Config {
	return gconfig.Config{
		gconfig.StringField("clusterName", &p.clusterName, "The EKS cluster name"),
		gconfig.StringField("namespace", &p.namespace, "The kubernetes cluster namespace"),
		gconfig.StringField("clusterRegion", &p.clusterRegion, "The region the EKS cluster is deployed"),
		gconfig.StringField("clusterAccessRoleArn", &p.clusterAccessRoleARN, "The ARN of the AWS IAM Role with permission to access the EKS cluster"),
		gconfig.StringField("identityStoreId", &p.identityStoreID, "The AWS SSO Identity Store ID"),
		gconfig.StringField("instanceArn", &p.instanceARN, "The AWS SSO Instance ARN"),
		gconfig.StringField("ssoRegion", &p.ssoRegion, "The region the AWS SSO instance is deployed to"),
		gconfig.StringField("ssoRoleARN", &p.ssoRoleARN, "The ARN of the AWS IAM Role with permission to administer SSO"),
	}
}

func (p *Provider) Init(ctx context.Context) error {
	ssoCfg, err := config.LoadDefaultConfig(ctx, config.WithRegion(p.ssoRegion.Get()), config.WithCredentialsProvider(cfaws.NewAssumeRoleCredentialsCache(ctx, p.ssoRoleARN.Get(), cfaws.WithRoleSessionName("accesshandler-eks-roles-sso"))))
	if err != nil {
		return err
	}
	eksCfg, err := config.LoadDefaultConfig(ctx, config.WithRegion(p.clusterRegion.Get()), config.WithCredentialsProvider(cfaws.NewAssumeRoleCredentialsCache(ctx, p.clusterAccessRoleARN.Get(), cfaws.WithRoleSessionName("accesshandler-eks-roles-sso"))))
	if err != nil {
		return err
	}
	eksClient := eks.NewFromConfig(eksCfg)
	r, err := eksClient.DescribeCluster(ctx, &eks.DescribeClusterInput{Name: aws.String(p.clusterName.Get())})
	if err != nil {
		return err
	}

	pem, err := base64.StdEncoding.DecodeString(aws.ToString(r.Cluster.CertificateAuthority.Data))
	if err != nil {
		return err
	}
	c := clientcmdapi.Config{
		APIVersion: "v1",
		Kind:       "Config",
		Clusters: map[string]*clientcmdapi.Cluster{
			"cluster": {

				Server:                   aws.ToString(r.Cluster.Endpoint),
				CertificateAuthorityData: pem,
			},
		},
		Contexts: map[string]*clientcmdapi.Context{
			"cluster": {
				Cluster: "cluster",
			},
		},
		CurrentContext: "cluster",
	}
	cc := clientcmd.NewDefaultClientConfig(c, nil)
	kubeConfig, err := cc.ClientConfig()
	if err != nil {
		return err
	}

	g, err := token.NewGenerator(false, false)
	if err != nil {
		return err
	}
	token, err := g.Get(p.clusterName.Get())
	if err != nil {
		return err
	}
	kubeConfig.BearerToken = token.Token
	client, err := kubernetes.NewForConfig(kubeConfig)
	if err != nil {
		return err
	}
	p.kubeClient = client

	p.ssoClient = ssoadmin.NewFromConfig(ssoCfg)
	p.orgClient = organizations.NewFromConfig(ssoCfg)
	p.idStoreClient = identitystore.NewFromConfig(ssoCfg)
	p.iamClient = iam.NewFromConfig(ssoCfg)
	stsClient := sts.NewFromConfig(eksCfg)
	res, err := stsClient.GetCallerIdentity(ctx, &sts.GetCallerIdentityInput{})
	if err != nil {
		return err
	}
	if res.Account == nil {
		return errors.New("aws accountID was nil in sts get caller id response")
	}
	p.eksClusterRoleAccountID = *res.Account

	return nil
}

func (p *Provider) ArgSchemaV2() providers.ArgSchema {
	arg := providers.ArgSchema{
		"role": {
			Id:          "role",
			Title:       "Role",
			FormElement: types.MULTISELECT,
		},
	}
	return arg
}
