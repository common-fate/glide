package groups

import (
	"github.com/common-fate/granted-approvals/cmd/gdeploy/commands/identity/groups/members"
	"github.com/urfave/cli/v2"
)

var GroupsCommand = cli.Command{
	Name:        "groups",
	Description: "Add or remove groups from the default cognito user pool.\nThese commands are only available when you are using the default Cognito user pool. If you have connected an SSO provider, like Okta, Google or AzureAD, use those tools to manage your users and groups instead.",
	Subcommands: []*cli.Command{&CreateCommand, &DeleteCommand, &members.MembersCommand},
	Action:      cli.ShowSubcommandHelp,
}
