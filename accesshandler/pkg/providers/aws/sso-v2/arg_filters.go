package ssov2

import (
	"context"
	"strings"

	"github.com/aws/aws-sdk-go-v2/aws"
	resourcegroupstaggingapitypes "github.com/aws/aws-sdk-go-v2/service/resourcegroupstaggingapi/types"
	"github.com/common-fate/granted-approvals/accesshandler/pkg/providers"
)

func (p *Provider) ArgOptionGroupValues(ctx context.Context, argId string, groupID string, groupValues []string) ([]string, error) {
	switch argId {
	case "accountId":
		switch groupID {
		case "organizationalUnit":
			var accountIDs []string
			for _, groupValue := range groupValues {
				accounts, err := p.listChildAccountsForParent(ctx, groupValue)
				if err != nil {
					return nil, err
				}
				for _, account := range accounts {
					accountIDs = append(accountIDs, aws.ToString(account.Id))
				}
			}
			return accountIDs, nil
		case "tag":
			var tags []resourcegroupstaggingapitypes.TagFilter
			for _, gv := range groupValues {
				kv := strings.SplitN(gv, ":", 1)
				if len(kv) != 2 {
					return nil, &providers.InvalidGroupValueError{GroupID: groupID, GroupValue: gv}
				}
				tags = append(tags, resourcegroupstaggingapitypes.TagFilter{
					Key:    aws.String(kv[0]),
					Values: []string{kv[1]},
				})
			}
			return p.listAccountsWithTag(ctx, tags)
		default:
			return nil, &providers.InvalidGroupIDError{GroupID: groupID}
		}
	default:
		return nil, &providers.InvalidArgumentError{Arg: argId}
	}
}
