package request

import (
	"os"
	"time"

	"github.com/benbjohnson/clock"
	"github.com/common-fate/clio"
	"github.com/common-fate/common-fate/pkg/eventhandler"
	"github.com/common-fate/common-fate/pkg/service/accesssvc"
	"github.com/common-fate/common-fate/pkg/service/preflightsvc"
	"github.com/common-fate/common-fate/pkg/service/requestroutersvc"
	"github.com/common-fate/common-fate/pkg/service/rulesvc"
	"github.com/common-fate/common-fate/pkg/service/workflowsvc"
	"github.com/common-fate/common-fate/pkg/service/workflowsvc/runtimes/local"
	"github.com/common-fate/common-fate/pkg/storage"
	"github.com/common-fate/common-fate/pkg/targetgroupgranter"
	"github.com/common-fate/common-fate/pkg/types"
	"github.com/common-fate/ddb"
	"github.com/joho/godotenv"
	"github.com/urfave/cli/v2"
)

var RequestCommand = cli.Command{
	Name:        "request",
	Subcommands: []*cli.Command{&submitCommand},
	Action:      cli.ShowSubcommandHelp,
}

var submitCommand = cli.Command{
	Name: "submit",
	Flags: []cli.Flag{
		&cli.StringFlag{Name: "subject", Required: true},
	},

	Action: func(c *cli.Context) error {
		ctx := c.Context

		clk := clock.New()
		_ = godotenv.Load()
		db, err := ddb.New(ctx, os.Getenv("COMMONFATE_TABLE_NAME"))
		if err != nil {
			return err
		}
		eh := &eventhandler.EventHandler{
			DB: db,
		}
		wf := &workflowsvc.Service{
			Runtime: &local.Runtime{
				DB: db,
				Granter: &targetgroupgranter.Granter{
					DB:          db,
					EventPutter: eh,
					RequestRouter: &requestroutersvc.Service{
						DB: db,
					},
				},
			},
			DB:       db,
			Clk:      clk,
			Eventbus: eh,
		}
		eh.Eventbus = eh
		eh.Workflow = wf

		accsvc := &accesssvc.Service{
			Clock:       clk,
			DB:          db,
			EventPutter: eh,
			Rules: &rulesvc.Service{
				Clock: clk,
				DB:    db,
			},
		}
		presvc := &preflightsvc.Service{
			DB:    db,
			Clock: clk,
		}

		uq := storage.GetUserByEmail{
			Email: c.String("subject"),
		}
		_, err = db.Query(ctx, &uq)
		if err != nil {
			return err
		}

		targets := storage.ListCachedTargets{}
		_, err = db.Query(ctx, &targets, ddb.Limit(1))
		if err != nil {
			return err
		}
		clio.Infow("found targets", "target", targets)
		pre, err := presvc.ProcessPreflight(ctx, *uq.Result, types.CreatePreflightRequest{
			Targets: []string{targets.Result[0].ID()},
		})
		if err != nil {
			return err
		}

		req, err := accsvc.CreateRequest(ctx, *uq.Result, types.CreateAccessRequestRequest{
			PreflightId: pre.ID,
			GroupOptions: []types.CreateAccessRequestGroupOptions{
				{
					Id:     pre.AccessGroups[0].ID,
					Timing: types.RequestAccessGroupTiming{DurationSeconds: 60 * 10}, // 10 minutes
				},
			},
		})
		if err != nil {
			return err
		}
		clio.Infow("created request ", "request", req)
		time.Sleep(11 * time.Minute)
		return nil
	},
}
