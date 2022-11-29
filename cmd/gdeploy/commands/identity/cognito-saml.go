package identity

// Command to generate a SAML Cognito url from gdeploy outputs

import (
	"fmt"

	"github.com/common-fate/common-fate/pkg/deploy"
	"github.com/urfave/cli/v2"
)

var CognitoSamlCommand = cli.Command{
	Name:        "cognito-saml",
	Description: "Generate a SAML Cognito url that can be used to login to Common Fate",
	Action: cli.ActionFunc(func(c *cli.Context) error {

		ctx := c.Context

		dc, err := deploy.ConfigFromContext(ctx)
		if err != nil {
			return err
		}
		o, err := dc.LoadOutput(ctx)
		if err != nil {
			return err
		}
		UserPoolDomain := o.UserPoolDomain
		SAMLIdentityProviderName := o.SAMLIdentityProviderName
		CognitoClientID := o.CognitoClientID
		FrontendDomainOutput := o.FrontendDomainOutput

		// Check if any are nil
		if UserPoolDomain == "" || SAMLIdentityProviderName == "" || CognitoClientID == "" || FrontendDomainOutput == "" {
			return fmt.Errorf("missing required output values")
		}

		// Url format:
		// https://$(UserPoolDomain)/authorize?response_type=code&identity_provider=$(SAMLIdentityProviderName)&client_id=$(CognitoClientID)&redirect_uri=https://$(FrontendDomainOutput)

		url := fmt.Sprintf("https://%s/authorize?response_type=code&identity_provider=%s&client_id=%s&redirect_uri=https://%s", UserPoolDomain, SAMLIdentityProviderName, CognitoClientID, FrontendDomainOutput)

		fmt.Println(url)

		return nil
	}),
}
