package deploy

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/feature/s3/manager"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/common-fate/granted-approvals/pkg/cfaws"
	"github.com/common-fate/granted-approvals/pkg/clio"
	"go.uber.org/zap"
)

type Manifest struct {
	// Version is the version of the manifest itself. Used for forwards-compatibility.
	Version                 int    `json:"manifestVersion"`
	LatestDeploymentVersion string `json:"latestDeploymentVersion"`
}

// PublishManifest updates the manifest.json file in the release bucket.
func PublishManifest(ctx context.Context, releaseBucket, version string) error {
	m := Manifest{
		LatestDeploymentVersion: version,
		Version:                 1,
	}
	mJSON, err := json.Marshal(m)
	if err != nil {
		return err
	}

	cfg, err := cfaws.ConfigFromContextOrDefault(ctx)
	if err != nil {
		return err
	}

	zap.S().Infow("writing manifest", "bucket", releaseBucket, "manifest", m)

	client := s3.NewFromConfig(cfg)
	_, err = client.PutObject(ctx, &s3.PutObjectInput{
		Bucket: &releaseBucket,
		Key:    aws.String("manifest.json"),
		Body:   bytes.NewBuffer(mJSON),
	})
	return err
}

// GetManifest retrieves the manifest.json file for the current deployment region
func GetManifest(ctx context.Context, region string) (Manifest, error) {
	cfg, err := cfaws.ConfigFromContextOrDefault(ctx)
	if err != nil {
		return Manifest{}, err
	}
	client := s3.NewFromConfig(cfg)

	buffer := manager.NewWriteAtBuffer([]byte{})

	bucket := fmt.Sprintf("granted-releases-%s", region)
	key := "manifest.json"

	clio.Debug("fetching manifest, bucket=%s key=%s", bucket, key)

	downloader := manager.NewDownloader(client)
	_, err = downloader.Download(ctx, buffer, &s3.GetObjectInput{
		Bucket: &bucket,
		Key:    &key,
	})
	if err != nil {
		return Manifest{}, err
	}

	var m Manifest
	err = json.Unmarshal(buffer.Bytes(), &m)
	if err != nil {
		return Manifest{}, err
	}
	return m, nil
}
