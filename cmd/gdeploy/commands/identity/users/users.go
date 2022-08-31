package users

import (
	"github.com/urfave/cli/v2"
)

var UsersCommand = cli.Command{
	Name:        "users",
	Description: "Add or remove users from default cognito identity provider",
	Subcommands: []*cli.Command{&CreateCommand, &DeleteCommand},
	Action:      cli.ShowSubcommandHelp,
}
