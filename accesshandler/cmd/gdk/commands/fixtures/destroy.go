package fixtures

import (
	"os"
	"path/filepath"

	"github.com/common-fate/granted-approvals/pkg/deploy"
	"github.com/common-fate/granted-approvals/pkg/gconfig"
	"github.com/joho/godotenv"
	"github.com/urfave/cli/v2"
	"go.uber.org/zap"
)

var DestroyCommand = cli.Command{
	Name: "destroy",
	Flags: []cli.Flag{
		&cli.PathFlag{Name: "path", Value: "fixtures", Usage: "The path to the fixture JSON file to read or write to"},
		&cli.StringFlag{Name: "name", Aliases: []string{"n"}, Usage: "The name of the provider to generate fixtures for", Required: true},
	},
	Action: func(c *cli.Context) error {
		_ = godotenv.Load()

		ctx := c.Context

		name := c.String("name")
		g, err := LookupGenerator(name)
		if err != nil {
			return err
		}
		ac := deploy.EnvDeploymentConfig{}
		pc, err := ac.ReadProviders(ctx)
		if err != nil {
			return err
		}

		// configure the generator if it supports it
		if configer, ok := g.(gconfig.Configer); ok {
			p := pc[name]
			err = configer.Config().Load(ctx, &gconfig.MapLoader{Values: p.With})
			if err != nil {
				return err
			}
		}

		// init the generator if it supports it
		if configurer, ok := g.(gconfig.Initer); ok {
			err = configurer.Init(ctx)
			if err != nil {
				return err
			}
		}

		p := c.Path("path")
		fixturePath := filepath.Join(p, name+".json")

		data, err := os.ReadFile(fixturePath)
		if err != nil {
			return err
		}

		err = g.Destroy(ctx, data)
		if err != nil {
			return err
		}

		err = os.Remove(fixturePath)
		if err != nil {
			return err
		}

		zap.S().Infow("destroyed fixture", "file", fixturePath)

		return nil
	},
}
