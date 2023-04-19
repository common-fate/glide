package requests

import (
	"time"

	"github.com/common-fate/common-fate/pkg/identity"
	"github.com/common-fate/common-fate/pkg/types"
)

type RequestContext struct {
	Purpose  string `json:"purpose" dynamodbav:"purpose"`
	Metadata string `json:"metadata" dynamodbav:"metadata"`

	Reason *string `json:"reason" dynamodbav:"reason"`
}

func (c *RequestContext) ToAPI() types.RequestContext {
	return types.RequestContext{
		Purpose: struct {
			Reason string "json:\"reason\""
		}{*c.Reason},
	}
}

// // RequestData is information provided by the user when they make the request,
// // through filling in form fields in the web application.
type RequestData struct {
	Reason *string `json:"reason,omitempty" dynamodbav:"reason,omitempty"`
}

// // WithNow allows you to override the now time used by getInterval
func WithNow(t time.Time) func(o *GetIntervalOpts) {
	return func(o *GetIntervalOpts) { o.Now = t }
}

type Requestv2 struct {
	types.Request
	Context     RequestContext
	RequestedBy identity.User
}
