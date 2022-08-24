package identity

import (
	"github.com/common-fate/granted-approvals/cmd/gdeploy/commands/identity/groups"
	"github.com/common-fate/granted-approvals/cmd/gdeploy/commands/identity/sso"
	"github.com/common-fate/granted-approvals/cmd/gdeploy/commands/identity/sync"
	"github.com/common-fate/granted-approvals/cmd/gdeploy/commands/identity/users"
	"github.com/urfave/cli/v2"
)

// TODO: Update description and usage text.
var Command = cli.Command{
	Name:        "identity",
	Description: "configure identity",
	Usage:       "Add an identity provider",
	Action:      cli.ShowSubcommandHelp,
	Subcommands: []*cli.Command{&sso.SSOCommand, &users.UsersCommand, &groups.GroupsCommand, &sync.SyncCommand},
}
