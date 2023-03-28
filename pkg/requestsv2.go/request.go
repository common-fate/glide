package requestsv2

import (
	"os/user"

	"github.com/common-fate/common-fate/pkg/types"
	"github.com/common-fate/provider-registry-sdk-go/pkg/providerregistrysdk"
)

type TargetFrom struct {
	Publisher string `json:"publisher" dynamodbav:"publisher"`
	Name      string `json:"name" dynamodbav:"name"`
	Version   string `json:"version" dynamodbav:"version"`
	Kind      string `json:"kind" dynamodbav:"kind"`
}

type Entitlement struct {
	ID           string                                     `json:"id" dynamodbav:"id"`
	Provider     TargetFrom                                 `json:"provider" dynamodbav:"provider"`
	Description  string                                     `json:"description" dynamodbav:"description"`
	OptionSchema map[string]providerregistrysdk.TargetField `json:"optionSchema" dynamodbav:"optionSchema"`
	User         string                                     `json:"user" dynamodbav:"user"`
}

// Status of an Access Request.
type Status string

const (
	APPROVED  Status = "APPROVED"
	DECLINED  Status = "DECLINED"
	CANCELLED Status = "CANCELLED"
	PENDING   Status = "PENDING"
)

type Grantv2 struct {
	ID          string      `json:"id" dynamodbav:"id"`
	User        user.User   `json:"user" dynamodbav:"user"`
	Entitlement Entitlement `json:"entitlement" dynamodbav:"entitlement"`
	Status      Status      `json:"status" dynamodbav:"status"`
}

type Option struct {
	Value       string     `json:"value" dynamodbav:"value"`
	Label       string     `json:"label" dynamodbav:"label"`
	Description *string    `json:"description" dynamodbav:"description"`
	Provider    TargetFrom `json:"provider" dynamodbav:"provider"`
}

type Requestv2 struct {
	// ID is a read-only field after the request has been created.
	ID      string         `json:"id" dynamodbav:"id"`
	Groups  []AccessGroup  `json:"groups" dynamodbav:"groups"`
	Context RequestContext `json:"context" dynamodbav:"context"`
}

type RequestContext struct {
	Purpose  string `json:"purpose" dynamodbav:"purpose"`
	Metadata string `json:"metadata" dynamodbav:"metadata"`
}

type AccessGroup struct {
	// ID is a read-only field after the request has been created.
	ID              string                `json:"id" dynamodbav:"id"`
	Grants          []Grantv2             `json:"grants" dynamodbav:"grants"`
	TimeConstraints types.TimeConstraints `json:"timeConstraints" dynamodbav:"timeConstraints"`
	Approval        string                `json:"Approval" dynamodbav:"Approval"`
}
