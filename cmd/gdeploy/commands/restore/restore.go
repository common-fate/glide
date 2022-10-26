package restore

import (
	"fmt"
	"regexp"

	"github.com/AlecAivazis/survey/v2"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/common-fate/clio"
	"github.com/common-fate/granted-approvals/pkg/deploy"
	"github.com/urfave/cli/v2"
)

var Command = cli.Command{
	Name:        "restore",
	Description: "Restore a DynamoDB backup to a new table",
	Usage:       "Restore a DynamoDB backup to a new table",
	Flags: []cli.Flag{
		&cli.BoolFlag{Name: "confirm", Aliases: []string{"y"}, Usage: "If provided, will automatically continue without asking for confirmation"},
		&cli.StringFlag{Name: "arn", Usage: "The ARN of the backup to restore"},
		&cli.StringFlag{Name: "table-name", Usage: "The name of a new table to restore this backup to"},
	},
	Subcommands: []*cli.Command{&Status},
	Action: func(c *cli.Context) error {
		arn := c.String("arn")
		// checking this here rather than setting it as a required flag so that the "status" subcommand can run without needing an arn
		if arn == "" {
			return fmt.Errorf(`required flag "arn" not set`)
		}
		ctx := c.Context
		bs, err := deploy.BackupStatus(ctx, arn)
		if err != nil {
			return err
		}
		tableName := c.String("table-name")
		if tableName == "" {
			p := &survey.Input{
				Message: "Enter a new table name to restore the backup",
			}
			err = survey.AskOne(p, &tableName, survey.WithValidator(func(ans interface{}) error { return TableNameValidator(ans.(string)) }))
			if err != nil {
				return err
			}
		}
		if err := TableNameValidator(tableName); err != nil {
			return err
		}

		clio.Infof("Restoring Granted Approvals backup: %s to table: %s", aws.ToString(bs.BackupDetails.BackupName), tableName)
		confirm := c.Bool("confirm")
		if !confirm {
			cp := &survey.Confirm{Message: "Do you wish to continue?", Default: true}
			err = survey.AskOne(cp, &confirm)
			if err != nil {
				return err
			}
		}

		_, err = deploy.StartRestore(ctx, arn, tableName)
		if err != nil {
			return err
		}
		clio.Success("Successfully started restoration")
		clio.Successf("To check the status of a restoration run `gdeploy restore status --table-name=%s`", tableName)

		return nil
	},
}

func TableNameValidator(s string) error {
	r := regexp.MustCompile(`[^a-zA-Z0-9_.-]`)
	match := r.MatchString(s)
	if match {
		return fmt.Errorf("value: `%s` must satisfy regular expression pattern: [a-zA-Z0-9_.-]+", s)
	}
	return nil
}
