package access

import (
	"time"

	"github.com/common-fate/ddb"
	"github.com/common-fate/granted-approvals/pkg/storage/keys"
	"github.com/common-fate/granted-approvals/pkg/types"
)

type Bookmark struct {
	// ID
	ID     string `json:"id" dynamodbav:"id"`
	UserID string `json:"userId" dynamodbav:"userId"`
	Name   string `json:"name" dynamodbav:"name"`
	// Rule is the ID of the Access Rule which the request relates to.
	Rule            string                `json:"rule" dynamodbav:"rule"`
	Data            RequestData           `json:"data" dynamodbav:"data"`
	RequestedTiming Timing                `json:"requestedTiming" dynamodbav:"requestedTiming"`
	With            []map[string][]string `json:"with" dynamodbav:"with"`
	CreatedAt       time.Time             `json:"createdAt" dynamodbav:"createdAt"`
	UpdatedAt       time.Time             `json:"updatedAt" dynamodbav:"updatedAt"`
}

func (b Bookmark) ToAPI() types.Bookmark {
	return types.Bookmark{
		Id:     b.ID,
		Name:   b.Name,
		RuleId: b.Rule,
	}
}

func (b Bookmark) ToAPIDetail() types.BookmarkDetail {
	bm := types.BookmarkDetail{
		Id:     b.ID,
		Name:   b.Name,
		Reason: b.Data.Reason,
		Timing: b.RequestedTiming.ToAPI(),
	}

	for _, w := range b.With {
		bm.With = append(bm.With, types.CreateRequestWith{
			AdditionalProperties: w,
		})
	}
	return bm
}

func (b *Bookmark) DDBKeys() (ddb.Keys, error) {
	keys := ddb.Keys{
		PK: keys.Bookmark.PK1,
		SK: keys.Bookmark.SK1(b.UserID, b.ID),
	}
	return keys, nil
}
