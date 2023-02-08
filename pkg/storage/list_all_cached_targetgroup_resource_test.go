package storage

import (
	"testing"

	"github.com/common-fate/common-fate/pkg/cache"
	"github.com/common-fate/ddb/ddbtest"
)

func TestListCachedTargetGroupResource(t *testing.T) {
	s := newTestingStorage(t)
	re1 := cache.TargateGroupResource{TargetGroupID: "test", Resource: cache.Resource{ID: "value1", Name: "test"}, ResourceType: "testType"}
	re2 := cache.TargateGroupResource{TargetGroupID: "test", Resource: cache.Resource{ID: "value2", Name: "test"}, ResourceType: "testType"}
	re3 := cache.TargateGroupResource{TargetGroupID: "test", Resource: cache.Resource{ID: "value3", Name: "test"}, ResourceType: "testType2"}
	ddbtest.PutFixtures(t, s, []*cache.TargateGroupResource{&re1, &re2, &re3})

	tc := []ddbtest.QueryTestCase{
		{
			Name:  "ok type 1",
			Query: &ListCachedTargetGroupResource{TargetGroupID: "test", ResourceType: "testType"},
			Want:  &ListCachedTargetGroupResource{TargetGroupID: "test", ResourceType: "testType", Result: []cache.TargateGroupResource{re1, re2}},
		},
		{
			Name:  "ok type 2",
			Query: &ListCachedTargetGroupResource{TargetGroupID: "test", ResourceType: "testType2"},
			Want:  &ListCachedTargetGroupResource{TargetGroupID: "test", ResourceType: "testType2", Result: []cache.TargateGroupResource{re3}},
		},
	}

	ddbtest.RunQueryTests(t, s, tc)
}
