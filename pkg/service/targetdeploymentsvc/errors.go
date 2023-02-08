package targetdeploymentsvc

import "errors"

var (
	ErrTargetGroupDeploymentIdAlreadyExists = errors.New("target group deployment id already exists")
	ErrInvalidAwsAccountNumber              = errors.New("invalid aws account number")
)
