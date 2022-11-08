package slack

import (
	"bytes"
	"embed"
	"encoding/json"
	"fmt"
	"net/url"
	"strings"
	"text/template"

	"github.com/common-fate/clio"
	"github.com/common-fate/granted-approvals/pkg/deploy"
	"github.com/common-fate/granted-approvals/pkg/gconfig"
	slacknotifier "github.com/common-fate/granted-approvals/pkg/notifiers/slack"
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

		dc, err := deploy.ConfigFromContext(ctx)
		if err != nil {
			return err
		}
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

		var slack slacknotifier.SlackDirectMessage
		if dc.Deployment.Parameters.NotificationsConfiguration == nil {
			dc.Deployment.Parameters.NotificationsConfiguration = &deploy.Notifications{}
		}
		cfg := slack.Config()
		if dc.Deployment.Parameters.NotificationsConfiguration.Slack != nil {
			err = cfg.Load(ctx, &gconfig.MapLoader{Values: dc.Deployment.Parameters.NotificationsConfiguration.Slack})
			if err != nil {
				return err
			}
		}

		for _, v := range cfg {
			err := deploy.CLIPrompt(v)
			if err != nil {
				return err
			}
		}

		err = deploy.RunConfigTest(ctx, &slack)
		if err != nil {
			return err
		}

		// if tests pass, dump the config and update in the deployment config
		newConfig, err := cfg.Dump(ctx, gconfig.SSMDumper{Suffix: dc.Deployment.Parameters.DeploymentSuffix})
		if err != nil {
			return err
		}

		dc.Deployment.Parameters.NotificationsConfiguration.Slack = newConfig

		err = dc.Save(f)
		if err != nil {
			return err
		}

		clio.Success("Successfully configured Slack")
		clio.Warn("Your changes won't be applied until you redeploy. Run 'gdeploy update' to apply the changes to your CloudFormation deployment.")
		clio.Warn("Run: `gdeploy notifications slack test --email=<your_slack_email>` to send a test DM")

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
