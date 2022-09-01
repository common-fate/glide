package users

import (
	"github.com/urfave/cli/v2"
)

var UsersCommand = cli.Command{
	Name:        "users",
	Description: "Add or remove users from the default cognito user pool.\nThese commands are only available when you are using the default Cognito user pool. If you have connected an SSO provider, like Okta, Google or AzureAD, use those tools to manage your users and groups instead.",
	Subcommands: []*cli.Command{&CreateCommand, &DeleteCommand},
	Action:      cli.ShowSubcommandHelp,
}
