package workflowsvc

import (
	"context"

	"github.com/benbjohnson/clock"
	"github.com/common-fate/common-fate/pkg/access"
	"github.com/common-fate/common-fate/pkg/gevent"
	"github.com/common-fate/ddb"
)

// //go:generate go run github.com/golang/mock/mockgen -destination=mocks/runtime.go -package=mocks . Runtime
type Runtime interface {
	// grant is expected to be asyncronous
	Grant(ctx context.Context, access_group access.GroupTarget) error
	// isForTargetGroup tells the runtime how to process the request
	// revoke is expected to be syncronous
	Revoke(ctx context.Context, grantID string) error
}

// //go:generate go run github.com/golang/mock/mockgen -destination=mocks/eventputter.go -package=mocks . EventPutter
type EventPutter interface {
	Put(ctx context.Context, detail gevent.EventTyper) error
}
type Service struct {
	Runtime  Runtime
	DB       ddb.Storage
	Clk      clock.Clock
	Eventbus EventPutter
}

func (s *Service) Grant(ctx context.Context, access_group access.GroupTarget, subject string) ([]access.GroupTarget, error) {
	// Contains logic for preparing a grant and emitting events

	// err := s.Runtime.Grant(ctx, access_group)
	// if err != nil {
	// 	return nil, err
	// }

	// //TODO: Grant created event here

	// return grants.Result, nil
	return nil, nil
}

// // Revoke attepmts to syncronously revoke access to a request
// // If it is successful, the request is updated in the database, and the updated request is returned from this method
// func (s *Service) Revoke(ctx context.Context, request requests.Requestv2, revokerID string, revokerEmail string) (*requests.Requestv2, error) {

// 	accessGroups := storage.ListAccessGroups{RequestID: request.ID}
// 	_, err := s.DB.Query(ctx, &accessGroups)
// 	if err != nil {
// 		return nil, err
// 	}

// 	for _, access_group := range accessGroups.Result {

// 		//get all grants

// 		grants := storage.ListGrantsV2{GroupID: access_group.ID}

// 		_, err = s.DB.Query(ctx, &grants)
// 		if err == ddb.ErrNoItems {
// 			return nil, ErrNoGrant
// 		}
// 		if err != nil {
// 			return nil, err
// 		}

// 		for _, grant := range grants.Result {

// 			//Cannot request to revoke/cancel grant if it is not active or pending (state function has been created and executed)
// 			canRevoke := grant.Status == types.GrantStatusACTIVE || grant.Status == types.GrantStatusPENDING

// 			if !canRevoke || grant.End.Before(s.Clk.Now()) {
// 				return nil, ErrGrantInactive
// 			}

// 			q := storage.GetAccessRule{
// 				ID: access_group.AccessRule.ID,
// 			}

// 			_, err := s.DB.Query(ctx, &q)
// 			if err != nil {
// 				return nil, err
// 			}
// 			err = s.Runtime.Revoke(ctx, grant.ID)
// 			if err != nil {
// 				return nil, err
// 			}

// 			// previousStatus := grant.Status
// 			grant.Status = types.GrantStatusREVOKED
// 			grant.UpdatedAt = s.Clk.Now()

// 			//todo: update request item
// 			// items, err := dbupdate.GetUpdateRequestItems(ctx, s.DB, request)
// 			// if err != nil {
// 			// 	return nil, err
// 			// }

// 			// //create a request event for audit loggging request change
// 			// requestEvent := access.NewGrantStatusChangeEvent(request.ID, grant.UpdatedAt, &revokerID, previousStatus, grant.Status)

// 			// items = append(items, &requestEvent)

// 			// err = s.DB.PutBatch(ctx, items...)
// 			// if err != nil {
// 			// 	return nil, err
// 			// }
// 			// // Emit an event for the grant revoke
// 			// // We have chosen to emit events from the Common Fate app for grant revocation rather than from the access handler because we are using a syncronous API.
// 			// // All effects from revoking will be implemented in this syncronous api rather than triggered from the events.
// 			// // So we update the grant status here and save the grant before emitting the event
// 			// err = s.Eventbus.Put(ctx, gevent.GrantRevoked{Grant: grant.ToAPI(), Actor: revokerID, RevokerEmail: revokerEmail})
// 			// if err != nil {
// 			// 	return nil, err
// 			// }
// 		}
// 	}

// 	return &request, nil
// }

// // prepareCreateGrantRequest prepares the data for requesting
// func (s *Service) prepareCreateGrantRequest(ctx context.Context, groupTarget access.GroupTarget) (types.CreateGrant, error) {

// 	start, end := requestTiming.GetInterval(requests.WithNow(s.Clk.Now()))

// 	req := types.CreateGrant{
// 		// Id:       CreateGrantIdHash(subject, iso8601.New(start).Time, accessRule.Target.TargetGroupID),
// 		Id:       types.NewGrantID(),
// 		Provider: accessRule.Target.TargetGroupID,
// 		With: types.CreateGrant_With{
// 			AdditionalProperties: make(map[string]string),
// 		},
// 		Subject: openapi_types.Email(subject),
// 		Start:   iso8601.New(start),
// 		End:     iso8601.New(end),
// 	}

// 	//todo: rework this to be used safely around the codebase
// 	for k, v := range target.Fields {
// 		req.With.AdditionalProperties[k] = v.Value.Value
// 	}
// 	// for k, v := range request.SelectedWith {
// 	// 	req.With.AdditionalProperties[k] = v.Value
// 	// }

// 	return req, nil
// }

// // Due to multiple grants being created from one access group I am opting to create dynamic ID's
// // This helps with testing workflow grant creation
// // func CreateGrantIdHash(subject string, startTime time.Time, target string) string {
// // 	h := sha256.New()
// // 	hashString := subject + strconv.Itoa(int(startTime.Unix())) + target

// // 	h.Write([]byte(hashString))

// // 	bs := h.Sum(nil)
// // 	return fmt.Sprintf("%x", bs)

// // }
