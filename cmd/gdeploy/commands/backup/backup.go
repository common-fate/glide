package backup

import (
	"errors"
	"fmt"
	"regexp"

	"github.com/AlecAivazis/survey/v2"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/common-fate/clio"
	"github.com/common-fate/common-fate/pkg/deploy"
	"github.com/urfave/cli/v2"
)

var Command = cli.Command{
	Name:        "backup",
	Description: "Backup the DynamoDB table",
	Usage:       "Backup the DynamoDB table",
	Flags: []cli.Flag{
		&cli.BoolFlag{Name: "confirm", Aliases: []string{"y"}, Usage: "If provided, will automatically continue without asking for confirmation"},
	},
	Subcommands: []*cli.Command{&BackupStatus},
	Action: func(c *cli.Context) error {
		ctx := c.Context
		dc, err := deploy.ConfigFromContext(ctx)
		if err != nil {
			return err
		}

		stackOutput, err := dc.LoadOutput(ctx)
		if err != nil {
			return err
		}

		p := &survey.Input{
			Message: "Enter a backup name",
		}
		var backupName string
		err = survey.AskOne(p, &backupName, survey.WithValidator(func(ans interface{}) error {
			a := ans.(string)
			r := regexp.MustCompile(`[^a-zA-Z0-9_.-]`)
			match := r.MatchString(a)
			if match {
				return fmt.Errorf("value: `%s` must satisfy regular expression pattern: [a-zA-Z0-9_.-]+", a)
			}
			return nil
		}))
		if err != nil {
			return err
		}

		clio.Infof("Creating backup of Common Fate dynamoDB table: %s", stackOutput.DynamoDBTable)
		confirm := c.Bool("confirm")
		if !confirm {
			cp := &survey.Confirm{Message: "Do you wish to continue?", Default: true}
			err = survey.AskOne(cp, &confirm)
			if err != nil {
				return err
			}
		}

		if !confirm {
			return errors.New("user cancelled backup")
		}
		backupOutput, err := deploy.StartBackup(ctx, stackOutput.DynamoDBTable, backupName)
		if err != nil {
			return err
		}
		clio.Successf("Successfully started a backup of Common Fate dynamoDB table: %s", stackOutput.DynamoDBTable)
		clio.Infof("Backup details\n%s", deploy.BackupDetailsToString(backupOutput))
		clio.Infof("To view the status of this backup, run `gdeploy backup status --arn=%s`", aws.ToString(backupOutput.BackupArn))
		clio.Infof("To restore from this backup, run `gdeploy restore --arn=%s`", aws.ToString(backupOutput.BackupArn))

		return nil
	},
}
