package pdk

import "encoding/json"

type payloadType string

const (
	payloadTypeGrant  payloadType = "grant"
	payloadTypeRevoke payloadType = "revoke"
	payloadTypeSchema payloadType = "schema"
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
