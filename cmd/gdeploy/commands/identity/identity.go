package identity

import (
	"github.com/common-fate/granted-approvals/cmd/gdeploy/commands/identity/groups"
	"github.com/common-fate/granted-approvals/cmd/gdeploy/commands/identity/sso"
	"github.com/common-fate/granted-approvals/cmd/gdeploy/commands/identity/sync"
	"github.com/common-fate/granted-approvals/cmd/gdeploy/commands/identity/users"
	"github.com/urfave/cli/v2"
)

var Command = cli.Command{
	Name:        "identity",
	Description: "Identity commands are used to manage how your users login to Granted Approvals.\nYou can manage users and groups in the default Cognito user pool or configure your corporate SSO provider.",
	Usage:       "Configure how your users login to Granted Approvals",
	Action:      cli.ShowSubcommandHelp,
	Subcommands: []*cli.Command{&sso.SSOCommand, &users.UsersCommand, &groups.GroupsCommand, &sync.SyncCommand},
}
