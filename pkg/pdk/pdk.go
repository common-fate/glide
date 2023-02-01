package pdk

import "encoding/json"

type payloadType string

const (
	payloadTypeGrant         payloadType = "grant"
	payloadTypeRevoke        payloadType = "revoke"
	payloadTypeSchema        payloadType = "schema"
	payloadTypeLoadResources payloadType = "loadResources"
)

type payload struct {
	Type payloadType `json:"type"`
	Data any         `json:"data,omitempty"`
}

func (p payload) Marshal() ([]byte, error) {
	return json.Marshal(p)
}

type grantData struct {
	Subject string            `json:"subject"`
	Target  map[string]string `json:"args"`
}

func NewGrantEvent(subject string, target map[string]string) payload {
	return payload{
		Type: payloadTypeGrant,
		Data: grantData{
			Subject: subject,
			Target:  target,
		},
	}
}

type revokeData struct {
	Subject string            `json:"subject"`
	Target  map[string]string `json:"args"`
}

func NewRevokeEvent(subject string, target map[string]string) payload {
	return payload{
		Type: payloadTypeRevoke,
		Data: revokeData{
			Subject: subject,
			Target:  target,
		},
	}
}

func NewSchemaEvent() payload {
	return payload{
		Type: payloadTypeSchema,
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
