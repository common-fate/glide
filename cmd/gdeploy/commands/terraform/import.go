package terraform

import (
	"bytes"
	"embed"
	"os"
	"text/template"

	"github.com/common-fate/clio"
	"github.com/common-fate/common-fate/pkg/deploy"
	"github.com/common-fate/common-fate/pkg/rule"
	"github.com/common-fate/common-fate/pkg/storage"
	"github.com/common-fate/ddb"
	"github.com/olekukonko/tablewriter"
	"github.com/urfave/cli/v2"
)

//go:embed templates
var templateFiles embed.FS

var importTerraformCommand = cli.Command{
	Name:        "import",
	Description: "import click ops access rules into hcl terraform format",
	Action: func(c *cli.Context) error {
		ctx := c.Context

		dc, err := deploy.ConfigFromContext(ctx)
		if err != nil {
			return err
		}
		o, err := dc.LoadOutput(ctx)
		if err != nil {
			return err
		}

		db, err := ddb.New(ctx, o.DynamoDBTable)
		if err != nil {
			return err
		}

		l := storage.ListAccessRulesForStatus{Status: rule.ACTIVE}
		_, err = db.Query(ctx, &l)
		if err != nil {
			return err
		}
		rules := l.Result

		out, err := os.Create("test.tf")
		if err != nil {
			return err
		}

		var output []byte
		for _, r := range rules {
			if r.Metadata.CreatedBy != "bot_governance_api" {

				tf_rule, err := WriteAccessRuleToHCL(r)

				output = append(output, tf_rule...)
				if err != nil {
					return err
				}
			}

		}

		_, err = out.Write(output)
		if err != nil {
			return err
		}
		clio.Success("Copied Access Rules to .tf file")
		clio.Info("Import these rules to your Terraform state, using the following ID's")
		writeAccessRuleTable(rules)

		return nil
	},
}

func WriteAccessRuleToHCL(ar rule.AccessRule) ([]byte, error) {

	t := template.New("t")
	t, err := t.Parse(templateString)
	if err != nil {
		return nil, err
	}
	var tpl bytes.Buffer

	err = t.Execute(&tpl, ar)
	if err != nil {
		return nil, err
	}

	return tpl.Bytes(), nil
}

func writeAccessRuleTable(ar []rule.AccessRule) {
	table := tablewriter.NewWriter(os.Stdout)
	table.SetHeader([]string{"Name", "ID"})
	table.SetAutoWrapText(false)
	table.SetAutoFormatHeaders(true)
	table.SetHeaderAlignment(tablewriter.ALIGN_LEFT)
	table.SetAlignment(tablewriter.ALIGN_LEFT)
	table.SetCenterSeparator("")
	table.SetColumnSeparator("")
	table.SetRowSeparator("")
	table.SetHeaderLine(false)
	table.SetBorder(false)

	for _, a := range ar {

		table.Append([]string{
			a.Name, a.ID,
		})
	}
	table.Render()
}

var templateString string = `resource "commonfate_access_rule" "{{ .Name }}" {
  name ="{{ .Name }}"
  description="{{ .Description }}"
  groups=[{{ range $index, $element := .Groups}}
        "{{$element}}",{{end}}
  ]
  target=[
   {{ range $key, $value := .Target.With}}
	{
        field="{{$key}}",
		value="{{$value}}"
	},
	{{end}}

	{{ range $key, $value := .Target.WithSelectable}}
	{
        field="{{$key}}",
		value=[
			"{{$value}}"
		]
	},
	{{end}}
  ]
  target_provider_id="{{ .Target.ProviderID }}"
  duration={{ .TimeConstraints.MaxDurationSeconds }}
}
`

//  target=[
//     {
//       field="accountId"
//       value=["632700053629"]
//     },
//     {
//       field="permissionSetArn"
//       value=["arn:aws:sso:::permissionSet/ssoins-825968feece9a0b6/ps-dda57372ebbfeb94"]
//     }
//   ]
