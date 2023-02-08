package config

type Config struct {
	Host              string `env:"COMMONFATE_HOST,default=0.0.0.0:8080"`
	LogLevel          string `env:"LOG_LEVEL,default=info"`
	DynamoTable       string `env:"COMMONFATE_TABLE_NAME,required"`
	CognitoUserPoolID string `env:"COMMONFATE_COGNITO_USER_POOL_ID,required"`
	Region            string `env:"AWS_REGION,required"`
	AdminGroup        string `env:"COMMONFATE_ADMIN_GROUP,required"`
	FrontendURL       string `env:"COMMONFATE_FRONTEND_URL,required"`
	AccessHandlerURL  string `env:"COMMONFATE_ACCESS_HANDLER_URL,default=http://0.0.0.0:9092"`
	GovernanceURL     string `env:"COMMONFATE_GOVERNANCE_URL,default=0.0.0.0:8889"`
	RunAccessHandler  bool   `env:"COMMONFATE_RUN_ACCESS_HANDLER,default=true"`
	MockAccessHandler bool   `env:"COMMONFATE_MOCK_ACCESS_HANDLER,default=false"`
	SentryDSN         string `env:"COMMONFATE_SENTRY_DSN"`
	EventBusArn       string `env:"COMMONFATE_EVENT_BUS_ARN,required"`
	EventBusSource    string `env:"COMMONFATE_EVENT_BUS_SOURCE,required"`
	IdpProvider       string `env:"COMMONFATE_IDENTITY_PROVIDER,required"`
	DeploymentSuffix  string `env:"COMMONFATE_DEPLOYMENT_SUFFIX"`
	// This should be an instance of deploy.FeatureMap which is a specific json format for this
	// Use deploy.UnmarshalFeatureMap to unmarshal this data into a FeatureMap
	IdentitySettings              string `env:"COMMONFATE_IDENTITY_SETTINGS,default={}"`
	PaginationKMSKeyARN           string `env:"COMMONFATE_PAGINATION_KMS_KEY_ARN,required"`
	AccessHandlerExecutionRoleARN string `env:"COMMONFATE_ACCESS_HANDLER_EXECUTION_ROLE_ARN,required"`
	RemoteConfigURL               string `env:"COMMONFATE_ACCESS_REMOTE_CONFIG_URL"`
	RemoteConfigHeaders           string `env:"COMMONFATE_REMOTE_CONFIG_HEADERS"`

	NoAuthEmail string `env:"NO_AUTH_EMAIL"`
}

type NotificationsConfig struct {
	LogLevel    string `env:"LOG_LEVEL,default=info"`
	DynamoTable string `env:"COMMONFATE_TABLE_NAME,required"`
	FrontendURL string `env:"COMMONFATE_FRONTEND_URL,required"`
	// This should be an instance of deploy.FeatureMap which is a specific json format for this
	// Use deploy.UnmarshalFeatureMap to unmarshal this data into a FeatureMap
	NotificationsConfig string `env:"COMMONFATE_NOTIFICATIONS_SETTINGS,default={}"`
	RemoteConfigURL     string `env:"COMMONFATE_ACCESS_REMOTE_CONFIG_URL"`
	RemoteConfigHeaders string `env:"COMMONFATE_REMOTE_CONFIG_HEADERS"`
}

type EventHandlerConfig struct {
	LogLevel    string `env:"LOG_LEVEL,default=info"`
	DynamoTable string `env:"COMMONFATE_TABLE_NAME,required"`
}

type SyncConfig struct {
	TableName   string `env:"COMMONFATE_TABLE_NAME,required"`
	IdpProvider string `env:"COMMONFATE_IDENTITY_PROVIDER,required"`
	UserPoolId  string `env:"COMMONFATE_COGNITO_USER_POOL_ID,required"`
	LogLevel    string `env:"LOG_LEVEL,default=info"`
	// This should be an instance of deploy.FeatureMap which is a specific json format for this
	// Use deploy.UnmarshalFeatureMap to unmarshal this data into a FeatureMap
	IdentitySettings string `env:"COMMONFATE_IDENTITY_SETTINGS,default={}"`
}
type CacheSyncConfig struct {
	TableName        string `env:"COMMONFATE_TABLE_NAME,required"`
	LogLevel         string `env:"LOG_LEVEL,default=info"`
	Region           string `env:"AWS_REGION,required"`
	AccessHandlerURL string `env:"COMMONFATE_ACCESS_HANDLER_URL,default=http://0.0.0.0:9092"`
}

type FrontendDeployerConfig struct {
	LogLevel                             string `env:"LOG_LEVEL,default=info"`
	Region                               string `env:"AWS_REGION,required"`
	CFReleasesBucket                     string `env:"CF_RELEASES_BUCKET,required"`
	CFReleasesFrontendAssetsObjectPrefix string `env:"CF_RELEASES_FRONTEND_ASSET_OBJECT_PREFIX,required"`
	FrontendBucket                       string `env:"COMMONFATE_FRONTEND_BUCKET,required"`
	UserPoolID                           string `env:"COMMONFATE_COGNITO_USER_POOL_ID,required"`
	CognitoClientID                      string `env:"COMMONFATE_COGNITO_CLIENT_ID,required"`
	UserPoolDomain                       string `env:"COMMONFATE_USER_POOL_DOMAIN,required"`
	FrontendDomain                       string `env:"COMMONFATE_FRONTEND_DOMAIN,required"`
	CloudFrontDistributionID             string `env:"COMMONFATE_CLOUDFRONT_DISTRIBUTION_ID,required"`
	APIURL                               string `env:"COMMONFATE_API_URL,required"`
}
