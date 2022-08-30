package commands

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"strings"

	"github.com/AlecAivazis/survey/v2"
	"github.com/aws/aws-sdk-go-v2/service/ssm"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/common-fate/granted-approvals/pkg/cfaws"
	"github.com/common-fate/granted-approvals/pkg/clio"
	"github.com/common-fate/granted-approvals/pkg/deploy"
	"github.com/common-fate/granted-approvals/pkg/identity/identitysync"
	slacknotifier "github.com/common-fate/granted-approvals/pkg/notifiers/slack"
	"github.com/urfave/cli/v2"
	"gopkg.in/yaml.v2"
)

var MigrateCommand = cli.Command{
	Name:        "migrate",
	Description: "migrate from config version 1 to config version 2",
	Usage:       "migrate from config version 1 to config version 2",
	Action: func(c *cli.Context) error {
		ctx := c.Context
		file := c.String("file")
		cfgv1, err := LoadConfigV1(file)
		if err != nil {
			return err
		}

		if cfgv1.Version != 1 {
			return fmt.Errorf("expected config to be version 1, found %d", cfgv1.Version)
		}
		clio.Info("detected config version 1")
		clio.Info("this command will update your config to match the new format")
		conf := survey.Confirm{Message: "Please confirm you want to migrate your config version from 1 -> 2"}
		var migrate bool
		err = survey.AskOne(&conf, &migrate)
		if err != nil {
			return err
		}
		if !migrate {
			clio.Info("cancelled migration")
			return nil
		}

		cfgv2 := deploy.Config{
			Version:    2,
			Deployment: cfgv1.Deployment,
		}
		// identity provider type
		cfgv2.Deployment.Parameters.IdentityProviderType = strings.ToLower(cfgv1.Deployment.Parameters.IdentityProviderType)

		// providers
		for k, v := range cfgv1.Providers {
			newWith := make(map[string]interface{})
			for k2, w := range v.With {
				// add a version to the provider ssm params
				if strings.HasPrefix(w, "awsssm://") {
					cfg, err := cfaws.ConfigFromContextOrDefault(ctx)
					if err != nil {
						return err
					}

					client := ssm.NewFromConfig(cfg)
					// get the lastest param version from ssm
					o, err := client.GetParameterHistory(ctx, &ssm.GetParameterHistoryInput{Name: aws.String(strings.TrimPrefix(w, "awsssm://"))})
					if err != nil {
						return err
					}
					version := int64(1)
					for _, p := range o.Parameters {
						if p.Version > int64(version) {
							version = p.Version
						}
					}
					// fix some bad casing from v1 config
					w += fmt.Sprintf(":%d", version)
				}

				newWith[strings.ReplaceAll(k2, "ID", "Id")] = w
			}
			err := cfgv2.Deployment.Parameters.ProviderConfiguration.Add(k, deploy.Provider{Uses: v.Uses, With: newWith})
			if err != nil {
				return err
			}
		}

		// Identity
		if cfgv1.Identity != nil {
			if cfgv1.Identity.Azure != nil {
				b, err := json.Marshal(cfgv1.Identity.Azure)
				if err != nil {
					return err
				}
				var newWith map[string]interface{}
				err = json.Unmarshal(b, &newWith)
				if err != nil {
					return err
				}
				cfgv2.Deployment.Parameters.IdentityConfiguration.Upsert(identitysync.IDPTypeAzureAD, newWith)
			}
			if cfgv1.Identity.Google != nil {
				b, err := json.Marshal(cfgv1.Identity.Google)
				if err != nil {
					return err
				}
				var newWith map[string]interface{}
				err = json.Unmarshal(b, &newWith)
				if err != nil {
					return err
				}
				cfgv2.Deployment.Parameters.IdentityConfiguration.Upsert(identitysync.IDPTypeGoogle, newWith)
			}
			if cfgv1.Identity.Okta != nil {
				b, err := json.Marshal(cfgv1.Identity.Okta)
				if err != nil {
					return err
				}
				var newWith map[string]interface{}
				err = json.Unmarshal(b, &newWith)
				if err != nil {
					return err
				}
				cfgv2.Deployment.Parameters.IdentityConfiguration.Upsert(identitysync.IDPTypeOkta, newWith)
			}
		}
		// Identity
		if cfgv1.Notifications != nil {
			if cfgv1.Notifications.Slack != nil {
				b, err := json.Marshal(cfgv1.Notifications.Slack)
				if err != nil {
					return err
				}
				var newWith map[string]interface{}
				err = json.Unmarshal(b, &newWith)
				if err != nil {
					return err
				}
				cfgv2.Deployment.Parameters.NotificationsConfiguration.Upsert(slacknotifier.NotificationsTypeSlack, newWith)
			}

		}
		err = cfgv2.Save(file)
		if err != nil {
			return err
		}
		clio.Success("Successfully migrated from config version 1 -> 2")
		return nil
	},
}

func LoadConfigV1(f string) (ConfigV1, error) {
	if _, err := os.Stat(f); errors.Is(err, os.ErrNotExist) {
		return ConfigV1{}, errors.New("config file does not exist")
	}

	fileRead, err := os.OpenFile(f, os.O_RDWR|os.O_CREATE, 0600)
	if err != nil {
		return ConfigV1{}, err
	}
	defer fileRead.Close()
	var dc ConfigV1
	err = yaml.NewDecoder(fileRead).Decode(&dc)
	if err != nil {
		return ConfigV1{}, err
	}
	return dc, nil
}

type ConfigV1 struct {
	Version       int                  `yaml:"version"`
	Deployment    deploy.Deployment    `yaml:"deployment"`
	Providers     map[string]Provider  `yaml:"providers,omitempty"`
	Notifications *NotificationsConfig `yaml:"notifications,omitempty"`
	Identity      *IdentityConfig      `yaml:"identity,omitempty"`
}
type Provider struct {
	Uses string            `yaml:"uses" json:"uses"`
	With map[string]string `yaml:"with" json:"with"`
}

type IdentityConfig struct {
	Google *Google `yaml:"google,omitempty" json:"google,omitempty"`
	Okta   *Okta   `yaml:"okta,omitempty" json:"okta,omitempty"`
	Azure  *Azure  `yaml:"azure,omitempty" json:"azure,omitempty"`
}

type Google struct {
	APIToken   string `yaml:"apiToken" json:"apiToken"`
	Domain     string `yaml:"domain" json:"domain"`
	AdminEmail string `yaml:"adminEmail" json:"adminEmail"`
}

type Okta struct {
	APIToken string `yaml:"apiToken" json:"apiToken"`
	OrgURL   string `yaml:"orgUrl" json:"orgUrl"`
}

type Azure struct {
	// V1 config had a casing issue with yaml
	TenantID     string `yaml:"tenantID" json:"tenantId"`
	ClientID     string `yaml:"clientID" json:"clientId"`
	ClientSecret string `yaml:"clientSecret" json:"clientSecret"`
}

type NotificationsConfig struct {
	Slack *SlackConfig `yaml:"slack,omitempty" json:"slack,omitempty"`
}

type SlackConfig struct {
	APIToken string `yaml:"apiToken" json:"apiToken"`
}
