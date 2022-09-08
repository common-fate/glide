package fixtures

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/common-fate/granted-approvals/pkg/deploy"
	"github.com/common-fate/granted-approvals/pkg/gconfig"
	"github.com/joho/godotenv"
	"github.com/pkg/errors"
	"github.com/urfave/cli/v2"
	"go.uber.org/zap"
)

var CreateCommand = cli.Command{
	Name: "create",
	Flags: []cli.Flag{
		&cli.PathFlag{Name: "path", Value: "fixtures", Usage: "The path to the fixture JSON file to read or write to"},
		&cli.StringFlag{Name: "name", Aliases: []string{"n"}, Usage: "The name of the provider to generate fixtures for", Required: true},
	},
	Action: func(c *cli.Context) error {
		ctx := c.Context

		_ = godotenv.Load()

		name := c.String("name")
		g, err := LookupGenerator(name)
		if err != nil {
			return err
		}

		// ensure that the fixture file doesn't already exist - return an error if it does to prevent
		// multiple fixtures being created.
		parentFolder := c.Path("path")

		fixturePath := filepath.Join(parentFolder, name+".json")
		if _, err := os.Stat(fixturePath); err == nil {
			return fmt.Errorf("fixture already exists (%s). Use 'gdk fixtures delete --name %s' to remove it before generating it again", fixturePath, name)
		}

		ac := deploy.EnvDeploymentConfig{}
		pc, err := ac.ReadProviders(ctx)
		if err != nil {
			return errors.Wrap(err, "reading providers")
		}

		// configure the generator if it supports it
		if configer, ok := g.(gconfig.Configer); ok {
			p := pc[name]
			err = configer.Config().Load(ctx, &gconfig.MapLoader{Values: p.With})
			if err != nil {
				return errors.Wrap(err, "loading config")
			}
		}

		// init the generator if it supports it
		if configurer, ok := g.(gconfig.Initer); ok {
			err = configurer.Init(ctx)
			if err != nil {
				return err
			}
		}

		fixtures, err := g.Generate(ctx)
		if err != nil {
			return err
		}

		err = os.MkdirAll(parentFolder, os.ModePerm)
		if err != nil {
			return err
		}

		err = os.WriteFile(fixturePath, fixtures, 0666)
		if err != nil {
			return err
		}

		zap.S().Infow("created fixture", "file", fixturePath)

		return nil
	},
}
