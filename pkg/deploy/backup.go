package deploy

import (
	"context"
	"fmt"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	ddbTypes "github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
)

func StartBackup(ctx context.Context, tableName string, backupName string) (*ddbTypes.BackupDetails, error) {

	cfg, err := config.LoadDefaultConfig(ctx)
	if err != nil {
		return nil, err
	}
	client := dynamodb.NewFromConfig(cfg)
	cbo, err := client.CreateBackup(ctx, &dynamodb.CreateBackupInput{
		BackupName: &backupName,
		TableName:  &tableName,
	})
	if err != nil {
		return nil, err
	}
	return cbo.BackupDetails, nil
}

func BackupStatus(ctx context.Context, backupARN string) (*ddbTypes.BackupDescription, error) {
	cfg, err := config.LoadDefaultConfig(ctx)
	if err != nil {
		return nil, err
	}
	client := dynamodb.NewFromConfig(cfg)
	b, err := client.DescribeBackup(ctx, &dynamodb.DescribeBackupInput{
		BackupArn: &backupARN,
	})
	if err != nil {
		return nil, err
	}
	return b.BackupDescription, nil
}

func BackupDetailsToString(b *ddbTypes.BackupDetails) string {
	if b == nil {
		return ""
	}
	var expiry string
	if b.BackupExpiryDateTime != nil {
		expiry = fmt.Sprintf("\nExpiry: %s", aws.ToTime(b.BackupExpiryDateTime).Local().Format(time.RFC3339))
	}
	return fmt.Sprintf("Backup: %s\nARN: %s\nStatus: %s%s", aws.ToString(b.BackupName), aws.ToString(b.BackupArn), b.BackupStatus, expiry)
}

func StartRestore(ctx context.Context, backupARN string, targetTableName string) (*ddbTypes.TableDescription, error) {
	cfg, err := config.LoadDefaultConfig(ctx)
	if err != nil {
		return nil, err
	}
	client := dynamodb.NewFromConfig(cfg)
	rbo, err := client.RestoreTableFromBackup(ctx, &dynamodb.RestoreTableFromBackupInput{
		BackupArn:       &backupARN,
		TargetTableName: &targetTableName,
	})
	if err != nil {
		return nil, err
	}
	return rbo.TableDescription, nil
}
func RestoreStatus(ctx context.Context, targetTableName string) (*ddbTypes.TableDescription, error) {
	cfg, err := config.LoadDefaultConfig(ctx)
	if err != nil {
		return nil, err
	}
	client := dynamodb.NewFromConfig(cfg)
	rbo, err := client.DescribeTable(ctx, &dynamodb.DescribeTableInput{
		TableName: &targetTableName,
	})

	if err != nil {
		return nil, err
	}
	return rbo.Table, nil
}

func RestoreSummaryToString(r *ddbTypes.RestoreSummary) string {
	if r == nil {
		return "No restoration in progress"
	}
	return fmt.Sprintf("Restoration in progress from backup: %s", aws.ToString(r.SourceBackupArn))
}
