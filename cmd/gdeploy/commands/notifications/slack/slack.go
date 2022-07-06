package slack

import (
	"github.com/urfave/cli/v2"
)

var Command = cli.Command{
	Name:        "slack",
	Description: "configure settings for slack integration",
	Subcommands: []*cli.Command{&configureSlackCommand, &testSlackCommand, &disableSlackCommand},
}
