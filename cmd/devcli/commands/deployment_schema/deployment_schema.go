package deployment_schema

import (
	"fmt"
	"os"

	"github.com/common-fate/clio"
	"github.com/common-fate/common-fate/pkg/deploy"
	"github.com/invopop/jsonschema"
	"github.com/urfave/cli/v2"
)

// DeploymentSchemaCommand
// to run in dev
// go run cmd/devcli/main.go deployment-schema
var DeploymentSchemaCommand = cli.Command{
	Name: "deployment-schema",
	Action: cli.ActionFunc(func(c *cli.Context) error {

		r := new(jsonschema.Reflector)
		s := r.Reflect(&deploy.Config{})

		marshalled, err := s.MarshalJSON()
		if err != nil {
			return err
		}

		// convert to string
		fmt.Printf("%s\n\n", marshalled)

		// save this to local fs using os in file name ./vscode/schemas/commonfate-deployment-schema.json
		// if it doesnt exist make it
		os.MkdirAll("./.vscode/schemas", 0755)
		err = os.WriteFile("./.vscode/schemas/commonfate-deployment-schema.json", marshalled, 0644)
		if err != nil {
			return err
		}
		clio.Success("wrote ./.vscode/schemas/commonfate-deployment-schema.json\n")

		clio.Successf(`To use this schema in vscode, add the following to your settings.json:\n\n
"yaml.schemas": {
	".vscode/schemas/commonfate-deployment-schema.json": "deployment.yml"
},
		`)

		// TODO: we coudl also add auto installation for this in vscode settings
		// err = os.MkdirAll("./.vscode", 0755)
		// if err != nil {
		// 	return err
		// }
		// // read first
		// f, err := os.ReadFile("./.vscode/settings.json")
		// if err != nil {
		// 	return err
		// }

		return nil
	}),
}
