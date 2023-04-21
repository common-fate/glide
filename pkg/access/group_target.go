package access

import (
	"time"

	"github.com/common-fate/common-fate/pkg/storage/keys"
	"github.com/common-fate/common-fate/pkg/target"
	"github.com/common-fate/common-fate/pkg/types"
	"github.com/common-fate/ddb"
)

type GroupTarget struct {
	ID        string `json:"id" dynamodbav:"id"`
	GroupID   string `json:"groupId" dynamodbav:"groupId"`
	RequestID string `json:"requestId" dynamodbav:"requestId"`
	// Also denormalised across all the request items
	RequestStatus types.RequestStatus `json:"requestStatus" dynamodbav:"requestStatus"`
	RequestedBy   string              `json:"requestedBy" dynamodbav:"requestedBy"`
	// The id of the cache.Target which was used to select this on the request.
	// the cache item is subject to be deleted so this cacheID may not always exist in the future after the grant is created
	TargetCacheID   string      `json:"cacheId" dynamodbav:"cacheId"`
	TargetGroupID   string      `json:"targetGroupId" dynamodbav:"targetGroupId"`
	TargetGroupFrom target.From `json:"targetGroupFrom" dynamodbav:"targetGroupFrom"`
	Fields          []Field     `json:"fields" dynamodbav:"fields"`
	// The grant will be populated when this target is submitted to be provisioned
	// The start and end time are calculated and stored on the grant when it is provisioned
	Grant     *Grant    `json:"grant" dynamodbav:"grant"`
	CreatedAt time.Time `json:"createdAt" dynamodbav:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt" dynamodbav:"updatedAt"`
	// request reviewers are users who have one or more groups to review on the request as a whole
	RequestReviewers []string `json:"requestReviewers" dynamodbav:"requestReviewers, set"`
}

type Grant struct {
	// The user email
	Subject string                               `json:"subject" dynamodbav:"subject"`
	Status  types.RequestAccessGroupTargetStatus `json:"status" dynamodbav:"status"`
	//the time which the grant starts
	Start time.Time `json:"start" dynamodbav:"start"`
	//the time the grant is scheduled to end
	End          time.Time `json:"end" dynamodbav:"end"`
	Instructions *string   `json:"instructions" dynamodbav:"instructions"`
}
type Field struct {
	ID               string     `json:"id" dynamodbav:"id"`
	FieldTitle       string     `json:"fieldTitle" dynamodbav:"fieldTitle"`
	FieldDescription *string    `json:"fieldDescription" dynamodbav:"fieldDescription"`
	ValueLabel       string     `json:"valueLabel" dynamodbav:"valueLabel"`
	ValueDescription *string    `json:"valueDescription" dynamodbav:"valueDescription"`
	Value            FieldValue `json:"value" dynamodbav:"value"`
}
type FieldValue struct {
	Type  string `json:"type" dynamodbav:"type"`
	Value string `json:"value" dynamodbav:"value"`
}

func (f *Field) ToAPI() types.TargetField {
	return types.TargetField{
		Id:               f.ID,
		FieldDescription: f.FieldDescription,
		FieldTitle:       f.FieldTitle,
		Value:            f.Value.Value,
		ValueDescription: f.ValueDescription,
		ValueLabel:       f.ValueLabel,
	}
}
func (g *GroupTarget) ToAPI() types.RequestAccessGroupTarget {
	grant := types.RequestAccessGroupTarget{
		AccessGroupId:   g.GroupID,
		Id:              g.ID,
		RequestId:       g.RequestID,
		Status:          types.RequestAccessGroupTargetStatusPENDINGPROVISIONING,
		Fields:          []types.TargetField{},
		TargetGroupFrom: g.TargetGroupFrom.ToAPI(),
		TargetGroupId:   g.TargetGroupID,
	}
	if g.Grant != nil {
		grant.Status = g.Grant.Status
	}
	for _, field := range g.Fields {
		grant.Fields = append(grant.Fields, field.ToAPI())
	}

	return grant
}
func (i *GroupTarget) DDBKeys() (ddb.Keys, error) {
	keys := ddb.Keys{
		PK:     keys.AccessRequestGroupTarget.PK1,
		SK:     keys.AccessRequestGroupTarget.SK1(i.RequestID, i.GroupID, i.ID),
		GSI1PK: keys.AccessRequestGroupTarget.GSI1PK(i.RequestedBy),
		GSI1SK: keys.AccessRequestGroupTarget.GSI1SK(i.RequestID, RequestStatusToPastOrUpcoming(i.RequestStatus), i.GroupID, i.ID),
	}
	return keys, nil
}

type Instructions struct {
	GroupTargetID string `json:"id" dynamodbav:"id"`
	RequestedBy   string `json:"requestedBy" dynamodbav:"requestedBy"`
	Instructions  string `json:"instructions" dynamodbav:"instructions"`
}

func (i *Instructions) DDBKeys() (ddb.Keys, error) {
	keys := ddb.Keys{
		PK: keys.AccessRequestGroupTargetInstructions.PK1,
		SK: keys.AccessRequestGroupTargetInstructions.SK1(i.GroupTargetID, i.RequestedBy),
	}
	return keys, nil
}
