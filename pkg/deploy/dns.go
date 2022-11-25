package deploy

import (
	"context"
	"fmt"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/route53"
	"github.com/aws/aws-sdk-go-v2/service/route53/types"
	"github.com/common-fate/common-fate/pkg/cfaws"
	"go.uber.org/zap"
)

func (c *Config) SetDNSRecord(ctx context.Context) error {
	cfg, err := cfaws.ConfigFromContextOrDefault(ctx)
	if err != nil {
		return err
	}
	o, err := c.LoadOutput(ctx)
	if err != nil {
		return err
	}

	name, err := c.GetDevStageName()
	if err != nil {
		return err
	}

	comment := fmt.Sprintf("Common Fate deployment %s", name)

	domain := fmt.Sprintf("%s.test.granted.run", name)

	client := route53.NewFromConfig(cfg)

	// the hosted zone ID is currently hardcoded to Common Fate's hosted zone for 'test.granted.run'.
	hostedZoneID := "Z0934863WRQ8PRDX856V"

	input := &route53.ChangeResourceRecordSetsInput{
		HostedZoneId: &hostedZoneID,
		ChangeBatch: &types.ChangeBatch{
			Comment: &comment,
			Changes: []types.Change{
				{
					Action: types.ChangeActionUpsert,
					ResourceRecordSet: &types.ResourceRecordSet{
						Name: &domain,
						Type: types.RRTypeCname,
						TTL:  aws.Int64(300),
						ResourceRecords: []types.ResourceRecord{
							{
								Value: &o.CloudFrontDomain,
							},
						},
					},
				},
			},
		},
	}

	zap.S().Infow("upserting route53 record", "input", input)
	_, err = client.ChangeResourceRecordSets(ctx, input)
	return err
}
