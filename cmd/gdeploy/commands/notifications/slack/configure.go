package slack

import (
	"bytes"
	"embed"
	"encoding/json"
	"fmt"
	"net/url"
	"strings"
	"text/template"

	"github.com/AlecAivazis/survey/v2"
	"github.com/common-fate/granted-approvals/pkg/clio"
	"github.com/common-fate/granted-approvals/pkg/config"
	"github.com/common-fate/granted-approvals/pkg/deploy"
	"github.com/pkg/errors"
	"github.com/urfave/cli/v2"
)

//go:embed templates
var templateFiles embed.FS

var configureSlackCommand = cli.Command{
	Name:        "configure",
	Description: "configure and enable slack integration",
	Action: func(c *cli.Context) error {
		ctx := c.Context
		f := c.Path("file")

		dc := deploy.MustLoadConfig(f)
		o, err := dc.LoadOutput(ctx)
		if err != nil {
			return err
		}
		apiUrl := o.APIURL

		appManifest, err := RenderSlackAppManifest(SlackManifestConfig{WebhookURL: strings.TrimSuffix(apiUrl, "/") + "/webhook/v1/slack/interactivity"})
		if err != nil {
			return err
		}

		appInstallURL := fmt.Sprintf("https://api.slack.com/apps?new_app=1&manifest_json=%s", url.QueryEscape(appManifest))
		clio.Info("Copy & paste the following link into your web browser to create a new Slack app for Granted Approvals:")
		fmt.Printf("\n\n%s\n\n", appInstallURL)
		clio.Info("After creating the app, install it to your workspace and find your Bot User OAuth Token in the OAuth & Permissions tab.")
		p := &survey.Password{
			Message: "Enter the Bot User OAuth Token for your Slack app",
		}
		var botUserToken string
		err = survey.AskOne(p, &botUserToken)
		if err != nil {
			return err
		}

		suffixedPath, version, err := config.PutSecretVersion(ctx, config.SlackTokenPath, dc.Deployment.Parameters.DeploymentSuffix, botUserToken)
		if err != nil {
			return errors.Wrap(err, "failed while setting Slack parameters in ssm")
		}

		dc.Notifications = &deploy.NotificationsConfig{
			Slack: &deploy.SlackConfig{
				APIToken: config.AWSSSMParamToken(suffixedPath, version),
			},
		}
		err = dc.Save(f)
		if err != nil {
			return err
		}
		clio.Warn("Your changes won't be applied until you redeploy. Run 'gdeploy update' to apply the changes to your CloudFormation deployment.")
		clio.Success("Successfully enabled Slack")

		return nil
	},
}

type SlackManifestConfig struct {
	WebhookURL string
}

func RenderSlackAppManifest(s SlackManifestConfig) (string, error) {
	tmpl, err := template.ParseFS(templateFiles, "templates/*")
	if err != nil {
		return "", err
	}

	buf := new(bytes.Buffer)
	err = tmpl.ExecuteTemplate(buf, "slack-app-maniftest.json.tmpl", s)
	if err != nil {
		return "", err
	}
	minibuf := new(bytes.Buffer)
	// compact removes whitespace from the json string
	// this allows much nicer escaped URLS
	err = json.Compact(minibuf, buf.Bytes())
	if err != nil {
		return "", err
	}
	return minibuf.String(), nil
}
