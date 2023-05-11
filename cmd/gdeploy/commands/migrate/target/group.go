package target

import (
	"time"

	"github.com/common-fate/provider-registry-sdk-go/pkg/providerregistrysdk"
)

type Group struct {
	//user defined e.g. 'okta'
	ID string `json:"id" dynamodbav:"id"`

	// From is a reference to the provider and kind from the registry
	// that the target group was created from
	From From `json:"from" dynamodbav:"from"`

	// Schema is denomalised and saved here for efficiency
	Schema providerregistrysdk.Target `json:"schema" dynamodbav:"schema"`

	// reference to the SVG icon for the target group
	Icon string `json:"icon" dynamodbav:"icon"`

	CreatedAt time.Time `json:"createdAt" dynamodbav:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt" dynamodbav:"updatedAt"`
}

type From struct {
	Publisher string `json:"publisher" dynamodbav:"publisher"`
	Name      string `json:"name" dynamodbav:"name"`
	Version   string `json:"version" dynamodbav:"version"`
	Kind      string `json:"kind" dynamodbav:"kind"`
}
