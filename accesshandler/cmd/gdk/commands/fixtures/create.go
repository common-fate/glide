package fixtures

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"

	"github.com/common-fate/granted-approvals/accesshandler/pkg/config"
	"github.com/common-fate/granted-approvals/accesshandler/pkg/genv"
	"github.com/common-fate/granted-approvals/accesshandler/pkg/providers"
	"github.com/joho/godotenv"
	"github.com/urfave/cli/v2"
	"go.uber.org/zap"
)

var CreateCommand = cli.Command{
	Name: "create",
	Flags: []cli.Flag{
		&cli.PathFlag{Name: "path", Value: "fixtures", Usage: "the path to the fixture JSON file to read or write to"},
		&cli.StringFlag{Name: "name", Aliases: []string{"n"}, Usage: "the name of the provider to generate fixtures for", Required: true},
	},
	Action: func(c *cli.Context) error {
		ctx := c.Context

		err := godotenv.Load()
		if err != nil {
			return err
		}

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

		pc, err := config.ReadProviderConfig(ctx, "local")
		if err != nil {
			return err
		}
		var configMap map[string]json.RawMessage
		err = json.Unmarshal(pc, &configMap)
		if err != nil {
			return err
		}

		// configure the generator if it supports it
		if configer, ok := g.(providers.Configer); ok {
			err = configer.Config().Load(ctx, genv.JSONLoader{Data: configMap[name]})
			if err != nil {
				return err
			}
		}

		// init the generator if it supports it
		if configurer, ok := g.(providers.Initer); ok {
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

		err = ioutil.WriteFile(fixturePath, fixtures, 0666)
		if err != nil {
			return err
		}

		zap.S().Infow("created fixture", "file", fixturePath)

		return nil
	},
}
