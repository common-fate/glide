package backgroundtasks

import (
	"context"
	"net/http"
	"strings"
	"time"

	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/cloudformation"
	"github.com/aws/aws-sdk-go-v2/service/cloudformation/types"
	"github.com/common-fate/apikit/logger"
	cfTypes "github.com/common-fate/common-fate/pkg/types"
	"github.com/mitchellh/mapstructure"
)

type Service struct {
	CommonFate cfTypes.ClientWithResponsesInterface
}

func (s *Service) StartPollForDeploymentStatus(stackID string, provider cfTypes.ProviderV2) {
	go func() {
		ctx := context.Background()
		log := logger.Get(ctx)
		cfg, err := config.LoadDefaultConfig(ctx)
		if err != nil {
			log.Error(err)
			return
		}
		client := cloudformation.NewFromConfig(cfg)

		for {
			log.Infow("describing cloudformation stack", "stack", stackID)
			res, err := client.DescribeStacks(ctx, &cloudformation.DescribeStacksInput{StackName: &stackID})
			if err != nil {
				log.Error(err)
				return
			}
			log.Infow("cloudformation stacks", "result", res)
			if len(res.Stacks) != 1 {
				if err != nil {
					log.Error("expected 1 stack")
					return
				}
			}
			if strings.Contains(string(res.Stacks[0].StackStatus), "FAILED") || strings.Contains(string(res.Stacks[0].StackStatus), "ROLLBACK") {
				if err != nil {
					log.Error("stack failed or rolling back")
					return
				}
			}
			if res.Stacks[0].StackStatus == types.StackStatusCreateComplete || res.Stacks[0].StackStatus == types.StackStatusUpdateComplete {
				outputMap := make(map[string]string)
				for _, o := range res.Stacks[0].Outputs {
					outputMap[*o.OutputKey] = *o.OutputValue
				}
				type Output struct {
					ProviderFunctionARN     string `json:"ProviderFunctionARN"`
					ProviderFunctionRoleARN string `json:"ProviderFunctionRoleARN"`
				}
				var out Output
				decoder, err := mapstructure.NewDecoder(&mapstructure.DecoderConfig{TagName: "json", Result: &out})
				if err != nil {
					log.Error(err)
					return
				}
				err = decoder.Decode(outputMap)
				if err != nil {
					log.Error(err)
					return
				}
				updateRes, err := s.CommonFate.AdminUpdateProviderv2WithResponse(ctx, provider.Id, cfTypes.UpdateProviderV2{
					Alias:           provider.Alias,
					Status:          cfTypes.DEPLOYED,
					Version:         provider.Version,
					FunctionArn:     &out.ProviderFunctionARN,
					FunctionRoleArn: &out.ProviderFunctionRoleARN,
				})
				if err != nil {
					log.Error("stack failed or rolling back")
					return
				}
				if updateRes.StatusCode() != http.StatusOK {
					if err != nil {
						log.Errorw("provider update failed", "body", string(updateRes.Body))
						return
					}
				}
				log.Info("successfully updated status")

				return
			}
			time.Sleep(time.Second * 5)
		}
	}()

}
