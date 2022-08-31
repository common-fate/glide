package fixtures

import (
	"encoding/json"
	"os"
	"path/filepath"

	"github.com/common-fate/granted-approvals/accesshandler/pkg/config"
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
		pc, err := config.ReadProviderConfig(ctx)
		if err != nil {
			return err
		}
		var configMap map[string]json.RawMessage
		err = json.Unmarshal(pc, &configMap)
		if err != nil {
			return err
		}

		// configure the generator if it supports it
		if configer, ok := g.(gconfig.Configer); ok {
			err = configer.Config().Load(ctx, gconfig.JSONLoader{Data: configMap[name]})
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
