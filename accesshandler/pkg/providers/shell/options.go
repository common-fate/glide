package shell

import (
	"context"
	"encoding/json"
	"net/http"
	"path"

	"github.com/common-fate/granted-approvals/accesshandler/pkg/providers"
	"github.com/common-fate/granted-approvals/accesshandler/pkg/types"
	"go.uber.org/zap"
)

// List options for arg
func (p *Provider) Options(ctx context.Context, arg string) (*types.ArgOptionsResponse, error) {
	switch arg {
	case "service":
		log := zap.S().With("arg", arg)
		log.Info("getting relay service options")
		res, err := p.getServices(ctx)
		if err != nil {
			return nil, err
		}
		var opts types.ArgOptionsResponse
		for _, s := range res.Services {
			opts.Options = append(opts.Options, types.Option{Value: s.ID, Label: s.ID})
		}
		return &opts, nil
	}
	return nil, &providers.InvalidArgumentError{Arg: arg}
}

func (p *Provider) getServices(ctx context.Context) (*listServicesResponse, error) {
	u := p.adminURL
	u.Path = path.Join(u.Path, "services")

	zap.S().Info("fetching services from relay", "url", u.String())

	res, err := http.Get(u.String())
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	var lsr listServicesResponse
	err = json.NewDecoder(res.Body).Decode(&lsr)
	if err != nil {
		return nil, err
	}
	return &lsr, nil
}

type listServicesResponse struct {
	Services []serviceResponse `json:"services"`
}

type serviceResponse struct {
	ID string `json:"id"`
}
