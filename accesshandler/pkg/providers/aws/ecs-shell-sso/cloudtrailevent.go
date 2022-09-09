package ecsshellsso

import "time"

// Auto-generated since the AWS SDK wont return the complete type
type CloudTrailEvent struct {
	EventVersion string `json:"eventVersion"`
	UserIdentity struct {
		Type           string `json:"type"`
		PrincipalID    string `json:"principalId"`
		Arn            string `json:"arn"`
		AccountID      string `json:"accountId"`
		AccessKeyID    string `json:"accessKeyId"`
		SessionContext struct {
			SessionIssuer struct {
				Type        string `json:"type"`
				PrincipalID string `json:"principalId"`
				Arn         string `json:"arn"`
				AccountID   string `json:"accountId"`
				UserName    string `json:"userName"`
			} `json:"sessionIssuer"`
			WebIDFederationData struct {
			} `json:"webIdFederationData"`
			Attributes struct {
				CreationDate     time.Time `json:"creationDate"`
				MfaAuthenticated string    `json:"mfaAuthenticated"`
			} `json:"attributes"`
		} `json:"sessionContext"`
		InvokedBy string `json:"invokedBy"`
	} `json:"userIdentity"`
	EventTime         time.Time `json:"eventTime"`
	EventSource       string    `json:"eventSource"`
	EventName         string    `json:"eventName"`
	AwsRegion         string    `json:"awsRegion"`
	SourceIPAddress   string    `json:"sourceIPAddress"`
	UserAgent         string    `json:"userAgent"`
	RequestParameters struct {
		Target       string `json:"target"`
		DocumentName string `json:"documentName"`
		Parameters   struct {
			CloudWatchEncryptionEnabled []string `json:"cloudWatchEncryptionEnabled"`
			S3EncryptionEnabled         []string `json:"s3EncryptionEnabled"`
			CloudWatchLogGroupName      []string `json:"cloudWatchLogGroupName"`
			Command                     []string `json:"command"`
		} `json:"parameters"`
	} `json:"requestParameters"`
	ResponseElements struct {
		SessionID  string `json:"sessionId"`
		TokenValue string `json:"tokenValue"`
		StreamURL  string `json:"streamUrl"`
	} `json:"responseElements"`
	RequestID          string `json:"requestID"`
	EventID            string `json:"eventID"`
	ReadOnly           bool   `json:"readOnly"`
	EventType          string `json:"eventType"`
	ManagementEvent    bool   `json:"managementEvent"`
	RecipientAccountID string `json:"recipientAccountId"`
	EventCategory      string `json:"eventCategory"`
}
