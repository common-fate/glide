package rulesvc

import (
	"context"

	"github.com/common-fate/analytics-go"
	"github.com/common-fate/common-fate/pkg/access"
	"github.com/common-fate/common-fate/pkg/identity"
	"github.com/common-fate/common-fate/pkg/rule"
	"github.com/common-fate/common-fate/pkg/storage"
	"github.com/common-fate/common-fate/pkg/storage/dbupdate"
	"github.com/common-fate/common-fate/pkg/types"
	"github.com/common-fate/ddb"
)

func (s *Service) ArchiveAccessRule(ctx context.Context, user *identity.User, in rule.AccessRule) (*rule.AccessRule, error) {
	if in.Status == rule.ARCHIVED {
		return nil, ErrAccessRuleAlreadyArchived
	}

	// make a copy of the existing version which will be mutated
	newVersion := in
	newVersion.Status = rule.ARCHIVED
	newVersion.Metadata.UpdatedAt = s.Clock.Now()
	newVersion.Version = types.NewVersionID()
	newVersion.Current = true

	// Set the existing version to not current
	in.Current = false

	// creates a new version entry as well as setting the current version
	items := []ddb.Keyer{&newVersion, &in}

	// pagination in case of many many pending requests
	hasMore := true
	var next string
	for hasMore {
		// list requests for rule and status
		// We currently don't have an access pattern for this directly, so we can fetch all the pending requests and filter them in go
		q := storage.ListRequestsForStatus{Status: access.PENDING}
		var opts []func(*ddb.QueryOpts)
		if next != "" {
			opts = append(opts, ddb.Page(next))
		}

		res, err := s.DB.Query(ctx, &q, opts...)
		if err != nil && err != ddb.ErrNoItems {
			return nil, err
		}
		next = res.NextPage
		hasMore = next != ""

		for _, r := range q.Result {
			if r.Rule == in.ID {
				r.Status = access.CANCELLED
				r.UpdatedAt = s.Clock.Now()
				updateItems, err := dbupdate.GetUpdateRequestItems(ctx, s.DB, r)
				if err != nil {
					return nil, err
				}
				items = append(items, updateItems...)
			}
		}
	}

	err := s.DB.PutBatch(ctx, items...)
	if err != nil {
		return nil, err
	}

	// analytics event
	analytics.FromContext(ctx).Track(&analytics.RuleArchived{
		ArchivedBy: user.ID,
		RuleID:     in.ID,
	})

	return &newVersion, nil
}
