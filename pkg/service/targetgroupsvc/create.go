package targetgroupsvc

import (
	"context"
	"errors"
	"strings"

	"github.com/common-fate/common-fate/pkg/targetgroup"
	"github.com/common-fate/common-fate/pkg/types"
	"go.uber.org/zap"
)

func (s *Service) CreateTargetGroup(ctx context.Context, req types.CreateTargetGroupRequest) (*targetgroup.TargetGroup, error) {
	log := zap.S()

	//look up target schema for the provider version

	//TODO: validate that the targetschema is something we are expecting

	splitKey := strings.Split(req.TargetSchema, "/")

	//the target schema we receive should be in the form team/provider/version and split into 3 keys
	if len(splitKey) != 3 {
		return nil, errors.New("target schema given in incorrect format")
	}

	resp, err := s.ProviderRegistryClient.GetProviderWithResponse(ctx, splitKey[0], splitKey[1], splitKey[2])
	if err != nil {
		return nil, err
	}
	now := s.Clock.Now()
	group := targetgroup.TargetGroup{
		ID:           req.ID,
		TargetSchema: targetgroup.GroupTargetSchema{From: req.TargetSchema, Schema: resp.JSON200.Schema.Target},
		CreatedAt:    now,
		UpdatedAt:    now,
	}
	//based on the target schema provider type set the Icon

	log.Debugw("saving target group", "group", group)
	// save the request.
	err = s.DB.Put(ctx, &group)
	if err != nil {
		return nil, err
	}
	return &group, nil
}
