package shell

import (
	"context"

	"github.com/common-fate/granted-approvals/accesshandler/pkg/providers"
)

func (p *Provider) Operations() map[string]providers.Operation {
	return map[string]providers.Operation{
		"get-socket": {
			Execute: func(ctx context.Context, opts providers.OperationOpts) (map[string]any, error) {
				u := p.userURL
				q := u.Query()
				q.Set("user", opts.Subject)
				q.Set("accessRequest", opts.GrantID)
				u.RawQuery = q.Encode()
				res := map[string]any{
					"url": u.String(),
				}
				return res, nil
			},
		},
	}
}
