// Package depid handles loading and saving deployment information.
package depid

import (
	"context"
	"os"

	"github.com/common-fate/analytics-go"
	"github.com/common-fate/ddb"
	"github.com/common-fate/granted-approvals/internal/build"
	"github.com/common-fate/granted-approvals/pkg/storage/keys"
	"github.com/common-fate/granted-approvals/pkg/types"
	"go.uber.org/zap"
)

type Deployment struct {
	ID         string `json:"id" dynamodbav:"id"`
	UserCount  int    `json:"userCount" dynamodbav:"userCount"`
	GroupCount int    `json:"groupCount" dynamodbav:"groupCount"`
	IDP        string `json:"idp" dynamodbav:"idp"`
}

func (d *Deployment) ToAnalytics() *analytics.DeploymentInfo {
	return &analytics.DeploymentInfo{
		ID:         d.ID,
		Version:    build.Version,
		Stage:      os.Getenv("CF_ANALYTICS_DEPLOYMENT_STAGE"),
		UserCount:  d.UserCount,
		GroupCount: d.GroupCount,
		IDP:        d.IDP,
	}
}

func (d *Deployment) DDBKeys() (ddb.Keys, error) {
	k := ddb.Keys{
		PK: keys.Deployment.PK1,
		SK: keys.Deployment.SK1,
	}
	return k, nil
}

type Loader struct {
	client ddb.Storage
	log    *zap.SugaredLogger
}

func New(client ddb.Storage, log *zap.SugaredLogger) *Loader {
	return &Loader{client: client, log: log.Named("depid")}
}

func (l *Loader) GetDeployment(ctx context.Context) (*Deployment, error) {
	return l.getOrCreateDeployment(ctx)
}

func (l *Loader) getOrCreateDeployment(ctx context.Context) (*Deployment, error) {
	var d Deployment
	_, err := l.client.Get(ctx, ddb.GetKey{PK: keys.Deployment.PK1, SK: keys.Deployment.SK1}, &d)
	if err == nil {
		l.log.Debugw("found existing deployment info", "deployment.id", d.ID)
		return &d, nil
	}
	if err != ddb.ErrNoItems {
		return nil, err
	}
	// if we get here, we got ddb.ErrNoItems.
	// this means the deployment isn't in the database, so provision it.
	d = Deployment{
		ID: types.NewDeploymentID(),
	}

	l.log.Infow("created deployment info", "deployment.id", d.ID)

	err = l.client.Put(ctx, &d)
	if err != nil {
		return nil, err
	}

	return &d, nil
}

type UserInfo struct {
	UserCount  int
	GroupCount int
}

func (l *Loader) SetUserInfo(ctx context.Context, ui UserInfo) (*Deployment, error) {
	d, err := l.getOrCreateDeployment(ctx)
	if err != nil {
		return nil, err
	}
	d.UserCount = ui.UserCount
	d.GroupCount = ui.GroupCount

	l.log.Infow("set deployment user info", "deployment.id", d.ID, "deployment.users", d.UserCount, "deployment.groups", d.GroupCount)

	err = l.client.Put(ctx, d)
	if err != nil {
		return nil, err
	}
	return d, nil
}
