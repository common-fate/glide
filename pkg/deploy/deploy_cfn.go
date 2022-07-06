package deploy

import (
	"context"
	"fmt"
	"os"
	"time"

	"github.com/AlecAivazis/survey/v2"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/briandowns/spinner"
	"github.com/common-fate/cloudform/cfn"
	"github.com/common-fate/cloudform/console"
	"github.com/common-fate/cloudform/ui"
	"github.com/common-fate/granted-approvals/pkg/clio"
	"github.com/pkg/errors"
)

const noChangeFoundMsg = "The submitted information didn't contain changes. Submit different information to create a change set."

// DeployCloudFormation creates a CloudFormation stack based on the config
func (c *Config) DeployCloudFormation(ctx context.Context, confirm bool) error {
	template := c.CfnTemplateURL()
	cfg, err := config.LoadDefaultConfig(ctx, config.WithRegion(c.Deployment.Region))
	if err != nil {
		return err
	}
	cfnClient := cfn.New(cfg)

	si := spinner.New(spinner.CharSets[14], 100*time.Millisecond)
	si.Suffix = " creating CloudFormation change set"
	si.Writer = os.Stderr
	si.Start()

	params, err := c.CfnParams()
	if err != nil {
		return err
	}

	changeSetName, createErr := cfnClient.CreateChangeSet(ctx, template, params, nil, c.Deployment.StackName, "")

	si.Stop()

	if createErr != nil {
		if createErr.Error() == noChangeFoundMsg {
			clio.Success("Change set was created, but there is no change. Deploy was skipped.")
			return nil
		} else {
			return errors.Wrap(createErr, "creating changeset")
		}
	}
	uiClient := ui.New(cfg)
	if !confirm {
		status, err := uiClient.FormatChangeSet(ctx, c.Deployment.StackName, changeSetName)
		if err != nil {
			return err
		}
		clio.Info("The following CloudFormation changes will be made:")
		fmt.Println(status)

		p := &survey.Confirm{Message: "Do you wish to continue?", Default: true}
		err = survey.AskOne(p, &confirm)
		if err != nil {
			return err
		}
		if !confirm {
			return errors.New("user cancelled deployment")
		}
	}

	err = cfnClient.ExecuteChangeSet(ctx, c.Deployment.StackName, changeSetName)
	if err != nil {
		return err
	}

	status, messages := uiClient.WaitForStackToSettle(ctx, c.Deployment.StackName)

	fmt.Println("Final stack status:", ui.ColouriseStatus(status))

	if len(messages) > 0 {
		fmt.Println(console.Yellow("Messages:"))
		for _, message := range messages {
			fmt.Printf("  - %s\n", message)
		}
	}
	return nil
}
