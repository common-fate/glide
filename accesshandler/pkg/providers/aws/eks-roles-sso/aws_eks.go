package eksrolessso

import (
	"context"
	"encoding/base64"
	"errors"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/eks"
	"github.com/aws/aws-sdk-go-v2/service/identitystore"
	"github.com/aws/aws-sdk-go-v2/service/organizations"
	"github.com/aws/aws-sdk-go-v2/service/ssoadmin"
	"github.com/common-fate/granted-approvals/pkg/gconfig"
	"github.com/invopop/jsonschema"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
	clientcmdapi "k8s.io/client-go/tools/clientcmd/api"
	"sigs.k8s.io/aws-iam-authenticator/pkg/token"
)

type Provider struct {
	kubeClient    *kubernetes.Clientset
	ssoClient     *ssoadmin.Client
	clusterName   gconfig.StringValue
	namespace     gconfig.StringValue
	idStoreClient *identitystore.Client
	orgClient     *organizations.Client
	instanceARN   gconfig.StringValue
	// The globally unique identifier for the identity store, such as d-1234567890.
	identityStoreID gconfig.StringValue
	// The aws region where the identity store runs
	region gconfig.OptionalStringValue
}

func (p *Provider) Config() gconfig.Config {
	return gconfig.Config{
		gconfig.StringField("clusterName", &p.clusterName, "The EKS cluster name"),
		gconfig.StringField("namespace", &p.namespace, "The kubernetes cluster namespace"),
		gconfig.StringField("identityStoreId", &p.identityStoreID, "the AWS SSO Identity Store ID"),
		gconfig.StringField("instanceArn", &p.instanceARN, "the AWS SSO Instance ARN"),
		gconfig.OptionalStringField("region", &p.region, "the region the AWS SSO instance is deployed to"),
	}
}

func (p *Provider) Init(ctx context.Context) error {
	cfg, err := config.LoadDefaultConfig(ctx)
	if err != nil {
		return err
	}
	creds, err := cfg.Credentials.Retrieve(ctx)
	if err != nil {
		return err
	}
	if creds.Expired() {
		return errors.New("AWS credentials are expired")
	}
	eksClient := eks.NewFromConfig(cfg)
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

	p.ssoClient = ssoadmin.NewFromConfig(cfg)
	p.orgClient = organizations.NewFromConfig(cfg)
	p.idStoreClient = identitystore.NewFromConfig(cfg)

	return nil
}

// ArgSchema returns the schema for the AWS SSO provider.
func (p *Provider) ArgSchema() *jsonschema.Schema {
	return jsonschema.Reflect(&Args{})
}
