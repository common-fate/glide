package storage

import (
	"testing"

	"github.com/common-fate/common-fate/pkg/cache"
	"github.com/common-fate/ddb/ddbtest"
)

func TestListCachedTargetGroupResourceForTargetGroupAndResourceType(t *testing.T) {
	ts := newTestingStorage(t)
	err := ts.deleteAll()
	if err != nil {
		t.Fatal(err)
	}
	re1 := cache.TargetGroupResource{TargetGroupID: "test", Resource: cache.Resource{ID: "value1", Name: "test"}, ResourceType: "testType"}
	re2 := cache.TargetGroupResource{TargetGroupID: "test", Resource: cache.Resource{ID: "value2", Name: "test"}, ResourceType: "testType"}
	re3 := cache.TargetGroupResource{TargetGroupID: "test", Resource: cache.Resource{ID: "value3", Name: "test"}, ResourceType: "testType2"}
	ddbtest.PutFixtures(t, ts.db, []*cache.TargetGroupResource{&re1, &re2, &re3})

	tc := []ddbtest.QueryTestCase{
		{
			Name:  "ok type 1",
			Query: &ListCachedTargetGroupResourceForTargetGroupAndResourceType{TargetGroupID: "test", ResourceType: "testType"},
			Want:  &ListCachedTargetGroupResourceForTargetGroupAndResourceType{TargetGroupID: "test", ResourceType: "testType", Result: []cache.TargetGroupResource{re1, re2}},
		},
		{
			Name:  "ok type 2",
			Query: &ListCachedTargetGroupResourceForTargetGroupAndResourceType{TargetGroupID: "test", ResourceType: "testType2"},
			Want:  &ListCachedTargetGroupResourceForTargetGroupAndResourceType{TargetGroupID: "test", ResourceType: "testType2", Result: []cache.TargetGroupResource{re3}},
		},
	}

	ddbtest.RunQueryTests(t, ts.db, tc)
}
