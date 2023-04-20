package requests

// import (
// 	"github.com/aws/aws-sdk-go-v2/aws"
// 	"github.com/common-fate/common-fate/pkg/config"
// 	"github.com/common-fate/common-fate/pkg/requests"
// 	"github.com/common-fate/common-fate/pkg/target"
// 	"github.com/common-fate/ddb"
// 	"github.com/common-fate/provider-registry-sdk-go/pkg/providerregistrysdk"
// 	"github.com/joho/godotenv"
// 	"github.com/sethvargo/go-envconfig"
// 	"github.com/urfave/cli/v2"
// )

// var SeedCommand = cli.Command{
// 	Name:        "seed",
// 	Description: "Seeds some dummy data into the dynamo table for testing new workflows",
// 	Action: cli.ActionFunc(func(c *cli.Context) error {
// 		ctx := c.Context
// 		// Read from the .env file
// 		var cfg config.HealthCheckerConfig
// 		_ = godotenv.Load()
// 		err := envconfig.Process(ctx, &cfg)
// 		if err != nil {
// 			return err
// 		}
// 		db, err := ddb.New(ctx, cfg.TableName)
// 		if err != nil {
// 			return err
// 		}

// 		items := []ddb.Keyer{}

// 		//create an entitlement

// 		ent := target.Group{
// 			ID: "aws-tg",
// 			From: target.From{
// 				Kind:      "Account",
// 				Name:      "AWS",
// 				Publisher: "common-fate",
// 				Version:   "v0.1.0",
// 			},
// 			Schema: providerregistrysdk.Target{
// 				Type: "",
// 				Properties: map[string]providerregistrysdk.TargetField{
// 					"accountId":        {Title: aws.String("Account")},
// 					"permissionSetArn": {Title: aws.String("Permission Set")},
// 				},
// 			},
// 		}

// 		// ent := requests.Entitlement{
// 		// 	ID: types.NewEntitlementID(),
// 		// 	Kind: requests.TargetFrom{
// 		// 		Kind:      "Account",
// 		// 		Name:      "AWS",
// 		// 		Publisher: "common-fate",
// 		// 		Version:   "v0.1.0",
// 		// 	},
// 		// 	OptionSchema: types.TargetSchema{
// 		// 		AdditionalProperties: map[string]types.TargetArgument{
// 		// 			"accountId":        {Title: "Account"},
// 		// 			"permissionSetArn": {Title: "Permission Set"},
// 		// 		},
// 		// 	},
// 		// }

// 		ent2 := target.Group{
// 			ID: "okta-tg",
// 			From: target.From{
// 				Kind:      "Group",
// 				Name:      "Okta",
// 				Publisher: "common-fate",
// 				Version:   "v0.1.0",
// 			},
// 			Schema: providerregistrysdk.Target{
// 				Type: "",
// 				Properties: map[string]providerregistrysdk.TargetField{
// 					"groupName": {Title: aws.String("Group Name")},
// 				},
// 			},
// 		}

// 		items = append(items, &ent)
// 		items = append(items, &ent2)

// 		//create some options
// 		opt1 := requests.ResourceOption{

// 			Label: "accountId",
// 			Value: "123456789012",
// 			Provider: requests.TargetFrom{
// 				Kind:      "Account",
// 				Name:      "AWS",
// 				Publisher: "common-fate",
// 				Version:   "v0.1.0",
// 			},
// 			Type:        "Account",
// 			TargetGroup: "test",
// 			AccessRules: []string{
// 				"test",
// 			},
// 		}
// 		opt1a := requests.ResourceOption{

// 			Label: "accountId",
// 			Value: "13579012345",
// 			Provider: requests.TargetFrom{
// 				Kind:      "Account",
// 				Name:      "AWS",
// 				Publisher: "common-fate",
// 				Version:   "v0.1.0",
// 			},
// 		}
// 		opt1b := requests.ResourceOption{

// 			Label: "accountId",
// 			Value: "583847583929",
// 			Provider: requests.TargetFrom{
// 				Kind:      "Account",
// 				Name:      "AWS",
// 				Publisher: "common-fate",
// 				Version:   "v0.1.0",
// 			},
// 			Type:        "Account",
// 			TargetGroup: "test",
// 			AccessRules: []string{
// 				"test",
// 			},
// 		}
// 		opt2 := requests.ResourceOption{
// 			Label: "permissionSetArn",
// 			Value: "123-abc",
// 			Provider: requests.TargetFrom{
// 				Kind:      "Account",
// 				Name:      "AWS",
// 				Publisher: "common-fate",
// 				Version:   "v0.1.0",
// 			},
// 			Type:        "Account",
// 			TargetGroup: "test",
// 			AccessRules: []string{
// 				"test",
// 			},
// 			RelatedTo: []string{"123456789012"},
// 		}
// 		opt2a := requests.ResourceOption{
// 			Label: "permissionSetArn",
// 			Value: "bar",
// 			Provider: requests.TargetFrom{
// 				Kind:      "Account",
// 				Name:      "AWS",
// 				Publisher: "common-fate",
// 				Version:   "v0.1.0",
// 			},
// 			Type:        "Account",
// 			TargetGroup: "test",
// 			AccessRules: []string{
// 				"test",
// 			},
// 			RelatedTo: []string{"123456789012"},
// 		}
// 		opt2b := requests.ResourceOption{
// 			Label: "permissionSetArn",
// 			Value: "foo",
// 			Provider: requests.TargetFrom{
// 				Kind:      "Account",
// 				Name:      "AWS",
// 				Publisher: "common-fate",
// 				Version:   "v0.1.0",
// 			},
// 			Type:        "Account",
// 			TargetGroup: "test",
// 			AccessRules: []string{
// 				"test",
// 			},
// 			RelatedTo: []string{"123456789012"},
// 		}

// 		opt3 := requests.ResourceOption{
// 			Label: "groupName",
// 			Value: "This is a okta group",
// 			Provider: requests.TargetFrom{
// 				Kind:      "Group",
// 				Name:      "Okta",
// 				Publisher: "common-fate",
// 				Version:   "v0.1.0",
// 			},
// 			Type:        "Account",
// 			TargetGroup: "test",
// 			AccessRules: []string{
// 				"diff-group",
// 			},
// 		}
// 		items = append(items, &opt1)
// 		items = append(items, &opt1a)
// 		items = append(items, &opt1b)
// 		items = append(items, &opt2)
// 		items = append(items, &opt2a)
// 		items = append(items, &opt2b)
// 		items = append(items, &opt3)

// 		err = db.PutBatch(ctx, items...)
// 		if err != nil {
// 			return err
// 		}
// 		return nil

// 	}),
// }
