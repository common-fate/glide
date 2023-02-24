package storage

import (
	"context"
	"testing"

	"github.com/common-fate/common-fate/pkg/handler"
	"github.com/common-fate/ddb/ddbtest"
	"github.com/stretchr/testify/assert"
)

func TestListHandlers(t *testing.T) {
	db := newTestingStorage(t)
	h := handler.TestHandler("test")
	ddbtest.PutFixtures(t, db, []handler.Handler{h})
	q := &ListHandlers{}
	_, err := db.Query(context.TODO(), q)
	assert.NoError(t, err)
	assert.Contains(t, q.Result, h)
}
