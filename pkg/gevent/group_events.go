package gevent

import "github.com/common-fate/common-fate/pkg/access"

//new AccessGroup Requests

const (
	AccessGroupGrant  = "accessGroup.grant"
	AccessGroupRevoke = "accessGroup.revoke"
	AccessGroupCancel = "accessGroup.cancel"
)

type AccessGroupGrantCreated struct {
	Group access.Group `json:"group"`
}

func (AccessGroupGrantCreated) EventType() string {
	return AccessGroupGrant
}

type AccessGroupGrantRevoked struct {
	Group access.Group `json:"group"`
}

func (AccessGroupGrantRevoked) EventType() string {
	return AccessGroupGrant
}

type AccessGroupGrantCancelled struct {
	Group access.Group `json:"group"`
}

func (AccessGroupGrantCancelled) EventType() string {
	return AccessGroupGrant
}
