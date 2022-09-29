package storage

import (
	"context"
	"testing"

	"github.com/common-fate/ddb/ddbtest"
	"github.com/common-fate/granted-approvals/pkg/rule"
	"github.com/common-fate/granted-approvals/pkg/types"
	"github.com/stretchr/testify/assert"
)

func TestListCurrentAccessRules(t *testing.T) {
	db := newTestingStorage(t)

	current := rule.TestAccessRule()
	archived := current
	archived.Current = false
	current.Version = types.NewVersionID()
	ddbtest.PutFixtures(t, db, []*rule.AccessRule{&current, &archived})
	q := &ListCurrentAccessRules{}
	_, err := db.Query(context.TODO(), q)
	assert.NoError(t, err)
	assert.NotContains(t, q.Result, archived)
	assert.Contains(t, q.Result, current)
}
