package dbupdate

import (
	"context"

	"github.com/common-fate/ddb"
	"github.com/common-fate/granted-approvals/pkg/access"
	"github.com/common-fate/granted-approvals/pkg/storage"
)

type UpdateRequestOpts struct {
	Reviewers []access.Reviewer
}

// WithReviewers allows reviewers to be passed in if they have already be fetched in a previous query
func WithReviewers(r []access.Reviewer) func(*UpdateRequestOpts) {
	return func(uro *UpdateRequestOpts) {
		uro.Reviewers = r
	}
}

// GetUpdateRequestItems returns a slice of ddb.keyers which needs to be written to update this request
// all the items returned have been updated with the input request
func GetUpdateRequestItems(ctx context.Context, db ddb.Storage, r access.Request, opts ...func(*UpdateRequestOpts)) ([]ddb.Keyer, error) {
	var o UpdateRequestOpts
	for _, opt := range opts {
		opt(&o)
	}

	if o.Reviewers == nil {
		rq := storage.ListRequestReviewers{RequestID: r.ID}
		_, err := db.Query(ctx, &rq)
		if err != nil {
			return nil, err
		}
		o.Reviewers = rq.Result
	}

	items := make([]ddb.Keyer, len(o.Reviewers)+1)
	items[0] = &r
	for i, rv := range o.Reviewers {
		rvc := rv
		rvc.Request = r
		items[1+i] = &rvc
	}
	return items, nil
}
