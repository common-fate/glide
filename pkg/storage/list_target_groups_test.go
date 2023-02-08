package storage

import (
	"context"
	"testing"

	"github.com/common-fate/common-fate/pkg/targetgroup"
	"github.com/common-fate/ddb/ddbtest"
	"github.com/stretchr/testify/assert"
)

func TestListTargetGroups(t *testing.T) {
	db := newTestingStorage(t)

	tg1 := targetgroup.TestTargetGroup()
	ddbtest.PutFixtures(t, db, &tg1)
	q := &ListTargetGroups{}
	_, err := db.Query(context.TODO(), q)
	assert.NoError(t, err)
	assert.Contains(t, q.Result, tg1)
}
