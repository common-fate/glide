package flask

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/ecs"
	"github.com/common-fate/granted-approvals/accesshandler/pkg/types"
)

// Options list the argument options for the provider
func (p *Provider) Options(ctx context.Context, arg string) ([]types.Option, error) {
	switch arg {
	case "taskdefinitionfamily":
		opts := []types.Option{}
		hasMore := true
		var nextToken *string

		for hasMore {

			tasks, err := p.ecsClient.ListTasks(ctx, &ecs.ListTasksInput{Cluster: aws.String(p.ecsClusterARN.Get()), NextToken: nextToken})
			if err != nil {
				return []types.Option{}, err
			}

			for _, t := range tasks.TaskArns {
				describedTask, err := p.ecsClient.DescribeTasks(ctx, &ecs.DescribeTasksInput{
					Tasks:   []string{t},
					Cluster: aws.String(p.ecsClusterARN.Get()),
				})
				if err != nil {
					return []types.Option{}, err
				}
				describedTasksDefn, err := p.ecsClient.DescribeTaskDefinition(ctx, &ecs.DescribeTaskDefinitionInput{
					TaskDefinition: describedTask.Tasks[0].TaskDefinitionArn,
				})
				if err != nil {
					return []types.Option{}, err
				}
				opts = append(opts, types.Option{Label: *describedTasksDefn.TaskDefinition.Family, Value: *describedTasksDefn.TaskDefinition.Family})
			}
			//exit the pagination
			nextToken = tasks.NextToken
			hasMore = nextToken != nil

		}

		return opts, nil
	}
	return nil, nil

}
