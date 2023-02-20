package deployment_schema

import (
	"fmt"
	"os"

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
		// if err := r.AddGoComments("github.com/invopop/jsonschema", "./"); err != nil {
		// 	return err
		// }
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
		fmt.Printf("wrote ./.vscode/schemas/commonfate-deployment-schema.json\n")
		return nil
	}),
}
