package ddbhelpers

import (
	"context"

	"github.com/common-fate/ddb"
)

func QueryPages(
	ctx context.Context, c ddb.Storage, qb ddb.QueryBuilder,
	f func(pageResult *ddb.QueryResult, pageQueryBuilder ddb.QueryBuilder, lastPage bool) bool,
	opts ...func(*ddb.QueryOpts),
) error {
	var nextPage string
	for {
		result, err := c.Query(ctx, qb, ddb.Page(nextPage))
		if err != nil {
			return err
		}
		lastPage := result.NextPage == ""
		if !f(result, qb, lastPage) {
			break
		}
		if lastPage {
			break
		}
		nextPage = result.NextPage
	}
	return nil
}
