package policy

import (
	"encoding/json"
	"fmt"
	"time"
)

type Policy struct {
	// 2012-10-17 or 2008-10-17 old policies, do NOT use this for new policies
	Version    string      `json:"Version"`
	Id         *string     `json:"Id,omitempty"`
	Statements []Statement `json:"Statement"`
}

func (p Policy) String() string {
	b, err := json.Marshal(p)
	if err != nil {
		return ""
	}
	return string(b)
}

type Statement struct {
	Sid          string           `json:"Sid,omitempty"`          // statement ID, service specific
	Effect       string           `json:"Effect"`                 // Allow or Deny
	Principal    map[string]Value `json:"Principal,omitempty"`    // principal that is allowed or denied
	NotPrincipal map[string]Value `json:"NotPrincipal,omitempty"` // exception to a list of principals
	Action       Value            `json:"Action,omitempty"`       // allowed or denied action
	NotAction    Value            `json:"NotAction,omitempty"`    // matches everything except
	Resource     Value            `json:"Resource,omitempty"`     // object or objects that the statement covers
	NotResource  Value            `json:"NotResource,omitempty"`  // matches everything except
	Condition    *ConditionEntry  `json:"Condition,omitempty"`    // conditions for when a policy is in effect
}

type ConditionEntry struct {
	DateGreaterThan AWSTime
	DateLessThan    AWSTime
}

type AWSTime struct {
	Time time.Time `json:"aws:CurrentTime"`
}

// AWS allows string or []string as value, we convert everything to []string to avoid casting
type Value []string

func (value *Value) UnmarshalJSON(b []byte) error {

	var raw interface{}
	err := json.Unmarshal(b, &raw)
	if err != nil {
		return err
	}

	var p []string
	//  value can be string or []string, convert everything to []string
	switch v := raw.(type) {
	case string:
		p = []string{v}
	case []interface{}:
		var items []string
		for _, item := range v {
			items = append(items, fmt.Sprintf("%v", item))
		}
		p = items
	default:
		return fmt.Errorf("invalid %s value element: allowed is only string or []string", value)
	}

	*value = p
	return nil
}

func AddExpiryCondition(p *Policy, expiresAt time.Time) {
	for i := range p.Statements {
		if p.Statements[i].Condition == nil {
			p.Statements[i].Condition = &ConditionEntry{}
		}

		p.Statements[i].Condition.DateLessThan = AWSTime{
			Time: expiresAt,
		}
	}
}
