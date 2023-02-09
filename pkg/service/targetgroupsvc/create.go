package targetgroupsvc

import (
	"context"
	"errors"
	"reflect"
	"strings"

	"github.com/common-fate/common-fate/pkg/storage"
	"github.com/common-fate/common-fate/pkg/targetgroup"
	"github.com/common-fate/common-fate/pkg/types"
	"github.com/common-fate/ddb"
	"go.uber.org/zap"
)

func (s *Service) CompareAndValidateProviderVersions(ctx context.Context, provider1 string, provider2 string) (bool, error) {
	splitKey := strings.Split(provider1, "/")

	//the target schema we receive should be in the form team/provider/version and split into 3 keys
	if len(splitKey) != 3 {
		return false, errors.New("target schema given in incorrect format")
	}

	provider1Resp, err := s.ProviderRegistryClient.GetProviderWithResponse(ctx, splitKey[0], splitKey[1], splitKey[2])
	if err != nil {
		return false, err
	}

	splitKey = strings.Split(provider2, "/")

	//the target schema we receive should be in the form team/provider/version and split into 3 keys
	if len(splitKey) != 3 {
		return false, errors.New("target schema given in incorrect format")
	}

	provider2Resp, err := s.ProviderRegistryClient.GetProviderWithResponse(ctx, splitKey[0], splitKey[1], splitKey[2])
	if err != nil {
		return false, err
	}

	return reflect.DeepEqual(provider1Resp.JSON200.Schema.Target, provider2Resp.JSON200.Schema.Target), nil

}

func (s *Service) CreateTargetGroupLink(ctx context.Context, req types.CreateTargetGroupLink, targetGroupId string) (*targetgroup.TargetGroup, error) {
	log := zap.S()

	//validate target group exists
	q := storage.GetTargetGroup{ID: targetGroupId}

	_, err := s.DB.Query(ctx, &q)

	if err == ddb.ErrNoItems {

		return nil, err
	}

	//lookup deployment
	p := storage.GetTargetGroupDeployment{ID: req.DeploymentId}

	_, err = s.DB.Query(ctx, &p)

	if err == ddb.ErrNoItems {
		return nil, err
	}

	//validate deployment target schema matches target group

	//update target group deployments to include new deployment
	q.Result.TargetDeployments = append(q.Result.TargetDeployments, targetgroup.DeploymentRegistration{
		ID:          p.Result.ID,
		Priority:    req.Priority,
		Diagnostics: p.Result.Diagnostics,
	})

	log.Debugw("Linking deployment to target group", "group", q.Result.ID)
	// save the request.
	err = s.DB.Put(ctx, &q.Result)
	if err != nil {
		return nil, err
	}

	return &q.Result, nil
}

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
