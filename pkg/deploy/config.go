package deploy

import (
	"context"
	"encoding/json"

	"fmt"
	"os"
	"strings"
	"time"

	"github.com/AlecAivazis/survey/v2"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/cloudformation/types"
	"github.com/aws/aws-sdk-go-v2/service/sts"
	"github.com/briandowns/spinner"
	"github.com/common-fate/clio"
	"github.com/common-fate/clio/clierr"
	"github.com/common-fate/common-fate/pkg/cfaws"
	"github.com/common-fate/common-fate/pkg/gconfig"
	"github.com/fatih/color"
	"github.com/pkg/errors"
	"github.com/urfave/cli/v2"
	"gopkg.in/yaml.v3"
)

type contextkey struct{}

var DeploymentConfigContextKey contextkey

func ConfigFromContext(ctx context.Context) (Config, error) {
	if cfg := ctx.Value(DeploymentConfigContextKey); cfg != nil {
		return cfg.(Config), nil
	} else {
		return Config{}, ErrConfigNotNotSetInContext
	}
}

func SetConfigInContext(ctx context.Context, cfg Config) context.Context {
	return context.WithValue(ctx, DeploymentConfigContextKey, cfg)
}

const DeprecatedDefaultFilename = "granted-deployment.yml"
const DefaultFilename = "deployment.yml"

var DeprecatedDefaultFilenameWarning = clierr.Warn("Since v0.11.0 the default deployment config file has been renamed from 'granted-deployment.yml' to 'deployment.yml'. To update, rename the file now or run this command to rename via the cli `mv granted-deployment.yml deployment.yml`")

const DefaultCommonFateAdministratorsGroup = "common_fate_administrators"

// AvailableRegions are the regions that we currently release CloudFormation templates to.
var AvailableRegions = []string{
	"ap-southeast-2",
	"us-west-2",
	"us-east-1",
	"eu-central-1",
}

type Config struct {
	Version      int        `yaml:"version"`
	Deployment   Deployment `yaml:"deployment"`
	cachedOutput *Output
}

type Deployment struct {
	StackName string `yaml:"stackName"`
	Account   string `yaml:"account"`
	Region    string `yaml:"region"`
	// Release may be one of two formats:
	//
	// 1. A release version tag (e.g. 'v0.1.0'). This uses a release
	// from Common Fate's release bucket.
	//
	// 2. A path to a CloudFormation template in S3, in the format
	// 'https://custom-bucket.s3.amazonaws.com/path/to/template.json'.
	// Note that the S3 bucket must be in the same region as the 'Region' parameter.
	Release string `yaml:"release"`
	// Dev is set to true for internal development deployments only.
	Dev        *bool             `yaml:"dev,omitempty"`
	Parameters Parameters        `yaml:"parameters"`
	Tags       map[string]string `yaml:"tags,omitempty"`
}

type Notifications struct {
	Slack                 map[string]string `yaml:"slack,omitempty" json:"slack,omitempty"`
	SlackIncomingWebhooks FeatureMap        `yaml:"slackIncomingWebhooks,omitempty" json:"slackIncomingWebhooks,omitempty"`
}

// Feature map represents the type used for features like identity and notifications
type FeatureMap map[string]map[string]string

// Upserts the feature in the map, if the map is not initialised, it initialises it first
func (f *FeatureMap) Upsert(id string, feature map[string]string) {
	// check if this is a nil map and initialise first if so
	// This is a trick to check the underlying maps from the alias' value
	if map[string]map[string]string(*f) == nil {
		*f = make(map[string]map[string]string)
	}
	(*f)[id] = feature
}

// Remove the feature in the map, if the map is not initialised, it does nothing
func (f FeatureMap) Remove(id string) {
	// check if this is a nil map and initialise first if so
	// This is a trick to check the underlying maps from the alias' value
	if map[string]map[string]string(f) == nil {
		return
	}
	delete(f, id)
}

type Parameters struct {
	CognitoDomainPrefix             string         `yaml:"CognitoDomainPrefix"`
	AdministratorGroupID            string         `yaml:"AdministratorGroupID"`
	DeploymentSuffix                string         `yaml:"DeploymentSuffix,omitempty"`
	IdentityProviderType            string         `yaml:"IdentityProviderType,omitempty"`
	SamlSSOMetadata                 string         `yaml:"SamlSSOMetadata,omitempty"`
	SamlSSOMetadataURL              string         `yaml:"SamlSSOMetadataURL,omitempty"`
	FrontendDomain                  string         `yaml:"FrontendDomain,omitempty"`
	FrontendCertificateARN          string         `yaml:"FrontendCertificateARN,omitempty"`
	CloudfrontWAFACLARN             string         `yaml:"CloudfrontWAFACLARN,omitempty"`
	APIGatewayWAFACLARN             string         `yaml:"APIGatewayWAFACLARN,omitempty"`
	ExperimentalRemoteConfigURL     string         `yaml:"ExperimentalRemoteConfigURL,omitempty"`
	ExperimentalRemoteConfigHeaders string         `yaml:"ExperimentalRemoteConfigHeaders,omitempty"`
	IdentityConfiguration           FeatureMap     `yaml:"IdentityConfiguration,omitempty"`
	NotificationsConfiguration      *Notifications `yaml:"NotificationsConfiguration,omitempty"`
	AnalyticsDisabled               string         `yaml:"AnalyticsDisabled,omitempty"`
	AnalyticsURL                    string         `yaml:"AnalyticsURL,omitempty"`
	AnalyticsLogLevel               string         `yaml:"AnalyticsLogLevel,omitempty"`
	AnalyticsDeploymentStage        string         `yaml:"AnalyticsDeploymentStage,omitempty"`
	IdentityGroupFilter             string         `yaml:"IdentityGroupFilter,omitempty"`
	EnableCronHealthCheckInDev      string         `yaml:"EnableCronHealthCheckInDev,omitempty"`
	IDPSyncTimeoutSeconds           string         `yaml:"IDPSyncTimeoutSeconds,omitempty"`
	IDPSyncSchedule                 string         `yaml:"IDPSyncSchedule,omitempty"`
	IDPSyncMemory                   string         `yaml:"IDPSyncMemory,omitempty"`
}

// UnmarshalFeatureMap parses the JSON configuration data and returns
// an initialised FeatureMap. If `data` is an empty string an empty
// FeatureMap is returned.
func UnmarshalFeatureMap(data string) (FeatureMap, error) {
	if data == "" {
		return make(FeatureMap), nil
	}
	// first remove any double backslashes which may have been added while loading from or to environment
	// the process of loading escaped strings into the environment can sometimes add double escapes which cannot be parsed correctly
	// unless removed
	data = strings.ReplaceAll(string(data), "\\", "")
	var i FeatureMap
	err := json.Unmarshal([]byte(data), &i)
	if err != nil {
		return nil, err
	}
	return i, nil
}

// UnmarshalNotifications parses the JSON configuration data and returns
// an initialised Notifications. If `data` is an empty string an empty
// Notifications is returned.
func UnmarshalNotifications(data string) (*Notifications, error) {
	if data == "" {
		return &Notifications{}, nil
	}
	// first remove any double backslashes which may have been added while loading from or to environment
	// the process of loading escaped strings into the environment can sometimes add double escapes which cannot be parsed correctly
	// unless removed
	data = strings.ReplaceAll(string(data), "\\", "")
	var i Notifications
	err := json.Unmarshal([]byte(data), &i)
	if err != nil {
		return nil, err
	}
	return &i, nil
}

// RunConfigTest runs ConfigTest() if it is implemented on the interface
func RunConfigTest(ctx context.Context, testable interface{}) error {
	if tester, ok := testable.(gconfig.Tester); ok {
		clio.Info("running tests using this configuration...")

		if initer, ok := testable.(gconfig.Initer); ok {
			err := initer.Init(ctx)
			if err != nil {
				return err
			}
		}
		err := tester.TestConfig(ctx)
		if err != nil {
			return err
		}
		clio.Success("Configuration tests passed!")
	}
	return nil
}

// Reset Identity Provider to cognito settings
func (c *Config) ResetIdentityProviderToCognito(filepath string) error {

	c.Deployment.Parameters.IdentityProviderType = ""
	c.Deployment.Parameters.AdministratorGroupID = "common_fate_administrators"
	c.Deployment.Parameters.IdentityConfiguration = nil
	c.Deployment.Parameters.SamlSSOMetadataURL = ""
	c.Deployment.Parameters.SamlSSOMetadata = ""

	if err := c.Save(filepath); err != nil {
		return err
	}

	return nil
}

// CLIPrompt prompts the user to enter a value for the config varsiable
// in a CLI context. If the config variable implements Defaulter, the
// default value is returned and the user is not prompted for any input.
func CLIPrompt(f *gconfig.Field) error {
	grey := color.New(color.FgHiBlack)
	msg := f.Key()
	if f.Description() != "" {
		msg = f.Description() + " " + grey.Sprintf("(%s)", msg)
	}

	// @TODO work out how to integrate the optional prompt here
	// you shoudl be able to choose to set or unset
	// if you choose to set, it should use a default if it exists
	// By design, we can't have an optional secret as they are mutually exclusive.
	var p survey.Prompt
	if f.IsSecret() && f.Get() != "" {
		confMsg := msg + " would you like to update this secret?"
		p = &survey.Confirm{
			Message: confMsg,
		}
		var doUpdate bool
		err := survey.AskOne(p, &doUpdate)
		if err != nil {
			return err
		}
		if !doUpdate {

			return nil
		}

	}

	//Handle different methods of cli prompt inputs.
	var val string
	switch f.CLIPrompt() {
	case gconfig.CLIPromptTypeString:
		defaultValue := f.Get()
		if defaultValue == "" {
			defaultValue = f.Default()
		}
		p = &survey.Input{
			Message: msg,
			Default: defaultValue,
		}

		err := survey.AskOne(p, &val)
		if err != nil {
			return err
		}

	case gconfig.CLIPromptTypePassword:
		p = &survey.Password{
			Message: msg,
		}
		err := survey.AskOne(p, &val)
		if err != nil {
			return err
		}

	case gconfig.CLIPromptTypeFile:
		p5 := &survey.Input{Message: "the file path to " + msg}
		var res string
		err := survey.AskOne(p5, &res, func(options *survey.AskOptions) error {
			options.Validators = append(options.Validators, func(ans interface{}) error {
				p := ans.(string)
				fileInfo, err := os.Stat(p)
				if err != nil {
					return err
				}
				if fileInfo.IsDir() {
					return fmt.Errorf("path is a directory, must be a file")
				}
				return nil
			})
			return nil
		})
		if err != nil {
			return err
		}
		b, err := os.ReadFile(res)
		if err != nil {
			return err
		}
		val = string(b)

	}

	// set the value.
	return f.Set(val)
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

	if c.Deployment.Parameters.NotificationsConfiguration != nil {
		if c.Deployment.Parameters.NotificationsConfiguration != nil {
			config, err := json.Marshal(c.Deployment.Parameters.NotificationsConfiguration)
			if err != nil {
				return nil, err
			}
			configStr := string(config)
			res = append(res, types.Parameter{
				ParameterKey:   aws.String("NotificationsConfiguration"),
				ParameterValue: &configStr,
			})
		}

	}
	if c.Deployment.Parameters.IdentityConfiguration != nil {
		config, err := json.Marshal(c.Deployment.Parameters.IdentityConfiguration)
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
	if c.Deployment.Parameters.APIGatewayWAFACLARN != "" {
		res = append(res, types.Parameter{
			ParameterKey:   aws.String("APIGatewayWAFACLARN"),
			ParameterValue: &p.APIGatewayWAFACLARN,
		})
	}
	if c.Deployment.Parameters.CloudfrontWAFACLARN != "" {
		res = append(res, types.Parameter{
			ParameterKey:   aws.String("CloudfrontWAFACLARN"),
			ParameterValue: &p.CloudfrontWAFACLARN,
		})
	}

	if c.Deployment.Parameters.ExperimentalRemoteConfigURL != "" {
		res = append(res, types.Parameter{
			ParameterKey:   aws.String("ExperimentalRemoteConfigURL"),
			ParameterValue: &p.ExperimentalRemoteConfigURL,
		})
	}
	if c.Deployment.Parameters.ExperimentalRemoteConfigHeaders != "" {
		res = append(res, types.Parameter{
			ParameterKey:   aws.String("ExperimentalRemoteConfigHeaders"),
			ParameterValue: &p.ExperimentalRemoteConfigHeaders,
		})
	}

	if c.Deployment.Parameters.AnalyticsDisabled != "" {
		res = append(res, types.Parameter{
			ParameterKey:   aws.String("AnalyticsDisabled"),
			ParameterValue: &p.AnalyticsDisabled,
		})
	}
	if c.Deployment.Parameters.AnalyticsURL != "" {
		res = append(res, types.Parameter{
			ParameterKey:   aws.String("AnalyticsURL"),
			ParameterValue: &p.AnalyticsURL,
		})
	}
	if c.Deployment.Parameters.AnalyticsLogLevel != "" {
		res = append(res, types.Parameter{
			ParameterKey:   aws.String("AnalyticsLogLevel"),
			ParameterValue: &p.AnalyticsLogLevel,
		})
	}
	if c.Deployment.Parameters.AnalyticsDeploymentStage != "" {
		res = append(res, types.Parameter{
			ParameterKey:   aws.String("AnalyticsDeploymentStage"),
			ParameterValue: &p.AnalyticsDeploymentStage,
		})
	}
	if c.Deployment.Parameters.IDPSyncMemory != "" {
		res = append(res, types.Parameter{
			ParameterKey:   aws.String("IDPSyncMemory"),
			ParameterValue: &p.IDPSyncMemory,
		})
	}
	if c.Deployment.Parameters.IDPSyncSchedule != "" {
		res = append(res, types.Parameter{
			ParameterKey:   aws.String("IDPSyncSchedule"),
			ParameterValue: &p.IDPSyncSchedule,
		})
	}
	if c.Deployment.Parameters.IDPSyncTimeoutSeconds != "" {
		res = append(res, types.Parameter{
			ParameterKey:   aws.String("IDPSyncTimeoutSeconds"),
			ParameterValue: &p.IDPSyncTimeoutSeconds,
		})
	}

	return res, nil
}

// LoadConfig attempts to load the config file at path f
// if it does not exist, returns ErrConfigNotExist
// else returns the config or any other error
//
// in CLI commands, it is preferable to use deploy.ConfigFromContext(ctx) where gdeploy.RequireDeploymentConfig has run as a before function for the command
// gdeploy.RequireDeploymentConfig will return a helpful cli error if there are any issues
func LoadConfig(f string) (Config, error) {
	if _, err := os.Stat(f); errors.Is(err, os.ErrNotExist) {
		return Config{}, ErrConfigNotExist
	}

	fileRead, err := os.OpenFile(f, os.O_RDWR|os.O_CREATE, 0600)
	if err != nil {
		return Config{}, err
	}
	defer fileRead.Close()
	var dc Config
	decoder := yaml.NewDecoder(fileRead)
	decoder.KnownFields(true)
	err = decoder.Decode(&dc)
	if err != nil {
		return Config{}, err
	}
	return dc, nil
}

// CfnTemplateURL returns the CloudFormation template URL.
// If the deployment release points to an S3 object (https://custom-bucket.s3.amazonaws.com/path/to/template.json)
// It is turned into a HTTPS URL. If a regular version number (v0.1.0) is used, we point to our official release bucket.
func (c *Config) CfnTemplateURL() string {
	// use a custom URL if it was provided
	if strings.HasPrefix(c.Deployment.Release, "https://") {
		return c.Deployment.Release
	}

	// otherwise, use the Common Fate release bucket
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
	acc, _ := TryGetCurrentAccountID(ctx)
	dev := true
	conf := Config{
		Deployment: Deployment{
			StackName: "common-fate-" + stage,
			Account:   acc,
			Dev:       &dev,

			Parameters: Parameters{
				AdministratorGroupID: "granted_administrators",
			},
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
	account, _ := TryGetCurrentAccountID(ctx)

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
		Version: 2,
		Deployment: Deployment{
			StackName: fmt.Sprintf("common-fate-%s", stage),
			Account:   account,
			Region:    region,
			Dev:       &dev,
			Parameters: Parameters{
				AdministratorGroupID: "common_fate_administrators",
				AnalyticsDisabled:    "true",
			},
		},
	}

	return &c, nil
}

// SetupReleaseConfig sets up the release configuration used in production deployments.
func SetupReleaseConfig(c *cli.Context) (*Config, error) {
	name := c.String("name")
	if name == "" {
		name = "CommonFate"
		p := &survey.Input{Message: "The name of the CloudFormation stack to create", Default: name}
		err := survey.AskOne(p, &name, survey.WithValidator(survey.MinLength(1)))
		if err != nil {
			return nil, err
		}
	}

	account := c.String("account")
	if account == "" {
		ctx := context.Background()
		//Setting account to a fresh variable when set, was being unallocated when using the account variable
		current, err := TryGetCurrentAccountID(ctx)
		if err != nil {
			return nil, err
		}
		p := &survey.Input{Message: "The account ID that you are deploying to", Default: current}
		err = survey.AskOne(p, &account)
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

	checkRegion(region)

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

		cognitoPrefix = fmt.Sprintf("common-fate-login-%s", company)
		p = &survey.Input{Message: "The prefix for the Cognito Sign in URL", Default: cognitoPrefix}
		err = survey.AskOne(p, &cognitoPrefix)
		if err != nil {
			return nil, err
		}
	}

	cfg := Config{
		Version: 2,
		Deployment: Deployment{
			StackName: name,
			Account:   account,
			Region:    region,
			Release:   version,
			Parameters: Parameters{
				CognitoDomainPrefix:  cognitoPrefix,
				AdministratorGroupID: "common_fate_administrators",
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
		clio.Infof("using deployment version %s", m.LatestDeploymentVersion)
		return m.LatestDeploymentVersion, nil
	}

	// we couldn't fetch the manifest for some reason, so allow the user to enter a version manually.
	if err != nil {
		clio.Errorf(`error fetching manifest: %s.
You can try and enter a deployment version manually now, but there's no guarantees we'll be able to deploy it.
`, err)
	}
	p := &survey.Input{Message: "The version of Common Fate to deploy"}
	err = survey.AskOne(p, &version, survey.WithValidator(survey.MinLength(1)))
	return version, err
}

// TryGetCurrentAccountID uses AWS STS to try and load the current account ID.
func TryGetCurrentAccountID(ctx context.Context) (string, error) {
	si := spinner.New(spinner.CharSets[14], 100*time.Millisecond)
	si.Suffix = " loading AWS account ID from your current profile"
	si.Writer = os.Stderr
	si.Start()
	defer si.Stop()

	cfg, err := cfaws.ConfigFromContextOrDefault(ctx)
	if err != nil {
		return "", errors.Wrap(err, "failed fetching config")
	}
	client := sts.NewFromConfig(cfg)
	res, err := client.GetCallerIdentity(ctx, &sts.GetCallerIdentityInput{})
	if err != nil {
		return "", errors.Wrap(err, "failed getting caller identity")
	}
	if res.Account == nil {
		return "", nil
	}
	return *res.Account, nil
}

// getDefaultAvailableRegion tries to match the AWS_REGION env var with one of the
// available regions for a Common Fate deployment.
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

func checkRegion(r string) {
	for _, ar := range AvailableRegions {
		if r == ar {
			return
		}
	}
	// print a warning here
	clio.Warnf("we don't yet support deployments hosted in %s. Our supported regions are: [%s]", r, strings.Join(AvailableRegions, ", "))
}
