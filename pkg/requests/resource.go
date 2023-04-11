package requests

import (
	"fmt"

	"github.com/common-fate/common-fate/pkg/storage/keys"
	"github.com/common-fate/common-fate/pkg/types"
	"github.com/common-fate/ddb"
)

// Status of an Access Request.
type Status string

const (
	APPROVED  Status = "APPROVED"
	DECLINED  Status = "DECLINED"
	CANCELLED Status = "CANCELLED"
	PENDING   Status = "PENDING"
)

type TargetFrom struct {
	Publisher string `json:"publisher" dynamodbav:"publisher"`
	Name      string `json:"name" dynamodbav:"name"`
	Version   string `json:"version" dynamodbav:"version"`
	Kind      string `json:"kind" dynamodbav:"kind"`
}

type ResourceOption struct {
	ID          string     `json:"id" dynamodbav:"id"`
	Type        string     `json:"type" dynamodbav:"type"`
	Value       string     `json:"value" dynamodbav:"value"`
	Label       string     `json:"label" dynamodbav:"label"`
	Description *string    `json:"description" dynamodbav:"description"`
	Provider    TargetFrom `json:"provider" dynamodbav:"provider"`
	TargetGroup string     `json:"targetGroup" dynamodbav:"targetGroup"`
	AccessRules []string   `json:"accessRules" dynamodbav:"accessRules"`
	RelatedTo   []string   `json:"childOf" dynamodbav:"childOf"`
}

func (o *TargetFrom) GetTargetFromString() string {
	return fmt.Sprintf("%s#%s#%s#%s", o.Kind, o.Publisher, o.Name, o.Version)
}

func (i *ResourceOption) DDBKeys() (ddb.Keys, error) {
	keys := ddb.Keys{
		PK: keys.OptionsV2.PK1(i.Label),
		SK: keys.OptionsV2.SK1(i.Provider.GetTargetFromString(), i.Value),
	}
	return keys, nil
}

func (e *ResourceOption) ToAPI() types.Resource {
	return types.Resource{
		Label: e.Label,
		Value: e.Value,
	}
}
