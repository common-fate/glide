package deploy

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/AlecAivazis/survey/v2"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/cloudformation/types"
	"github.com/aws/aws-sdk-go-v2/service/sts"
	"github.com/briandowns/spinner"
	"github.com/common-fate/granted-approvals/accesshandler/pkg/genv"
	"github.com/common-fate/granted-approvals/pkg/clio"
	"github.com/urfave/cli/v2"
	"go.uber.org/zap"
	"gopkg.in/yaml.v3"
)

var ErrConfigNotExist = errors.New("config does not exist")

const DefaultFilename = "granted-deployment.yml"

// AvailableRegions are the regions that we currently release CloudFormation templates to.
var AvailableRegions = []string{
	"ap-southeast-2",
	"us-west-2",
}

// AvailableSSOProviders are the currently implemented SSO providers.
var AvailableSSOProviders = []string{
	"Google",
	"Okta",
	"Azure",
}

type Config struct {
	Version       int                 `yaml:"version"`
	Deployment    Deployment          `yaml:"deployment"`
	Providers     map[string]Provider `yaml:"providers,omitempty"`
	Notifications *Notifications      `yaml:"notifications,omitempty"`
	Identity      *Identity           `yaml:"identity,omitempty"`

	cachedOutput     *Output
	cachedSAMLOutput *SAMLOutputs
}
type Identity struct {
	Google *Google `yaml:"google,omitempty" json:"google,omitempty"`
	Okta   *Okta   `yaml:"okta,omitempty" json:"okta,omitempty"`
	Azure  *Azure  `yaml:"azure,omitempty" json:"azure,omitempty"`
}

type Google struct {
	APIToken   string `yaml:"apiToken" json:"apiToken"`
	Domain     string `yaml:"domain" json:"domain"`
	AdminEmail string `yaml:"adminEmail" json:"adminEmail"`
}

// String redacts potentially sensitive token values
func (g Google) String() string {
	g.APIToken = "****"
	return fmt.Sprintf("{APIToken: %s, Domain: %s, AdminEmail %s}", g.APIToken, g.Domain, g.AdminEmail)
}

type Okta struct {
	APIToken string `yaml:"apiToken" json:"apiToken"`
	OrgURL   string `yaml:"orgUrl" json:"orgUrl"`
}

// String redacts potentially sensitive token values
func (o Okta) String() string {
	o.APIToken = "****"
	return fmt.Sprintf("{APIToken: %s, OrgURL: %s}", o.APIToken, o.OrgURL)
}

type Azure struct {
	TenantID     string `yaml:"tenantID" json:"tenantID"`
	ClientID     string `yaml:"clientID" json:"clientID"`
	ClientSecret string `yaml:"clientSecret" json:"clientSecret"`
}

// String redacts potentially sensitive token values
func (a Azure) String() string {
	a.ClientSecret = "****"
	return fmt.Sprintf("{TenantID: %s, ClientID: %s, ClientSecret: %s}", a.TenantID, a.ClientID, a.ClientSecret)
}

type Notifications struct {
	Slack *Slack `yaml:"slack,omitempty" json:"slack,omitempty"`
}

type Slack struct {
	APIToken string `yaml:"apiToken" json:"apiToken"`
}

// String redacts potentially sensitive token values
func (s Slack) String() string {
	s.APIToken = "****"
	return fmt.Sprintf("{APIToken: %s}", s.APIToken)
}

type Deployment struct {
	StackName string `yaml:"stackName"`
	Account   string `yaml:"account"`
	Region    string `yaml:"region"`
	Release   string `yaml:"release"`
	// Dev is set to true for internal development deployments only.
	Dev        *bool      `yaml:"dev,omitempty"`
	Parameters Parameters `yaml:"parameters"`
}

type Provider struct {
	Uses string            `yaml:"uses" json:"uses"`
	With map[string]string `yaml:"with" json:"with"`
}

type Parameters struct {
	CognitoDomainPrefix    string `yaml:"CognitoDomainPrefix"`
	AdministratorGroupID   string `yaml:"AdministratorGroupID"`
	DeploymentSuffix       string `yaml:"DeploymentSuffix,omitempty"`
	IdentityProviderType   string `yaml:"IdentityProviderType,omitempty"`
	SamlSSOMetadata        string `yaml:"SamlSSOMetadata,omitempty"`
	SamlSSOMetadataURL     string `yaml:"SamlSSOMetadataURL,omitempty"`
	FrontendDomain         string `yaml:"FrontendDomain,omitempty"`
	FrontendCertificateARN string `yaml:"FrontendCertificateARN,omitempty"`
}

// AddProvider adds a new provider to the deployment configuration.
func (c *Config) AddProvider(id string, p Provider) error {
	if c.Providers == nil {
		c.Providers = make(map[string]Provider)
	}
	if _, ok := c.Providers[id]; ok {
		return fmt.Errorf("provider %s already exists in the config", id)
	}
	c.Providers[id] = p
	return nil
}

func ProviderFromLookup(id string, uses string, with genv.Config) Provider {
	p := Provider{
		Uses: uses,
		With: make(map[string]string),
	}

	for _, v := range with {
		val := v.Get()

		if s, ok := v.(genv.Secret); ok && s.IsSecret() {
			val = fmt.Sprintf("awsssm:///granted/providers/%s/%s", id, v.Key())
		}

		p.With[v.Key()] = val
	}

	return p
}

// CfnParams converts the parameters to types supported by CloudFormation deployments.
func (c *Config) CfnParams() ([]types.Parameter, error) {
	p := c.Deployment.Parameters
	res := []types.Parameter{
		{
			ParameterKey:   aws.String("CognitoDomainPrefix"),
			ParameterValue: &p.CognitoDomainPrefix,
		},
	}
	if p.DeploymentSuffix != "" {
		res = append(res, types.Parameter{
			ParameterKey:   aws.String("DeploymentSuffix"),
			ParameterValue: &p.DeploymentSuffix,
		})
	}
	if c.Providers != nil {
		config, err := json.Marshal(c.Providers)
		if err != nil {
			return nil, err
		}
		configStr := string(config)
		res = append(res, types.Parameter{
			ParameterKey:   aws.String("ProviderConfiguration"),
			ParameterValue: &configStr,
		})
	}
	if c.Notifications != nil {
		if c.Notifications.Slack != nil {
			config, err := json.Marshal(c.Notifications.Slack)
			if err != nil {
				return nil, err
			}
			configStr := string(config)
			res = append(res, types.Parameter{
				ParameterKey:   aws.String("SlackConfiguration"),
				ParameterValue: &configStr,
			})
		}

	}
	if c.Identity != nil {
		config, err := json.Marshal(c.Identity)
		if err != nil {
			return nil, err
		}
		configStr := string(config)
		res = append(res, types.Parameter{
			ParameterKey:   aws.String("IdentityConfiguration"),
			ParameterValue: &configStr,
		})
	}
	if p.AdministratorGroupID != "" {
		res = append(res, types.Parameter{
			ParameterKey:   aws.String("AdministratorGroupID"),
			ParameterValue: &p.AdministratorGroupID,
		})
	}

	if c.Deployment.Parameters.IdentityProviderType != "" {
		res = append(res, types.Parameter{
			ParameterKey:   aws.String("IdentityProviderType"),
			ParameterValue: &p.IdentityProviderType,
		})
	}

	if c.Deployment.Parameters.SamlSSOMetadata != "" {
		res = append(res, types.Parameter{
			ParameterKey:   aws.String("SamlSSOMetadata"),
			ParameterValue: &p.SamlSSOMetadata,
		})
	}

	if c.Deployment.Parameters.SamlSSOMetadataURL != "" {
		res = append(res, types.Parameter{
			ParameterKey:   aws.String("SamlSSOMetadataURL"),
			ParameterValue: &p.SamlSSOMetadataURL,
		})
	}

	if c.Deployment.Parameters.FrontendCertificateARN != "" {
		res = append(res, types.Parameter{
			ParameterKey:   aws.String("FrontendCertificateARN"),
			ParameterValue: &p.FrontendCertificateARN,
		})
	}

	if c.Deployment.Parameters.FrontendDomain != "" {
		res = append(res, types.Parameter{
			ParameterKey:   aws.String("FrontendDomain"),
			ParameterValue: &p.FrontendDomain,
		})
	}

	return res, nil
}

// MustLoadConfig attempts to load config.
//
// if it does not exist, it will log a helpful message then os.Exit(0)
//
// if there are any other errors, they are logged and os.Exit(1)
func MustLoadConfig(f string) *Config {
	c, err := LoadConfig(f)
	if err == ErrConfigNotExist {
		clio.Error("Tried to load Granted deployment configuration from %s but the file doesn't exist.", f)
		clio.Log(`
To fix this, take one of the following actions:
  a) run this command from a folder which contains a Granted deployment configuration file (like 'granted-deployment.yml')
  b) run 'gdeploy init' to set up a new deployment configuration file
`)
		os.Exit(0)
	}
	if err != nil {
		clio.Error("failed to load config with error: %s", err)
		os.Exit(1)

	}
	return c
}

// LoadConfig attempts to load the config file
// if it does not exist, returns ErrConfigNotExist
// else returns the config or any other error
//
// for CLI commands where a helpful message is required, use MustLoadConfig, which will log a message and os.Exit() if it does not exist
func LoadConfig(f string) (*Config, error) {
	if _, err := os.Stat(f); errors.Is(err, os.ErrNotExist) {
		return nil, ErrConfigNotExist
	}

	fileRead, err := os.OpenFile(f, os.O_RDWR|os.O_CREATE, 0600)
	if err != nil {
		return nil, err
	}
	defer fileRead.Close()
	var dc Config
	err = yaml.NewDecoder(fileRead).Decode(&dc)
	if err != nil {
		return nil, err
	}
	return &dc, nil
}

func (c *Config) CfnTemplateURL() string {
	return fmt.Sprintf("https://granted-releases-%s.s3.amazonaws.com/%s/Granted.template.json", c.Deployment.Region, c.Deployment.Release)
}

func (c *Config) Save(f string) error {
	fileWrite, err := os.OpenFile(f, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0600)
	if err != nil {
		return err
	}
	defer fileWrite.Close()
	// save the config
	enc := yaml.NewEncoder(fileWrite)
	enc.SetIndent(2)
	return enc.Encode(c)
}

// NewStagingConfig sets up a Config for staging deployments.
// These deployments currently still use the CDK rather than CloudFormation.
func NewStagingConfig(ctx context.Context, stage string) *Config {
	acc, _ := tryGetCurrentAccountID(ctx)
	dev := true
	conf := Config{
		Deployment: Deployment{
			StackName: "granted-approvals-" + stage,
			Account:   acc,
			Dev:       &dev,
		},
	}
	return &conf
}

// SetupDevConfig sets up the config for local development.
func SetupDevConfig() (*Config, error) {
	var stage string
	p := &survey.Input{Message: "Enter a name for your deployment (you can use your name e.g. josh)"}
	err := survey.AskOne(p, &stage, survey.WithValidator(survey.MinLength(1)))
	if err != nil {
		return nil, err
	}
	ctx := context.Background()
	account, _ := tryGetCurrentAccountID(ctx)

	p = &survey.Input{Message: "Enter the account ID that you are deploying to", Default: account}
	err = survey.AskOne(p, &account)
	if err != nil {
		return nil, err
	}

	region := os.Getenv("AWS_REGION")

	p = &survey.Input{Message: "Enter the AWS region that you are deploying to", Default: region}
	err = survey.AskOne(p, &region)
	if err != nil {
		return nil, err
	}

	dev := true
	c := Config{
		Version: 1,
		Deployment: Deployment{
			StackName: fmt.Sprintf("granted-approvals-%s", stage),
			Account:   account,
			Region:    region,
			Dev:       &dev,
			Parameters: Parameters{
				AdministratorGroupID: "granted_administrators",
			},
		},
	}

	return &c, nil
}

// SetupReleaseConfig sets up the release configuration used in production deployments.
func SetupReleaseConfig(c *cli.Context) (*Config, error) {
	name := c.String("name")
	if name == "" {
		name = "Granted"
		p := &survey.Input{Message: "The name of the CloudFormation stack to create", Default: name}
		err := survey.AskOne(p, &name, survey.WithValidator(survey.MinLength(1)))
		if err != nil {
			return nil, err
		}
	}

	account := c.String("account")
	if account == "" {
		ctx := context.Background()

		account = MustGetCurrentAccountID(ctx, WithWarnExpiryIfWithinDuration(time.Minute))

		p := &survey.Input{Message: "The account ID that you are deploying to", Default: account}
		err := survey.AskOne(p, &account)
		if err != nil {
			return nil, err
		}
	}

	region := c.String("region")
	if region == "" {
		region = getDefaultAvailableRegion()
		p2 := &survey.Select{Message: "The AWS region that you are deploying to", Default: region, Options: AvailableRegions}
		err := survey.AskOne(p2, &region)
		if err != nil {
			return nil, err
		}
	}
	err := checkRegion(region)
	if err != nil {
		return nil, err
	}

	version, err := getVersion(c, region)
	if err != nil {
		return nil, err
	}

	// set up stack parameters
	cognitoPrefix := c.String("cognito-domain-prefix")
	if cognitoPrefix == "" {
		var company string
		p := &survey.Input{Message: "Your company name"}
		err := survey.AskOne(p, &company)
		if err != nil {
			return nil, err
		}
		// turn the company name into lowercase and remove spaces, so that it will work
		// as a Cognito prefix domain
		company = strings.ReplaceAll(company, " ", "")
		company = strings.ToLower(company)

		cognitoPrefix = fmt.Sprintf("granted-login-%s", company)
		p = &survey.Input{Message: "The prefix for the Cognito Sign in URL", Default: cognitoPrefix}
		err = survey.AskOne(p, &cognitoPrefix)
		if err != nil {
			return nil, err
		}
	}

	cfg := Config{
		Version: 1,
		Deployment: Deployment{
			StackName: name,
			Account:   account,
			Region:    region,
			Release:   version,
			Parameters: Parameters{
				CognitoDomainPrefix:  cognitoPrefix,
				AdministratorGroupID: "granted_administrators",
			},
		},
	}

	return &cfg, nil
}

func getVersion(c *cli.Context, region string) (string, error) {
	version := c.String("version")
	if version != "" {
		return version, nil
	}

	si := spinner.New(spinner.CharSets[14], 100*time.Millisecond)
	si.Suffix = " fetching the latest deployment version"
	si.Writer = os.Stderr
	si.Start()

	m, err := GetManifest(c.Context, region)
	si.Stop()
	if err == nil {
		clio.Info("using deployment version %s", m.LatestDeploymentVersion)
		return m.LatestDeploymentVersion, nil
	}

	// we couldn't fetch the manifest for some reason, so allow the user to enter a version manually.
	if err != nil {
		clio.Error(`error fetching manifest: %s.
You can try and enter a deployment version manually now, but there's no guarantees we'll be able to deploy it.
`, err)
	}
	p := &survey.Input{Message: "The version of Granted Approvals to deploy"}
	err = survey.AskOne(p, &version, survey.WithValidator(survey.MinLength(1)))
	return version, err
}

// tryGetCurrentAccountID uses AWS STS to try and load the current account ID.
func tryGetCurrentAccountID(ctx context.Context) (string, error) {
	si := spinner.New(spinner.CharSets[14], 100*time.Millisecond)
	si.Suffix = " loading AWS account ID from your current profile"
	si.Writer = os.Stderr
	si.Start()
	defer si.Stop()

	cfg, err := config.LoadDefaultConfig(ctx)
	if err != nil {
		zap.S().Debugw("failed fetching config", zap.Error(err))
		return "", err
	}
	client := sts.NewFromConfig(cfg)
	res, err := client.GetCallerIdentity(ctx, &sts.GetCallerIdentityInput{})
	if err != nil {
		zap.S().Debugw("failed getting caller identity", zap.Error(err))
		return "", err
	}
	if res.Account == nil {
		return "", nil
	}
	return *res.Account, nil
}

type GetCurrentAccountIDOpts struct {
	WarnExpiryIfWithinDuration *time.Duration
}

func WithWarnExpiryIfWithinDuration(t time.Duration) func(*GetCurrentAccountIDOpts) {
	return func(gcai *GetCurrentAccountIDOpts) {
		gcai.WarnExpiryIfWithinDuration = &t
	}
}

// MustGetCurrentAccountID uses AWS STS to try and load the current account ID.
//
// if not credentials are available, logs and error and os.Exit(1)
func MustGetCurrentAccountID(ctx context.Context, opts ...func(*GetCurrentAccountIDOpts)) string {
	var o GetCurrentAccountIDOpts
	for _, opt := range opts {
		opt(&o)
	}
	si := spinner.New(spinner.CharSets[14], 100*time.Millisecond)
	si.Suffix = " loading AWS account ID from your current profile"
	si.Writer = os.Stderr
	si.Start()
	defer si.Stop()

	cfg, err := config.LoadDefaultConfig(ctx)
	if err != nil {
		si.Stop()
		clio.Debug("Encountered error while loading default aws config: %s", err)
		clio.Error("Failed to load AWS credentials.")
		os.Exit(1)
	}

	creds, err := cfg.Credentials.Retrieve(ctx)
	if err != nil {
		si.Stop()
		clio.Debug("Encountered error while loading default aws config: %s", err)
		clio.Error("Failed to load AWS credentials.")
		os.Exit(1)
	}

	if !creds.HasKeys() {
		si.Stop()
		clio.Error("Could not find AWS credentials. Please export valid AWS credentials to run this command.")
		os.Exit(1)
	}

	if creds.Expired() {
		si.Stop()
		clio.Error("AWS credentials are expired. Please export valid AWS credentials to run this command.")
		os.Exit(1)
	}

	if o.WarnExpiryIfWithinDuration != nil && creds.CanExpire && creds.Expires.Before(time.Now().Add(*o.WarnExpiryIfWithinDuration)) {
		clio.Warn("AWS credentials expire in less than %s, consider exporting fresh credentials to avoid issues.", o.WarnExpiryIfWithinDuration.String())
	}

	client := sts.NewFromConfig(cfg)
	res, err := client.GetCallerIdentity(ctx, &sts.GetCallerIdentityInput{})
	if err != nil {
		si.Stop()
		clio.Debug("Encountered error while getting caller identity: %s", err)
		clio.Error("Failed to get AWS caller identity. Check that you have exported credentials and that they are not expired.")

		os.Exit(1)
	}
	if res.Account == nil {
		si.Stop()
		clio.Debug("Encountered nil response getting caller identity: %s", err)
		clio.Error("Failed to load AWS credentials.")
		os.Exit(1)
	}
	return *res.Account
}

// getDefaultAvailableRegion tries to match the AWS_REGION env var with one of the
// available regions for a Granted Approvals deployment.
// If the AWS_REGION env var doesn't match any of our available regions, the first
// AvailableRegion is returned instead.
func getDefaultAvailableRegion() string {
	region := os.Getenv("AWS_REGION")
	for _, r := range AvailableRegions {
		if r == region {
			return r
		}
	}
	return AvailableRegions[0]
}

func checkRegion(r string) error {
	for _, ar := range AvailableRegions {
		if r == ar {
			return nil
		}
	}
	return fmt.Errorf("we don't yet support deployments hosted in %s. Our supported regions are: [%s]", r, strings.Join(AvailableRegions, ", "))
}
