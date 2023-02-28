package pdk

import "encoding/json"

type payloadType string

const (
	payloadTypeGrant         payloadType = "grant"
	payloadTypeRevoke        payloadType = "revoke"
	payloadTypeDescribe      payloadType = "describe"
	payloadTypeLoadResources payloadType = "loadResources"
)

type payload struct {
	Type payloadType `json:"type"`
	Data any         `json:"data,omitempty"`
}

func (p payload) Marshal() ([]byte, error) {
	return json.Marshal(p)
}

type Target struct {
	// Mode is defines which behaviour of the provider to use, e.g SSO or Group
	// The modes are defined by the provider schema, and each deployment is registered with its mode configuration in the database
	Mode      string            `json:"kind"`
	Arguments map[string]string `json:"arguments"`
}

// NewDefaultModeTarget creates a target from arguments for the Default mode
// providers will eventually support multiple modes, until then
// All providers use the "Default" mode
func NewDefaultModeTarget(args map[string]string) Target {
	return Target{
		Mode:      "Default",
		Arguments: args,
	}
}

type grantData struct {
	Subject string `json:"subject"`
	Target  Target `json:"target"`
}

func NewGrantEvent(subject string, target Target) payload {
	return payload{
		Type: payloadTypeGrant,
		Data: grantData{
			Subject: subject,
			Target:  target,
		},
	}
}

type revokeData struct {
	Subject string `json:"subject"`
	Target  Target `json:"target"`
}

func NewRevokeEvent(subject string, target Target) payload {
	return payload{
		Type: payloadTypeRevoke,
		Data: revokeData{
			Subject: subject,
			Target:  target,
		},
	}
}

type loadResourceData struct {
	Name string      `json:"name"`
	Ctx  interface{} `json:"ctx"`
}

func NewLoadResourcesEvent(name string, ctx interface{}) payload {
	return payload{
		Type: payloadTypeLoadResources,
		Data: loadResourceData{
			Name: name,
			Ctx:  ctx,
		},
	}
}

func NewProviderDescribeEvent() payload {
	return payload{
		Type: payloadTypeDescribe,
	}
}
