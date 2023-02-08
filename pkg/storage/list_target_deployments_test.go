package storage

import (
	"context"
	"testing"

	"github.com/common-fate/common-fate/pkg/targetgroup"
	"github.com/common-fate/ddb/ddbtest"
	"github.com/stretchr/testify/assert"
)

func TestListTargetGroupDeployments_BuildQuery(t *testing.T) {
	db := newTestingStorage(t)

	tg1 := targetgroup.TestTargetGroup()
	ddbtest.PutFixtures(t, db, &tg1)
	q := &ListTargetGroupDeployments{}
	_, err := db.Query(context.TODO(), q)
	assert.NoError(t, err)
	assert.Contains(t, q.Result, tg1)
}
