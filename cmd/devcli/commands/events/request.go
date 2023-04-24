package events

import (
	"github.com/common-fate/common-fate/pkg/deploy"
	"github.com/common-fate/common-fate/pkg/gevent"
	"github.com/common-fate/common-fate/pkg/storage"
	"github.com/common-fate/ddb"
	"github.com/pkg/errors"
	"github.com/urfave/cli/v2"
)

var requestCommand = cli.Command{
	Name:  "request.created",
	Usage: "emit a request created event",
	Flags: []cli.Flag{
		&cli.StringFlag{Name: "rule", Aliases: []string{"r"}, Usage: "The rule ID", Required: true},
		&cli.StringFlag{Name: "user", Aliases: []string{"u"}, Usage: "The email of the requestor", Required: true},
	},
	Action: func(c *cli.Context) error {
		ctx := c.Context
		do, err := deploy.LoadConfig(deploy.DefaultFilename)
		if err != nil {
			return err
		}
		o, err := do.LoadOutput(ctx)
		if err != nil {
			return err
		}

		// reason := "Deploying Terraform for CF-123"

		db, err := ddb.New(ctx, o.DynamoDBTable)
		if err != nil {
			return err
		}

		q := storage.GetAccessRule{ID: c.String("rule")}

		_, err = db.Query(ctx, &q)
		if err != nil {
			return errors.Wrap(err, "getting access rule")
		}

		u := storage.GetUserByEmail{Email: c.String("user")}

		_, err = db.Query(ctx, &u)
		if err != nil {
			return errors.Wrap(err, "getting requestor")
		}

		e := gevent.RequestCreated{
			// Request: requests.Requestv2{
			// 	// ID:        types.NewRequestID(),
			// 	// CreatedAt: time.Now(),
			// 	// UpdatedAt: time.Now(),

			// 	RequestedBy: *u.Result,
			// },
		}

		s, err := gevent.NewSender(c.Context, gevent.SenderOpts{
			EventBusARN: o.EventBusArn,
		})
		if err != nil {
			return err
		}

		return s.Put(c.Context, e)
	},
}
