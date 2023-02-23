package target

import (
	"github.com/common-fate/common-fate/pkg/storage/keys"
	"github.com/common-fate/common-fate/pkg/types"
	"github.com/common-fate/ddb"
)

type Route struct {
	Group       string       `json:"group" dynamodbav:"group"`
	Handler     string       `json:"handler" dynamodbav:"handler"`
	Mode        string       `json:"mode" dynamodbav:"mode"`
	Priority    int          `json:"priority" dynamodbav:"priority"`
	Valid       bool         `json:"valid" dynamodbav:"valid"`
	Diagnostics []Diagnostic `json:"diagnostics" dynamodbav:"diagnostics"`
}

type Diagnostic struct {
	Level   string `json:"level" dynamodbav:"level"`
	Code    string `json:"code" dynamodbav:"code"`
	Message string `json:"message" dynamodbav:"message"`
	// Allows diagnostics to be grouped by an arbitrary key, for use in UI
	// This can be set in a response from a provider
	GroupBy *string `json:"groupBy" dynamodbav:"groupBy"`
}

func (r *Route) DDBKeys() (ddb.Keys, error) {
	keys := ddb.Keys{
		PK:     keys.TargetRoute.PK1,
		SK:     keys.TargetRoute.SK1(r.Group, r.Handler, r.Mode),
		GSI1PK: keys.TargetRoute.GSI1PK(r.Group),
		GSI1SK: keys.TargetRoute.GSI1SK(r.Valid, r.Priority),
	}
	return keys, nil
}

func (r *Route) ToAPI() types.TargetRoute {
	diagnostics := make([]types.Diagnostic, len(r.Diagnostics))
	for i, d := range r.Diagnostics {
		diagnostics[i] = types.Diagnostic{
			Code:    d.Code,
			Level:   d.Level,
			Message: d.Message,
		}
	}
	return types.TargetRoute{
		TargetGroupId: r.Group,
		HandlerId:     r.Handler,
		Mode:          r.Mode,
		Priority:      r.Priority,
		Valid:         r.Valid,
		Diagnostics:   diagnostics,
	}
}
