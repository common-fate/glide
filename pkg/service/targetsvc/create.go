package targetsvc

import (
	"context"

	"fmt"
	"net/http"
	"strings"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/common-fate/common-fate/pkg/storage"
	"github.com/common-fate/common-fate/pkg/target"
	"github.com/pkg/errors"

	"github.com/common-fate/common-fate/pkg/types"
	"github.com/common-fate/ddb"
	"go.uber.org/zap"
)

type Provider struct {
	Publisher string
	Name      string
	Version   string
	Kind      *string
}

func SplitProviderString(s string) (Provider, error) {
	splitversion := strings.Split(s, "@")
	if len(splitversion) != 2 {
		return Provider{}, errors.New("target schema given in incorrect format")
	}
	var kind *string
	splitKind := strings.Split(splitversion[1], "/")
	if len(splitKind) == 2 {
		kind = aws.String(splitKind[1])
	}

	splitname := strings.Split(splitversion[0], "/")
	if len(splitname) != 2 {
		return Provider{}, errors.New("target schema given in incorrect format")
	}
	p := Provider{
		Publisher: splitname[0],
		Name:      splitname[1],
		Version:   splitKind[0],
		Kind:      kind,
	}
	return p, nil
}

func (s *Service) CreateGroup(ctx context.Context, req types.CreateTargetGroupRequest) (*target.Group, error) {
	log := zap.S()

	q := &storage.GetTargetGroup{
		ID: req.Id,
	}

	_, err := s.DB.Query(ctx, q)
	if err == nil {
		return nil, ErrTargetGroupIdAlreadyExists
	}
	if err != nil && err != ddb.ErrNoItems {
		return nil, err
	}
	// //look up target schema for the provider version
	// provider, err := SplitProviderString(req.TargetSchema)
	// if err != nil {
	// 	return nil, err
	// }

	if req.From.Kind == "" {
		return nil, ErrKindIsRequired
	}
	response, err := s.ProviderRegistryClient.GetProviderWithResponse(ctx, req.From.Publisher, req.From.Name, req.From.Version)
	if err != nil {
		return nil, err
	}

	switch response.StatusCode() {
	case http.StatusOK:
	case http.StatusNotFound:
		return nil, ErrProviderNotFoundInRegistry
	case http.StatusInternalServerError:
		return nil, errors.Wrap(fmt.Errorf(response.JSON500.Error), "received 500 error from registry service when fetching provider")
	default:
		return nil, fmt.Errorf("unhandled response code received from registry service when querying for a provider status Code: %d Body: %s", response.StatusCode(), string(response.Body))
	}

	targets := response.JSON200.Schema.Targets
	if targets == nil {
		return nil, errors.New("provider does not provide any targets")
	}

	schema, ok := (*targets)[req.From.Kind]

	if !ok {
		return nil, ErrProviderDoesNotImplementKind
	}

	var icon string
	if response.JSON200.Meta != nil && response.JSON200.Meta.Icon != nil {
		icon = *response.JSON200.Meta.Icon
	}

	now := s.Clock.Now()
	group := target.Group{
		ID:        req.Id,
		Schema:    schema,
		From:      target.FromFieldFromAPI(req.From),
		Icon:      icon,
		CreatedAt: now,
		UpdatedAt: now,
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
