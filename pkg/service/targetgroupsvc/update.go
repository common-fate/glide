package targetgroupsvc

import (
	"context"
	"strings"

	"github.com/common-fate/common-fate/pkg/storage"
	"github.com/common-fate/common-fate/pkg/targetgroup"
	"github.com/common-fate/common-fate/pkg/types"
	"github.com/common-fate/ddb"
	"go.uber.org/zap"
)

type UpdateOpts struct {
	TargetGroupID string
	UpdateRequest types.CreateTargetGroupRequest
}

func (s *Service) UpdateTargetGroup(ctx context.Context, req UpdateOpts) (*targetgroup.TargetGroup, error) {
	log := zap.S()

	//get target group for updating

	q := storage.GetTargetGroup{ID: req.TargetGroupID}

	_, err := s.DB.Query(ctx, &q)
	if err != nil && err != ddb.ErrNoItems {
		return nil, err
	}

	targetGroup := q.Result

	targetGroup.ID = req.UpdateRequest.ID

	if req.UpdateRequest.TargetSchema != targetGroup.TargetSchema.From {

		splitKey := strings.Split(req.UpdateRequest.TargetSchema, "/")

		resp, err := s.ProviderRegistryClient.GetProviderWithResponse(ctx, splitKey[0], splitKey[1], splitKey[2])
		if err != nil {
			return nil, err
		}

		targetGroup.TargetSchema.Schema = resp.JSON200.Schema.Target
		targetGroup.TargetSchema.From = req.UpdateRequest.TargetSchema
	}
	//look up target schema for the provider version

	//based on the target schema provider type set the Icon

	log.Debugw("updating target group", "group", targetGroup)
	// save the request.
	err = s.DB.Put(ctx, &targetGroup)
	if err != nil {
		return nil, err
	}
	return &targetGroup, nil
}
