package terraform

import (
	"github.com/urfave/cli/v2"
)

var Command = cli.Command{
	Name:        "terraform",
	Description: "configure settings for slack integration",
	Subcommands: []*cli.Command{&importTerraformCommand},
}
