package targetsvc

import (
	"context"
	"errors"
	"strings"

	"github.com/common-fate/common-fate/pkg/storage"
	"github.com/common-fate/common-fate/pkg/target"

	"github.com/common-fate/common-fate/pkg/types"
	"github.com/common-fate/ddb"
	"go.uber.org/zap"
)

type Provider struct {
	Publisher string
	Name      string
	Version   string
}

func SplitProviderString(s string) (Provider, error) {
	splitversion := strings.Split(s, "@")
	if len(splitversion) != 2 {
		return Provider{}, errors.New("target schema given in incorrect format")
	}

	splitname := strings.Split(splitversion[0], "/")
	if len(splitname) != 2 {
		return Provider{}, errors.New("target schema given in incorrect format")
	}
	p := Provider{
		Publisher: splitname[0],
		Name:      splitname[1],
		Version:   splitversion[1],
	}
	return p, nil
}

func (s *Service) CreateGroup(ctx context.Context, req types.CreateTargetGroupRequest) (*target.Group, error) {
	log := zap.S()

	q := &storage.GetTargetGroup{
		ID: req.ID,
	}

	_, err := s.DB.Query(ctx, q)
	if err == nil {
		return nil, ErrTargetGroupIdAlreadyExists
	}
	if err != nil && err != ddb.ErrNoItems {
		return nil, err
	}
	//look up target schema for the provider version
	provider, err := SplitProviderString(req.TargetSchema)
	if err != nil {
		return nil, err
	}
	result, err := s.ProviderRegistryClient.GetProviderWithResponse(ctx, provider.Publisher, provider.Name, provider.Version)
	if err != nil {
		return nil, err
	}

	switch result.StatusCode() {
	case 404:
		return nil, ddb.ErrNoItems
	case 500:
		return nil, errors.New(string(result.Body))

	}

	now := s.Clock.Now()
	group := target.Group{
		ID: req.ID,
		// The default mode here is a placeholder in our API until multi mode providers are supported fully by the framework
		// until it is changed, providers will always return the Default mode
		TargetSchema: target.GroupTargetSchema{From: req.TargetSchema, Schema: result.JSON200.Schema.Target.AdditionalProperties["Default"].Schema},
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
