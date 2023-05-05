package migrate

import (
	"context"

	"github.com/common-fate/clio"
	"github.com/common-fate/ddb"
	"github.com/urfave/cli/v2"
)

//
// Deprecated, repurpose this command if we have to run another migration later
//

var MigrateCommand = cli.Command{
	Name:        "migrate",
	Description: "Migrate from v0.15 to v1.0.0",
	Usage:       "Migrate from v0.15 to v1.0.0",
	Action: func(c *cli.Context) error {
		// ctx := c.Context

		clio.Success("Successfully migrated from config version 1 -> 2")
		return nil
	},
}

type Migrator struct {
	DB ddb.Storage
}

// upgrade target groups
// update the schema from the registry to the new type
func (m *Migrator) TargetGroups(ctx context.Context) {

}

// migrate requests
// list request
// create v2 requests
func (m *Migrator) Requests(ctx context.Context) {

}

// migrate access rules
// convert targets to multiple
func (m *Migrator) AccessRules(ctx context.Context) {

}
