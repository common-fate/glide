package action

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
	case "action":
		log := zap.S().With("arg", arg)
		log.Info("getting relay action options")
		res, err := p.getActions(ctx)
		if err != nil {
			return nil, err
		}
		var opts types.ArgOptionsResponse
		for _, s := range res.Actions {
			opts.Options = append(opts.Options, types.Option{Value: s.ID, Label: s.ID})
		}
		return &opts, nil
	}
	return nil, &providers.InvalidArgumentError{Arg: arg}
}

func (p *Provider) getActions(ctx context.Context) (*listActionsResponse, error) {
	u := p.adminURL
	u.Path = path.Join(u.Path, "actions")

	zap.S().Info("fetching actions from relay", "url", u.String())

	res, err := http.Get(u.String())
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	var lar listActionsResponse
	err = json.NewDecoder(res.Body).Decode(&lar)
	if err != nil {
		return nil, err
	}
	return &lar, nil
}

type listActionsResponse struct {
	Actions []actionResponse `json:"actions"`
}

type actionResponse struct {
	ID string `json:"id"`
}
