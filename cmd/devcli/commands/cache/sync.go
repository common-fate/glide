package cache

import (
	"errors"
	"os"
	"strings"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/common-fate/common-fate/pkg/cachesync"
	"github.com/common-fate/common-fate/pkg/handler"
	"github.com/common-fate/common-fate/pkg/rule"
	"github.com/common-fate/common-fate/pkg/service/cachesvc"
	"github.com/common-fate/common-fate/pkg/service/requestroutersvc"
	"github.com/common-fate/common-fate/pkg/storage"
	"github.com/common-fate/ddb"
	"github.com/common-fate/provider-registry-sdk-go/pkg/providerregistrysdk"
	"github.com/joho/godotenv"
	"github.com/urfave/cli/v2"
)

var CacheCommand = cli.Command{
	Name:        "cache",
	Subcommands: []*cli.Command{&syncCommand, &targetsCommand},
	Action:      cli.ShowSubcommandHelp,
}

var syncCommand = cli.Command{
	Name:        "sync",
	Flags:       []cli.Flag{&cli.StringSliceFlag{Name: "deployment-mappings"}},
	Description: "Sync cache",
	Action: func(c *cli.Context) error {
		ctx := c.Context
		_ = godotenv.Load()
		db, err := ddb.New(ctx, os.Getenv("COMMONFATE_TABLE_NAME"))
		if err != nil {
			return err
		}

		for _, dm := range c.StringSlice("deployment-mappings") {
			kv := strings.Split(dm, ":")
			if len(kv) != 2 {
				return errors.New("deployment-mapping is invalid")
			}
			handler.LocalDeploymentMap[kv[0]] = kv[1]
		}

		// this configuration means the pdk will use the local test runtime instead of calling out to lambda
		syncer := cachesync.CacheSyncer{
			DB: db,
			Cache: cachesvc.Service{
				DB: db,
				RequestRouter: &requestroutersvc.Service{
					DB: db,
				},
			},
		}

		err = syncer.Sync(ctx)
		if err != nil {
			return err
		}

		return nil
	},
}

var targetsCommand = cli.Command{
	Name:        "targets",
	Flags:       []cli.Flag{&cli.StringSliceFlag{Name: "deployment-mappings"}},
	Description: "generate targets",
	Action: func(c *cli.Context) error {
		ctx := c.Context
		_ = godotenv.Load()

		db, err := ddb.New(ctx, os.Getenv("COMMONFATE_TABLE_NAME"))
		if err != nil {
			return err
		}

		for _, dm := range c.StringSlice("deployment-mappings") {
			kv := strings.Split(dm, ":")
			if len(kv) != 2 {
				return errors.New("deployment-mapping is invalid")
			}
			handler.LocalDeploymentMap[kv[0]] = kv[1]
		}

		s := cachesvc.Service{
			DB: db,
		}

		q := storage.GetTargetGroup{
			ID: "cloudwatch",
		}
		// q := storage.ListTargetGroups{}
		_, err = db.Query(ctx, &q)
		if err != nil {
			return err
		}
		q.Result.Schema.Properties = map[string]providerregistrysdk.TargetField{"log_group": {
			Resource: aws.String("LogGroup"),
		}}

		ar := rule.AccessRule{
			ID: "demo-rule",
			Targets: []rule.Target{
				{
					TargetGroup: *q.Result,
					FieldFilterExpessions: map[string]rule.FieldFilterExpessions{
						"accountId":        {},
						"permissionSetArn": {},
					},
				},
			},
		}
		err = db.Put(ctx, &ar)
		if err != nil {
			return err
		}

		err = s.RefreshCachedTargets(ctx)
		if err != nil {
			return err
		}

		return nil
	},
}
