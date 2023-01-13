package deploy

import (
	"encoding/json"
	"fmt"
	"strings"
)

// CDKContextArgs returns the CDK context arguments
// in the form "-c" "ArgName=ArgValue"
//
// This should only be used in development, where the StackName variable is always of
// the form "common-fate-$STAGE". It panics if this is not the case.
func (c Config) CDKContextArgs() []string {
	name, err := c.GetDevStageName()
	if err != nil {
		panic(err)
	}

	name = CleanName(name)
	var args []string
	// pass context variables through as CLI arguments. This will eventually allow them to be
	// overridden in automated deployment workflows like in CI pipelines.
	args = append(args, "-c", fmt.Sprintf("stage=%s", name))
	args = append(args, "-c", fmt.Sprintf("cognitoDomainPrefix=cf-granted-%s", name))

	if c.Deployment.Parameters.ProviderConfiguration != nil {
		cfg, err := json.Marshal(c.Deployment.Parameters.ProviderConfiguration)
		if err != nil {
			panic(err)
		}

		args = append(args, "-c", fmt.Sprintf("providerConfiguration=%s", string(cfg)))
	}
	if c.Deployment.Parameters.IdentityConfiguration != nil {
		cfg, err := json.Marshal(c.Deployment.Parameters.IdentityConfiguration)
		if err != nil {
			panic(err)
		}

		args = append(args, "-c", fmt.Sprintf("identityConfiguration=%s", string(cfg)))
	}

	if c.Deployment.Parameters.NotificationsConfiguration != nil {
		cfg, err := json.Marshal(c.Deployment.Parameters.NotificationsConfiguration)
		if err != nil {
			panic(err)
		}
		args = append(args, "-c", fmt.Sprintf("notificationsConfiguration=%s", string(cfg)))
	}

	if c.Deployment.Parameters.IdentityProviderType != "" {
		args = append(args, "-c", fmt.Sprintf("idpType=%s", string(c.Deployment.Parameters.IdentityProviderType)))
	}
	if c.Deployment.Parameters.AdministratorGroupID != "" {
		args = append(args, "-c", fmt.Sprintf("adminGroupId=%s", string(c.Deployment.Parameters.AdministratorGroupID)))
	}
	if c.Deployment.Parameters.SamlSSOMetadata != "" {
		args = append(args, "-c", fmt.Sprintf("samlMetadata=%s", string(c.Deployment.Parameters.SamlSSOMetadata)))
	}
	if c.Deployment.Parameters.SamlSSOMetadataURL != "" {
		args = append(args, "-c", fmt.Sprintf("samlMetadataUrl=%s", string(c.Deployment.Parameters.SamlSSOMetadataURL)))
	}
	if c.Deployment.Parameters.CloudfrontWAFACLARN != "" {
		args = append(args, "-c", fmt.Sprintf("cloudfrontWafAclArn=%s", string(c.Deployment.Parameters.CloudfrontWAFACLARN)))
	}
	if c.Deployment.Parameters.APIGatewayWAFACLARN != "" {
		args = append(args, "-c", fmt.Sprintf("apiGatewayWafAclArn=%s", string(c.Deployment.Parameters.APIGatewayWAFACLARN)))
	}
	if c.Deployment.Parameters.ExperimentalRemoteConfigURL != "" {
		args = append(args, "-c", fmt.Sprintf("experimentalRemoteConfigUrl=%s", string(c.Deployment.Parameters.ExperimentalRemoteConfigURL)))
	}
	if c.Deployment.Parameters.ExperimentalRemoteConfigHeaders != "" {
		args = append(args, "-c", fmt.Sprintf("experimentalRemoteConfigHeaders=%s", string(c.Deployment.Parameters.ExperimentalRemoteConfigHeaders)))
	}
	if c.Deployment.Parameters.AnalyticsDisabled != "" {
		args = append(args, "-c", fmt.Sprintf("analyticsDisabled=%s", string(c.Deployment.Parameters.AnalyticsDisabled)))
	}
	if c.Deployment.Parameters.AnalyticsLogLevel != "" {
		args = append(args, "-c", fmt.Sprintf("analyticsLogLevel=%s", string(c.Deployment.Parameters.AnalyticsLogLevel)))
	}
	if c.Deployment.Parameters.SDKAPIJWTAudience != "" {
		args = append(args, "-c", fmt.Sprintf("sdkApiJwtAudience=%s", string(c.Deployment.Parameters.SDKAPIJWTAudience)))
	}
	if c.Deployment.Parameters.SDKAPIJWTIssuer != "" {
		args = append(args, "-c", fmt.Sprintf("sdkApiJwtIssuer=%s", string(c.Deployment.Parameters.SDKAPIJWTIssuer)))
	}

	// CDK deploys always use the dev analytics endpoint and debug mode
	args = append(args, "-c", "analyticsUrl=https://t-dev.commonfate.io")
	args = append(args, "-c", "analyticsDeploymentStage=dev")

	return args
}

// GetDevStageName returns the stage name to be used in a CDK deployment.
// It expects that the stack name is in the form "common-fate--$STAGE".
func (c Config) GetDevStageName() (string, error) {
	pre := "common-fate-"
	if !strings.HasPrefix(c.Deployment.StackName, pre) {
		return "", fmt.Errorf("stack name %s must start with %s for development", c.Deployment.StackName, pre)
	}
	return strings.TrimPrefix(c.Deployment.StackName, pre), nil
}
