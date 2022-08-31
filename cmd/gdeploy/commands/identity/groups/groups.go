package groups

import (
	"github.com/common-fate/granted-approvals/cmd/gdeploy/commands/identity/groups/members"
	"github.com/urfave/cli/v2"
)

var GroupsCommand = cli.Command{
	Name:        "group",
	Subcommands: []*cli.Command{&CreateCommand, &DeleteCommand, &members.MembersCommand},
	Action:      cli.ShowSubcommandHelp,
}
