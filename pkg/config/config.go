package config

type Config struct {
	Host              string `env:"APPROVALS_HOST,default=0.0.0.0:8080"`
	LogLevel          string `env:"LOG_LEVEL,default=info"`
	DynamoTable       string `env:"APPROVALS_TABLE_NAME,required"`
	CognitoUserPoolID string `env:"APPROVALS_COGNITO_USER_POOL_ID,required"`
	Region            string `env:"AWS_REGION,required"`
	AdminGroup        string `env:"APPROVALS_ADMIN_GROUP,required"`
	FrontendURL       string `env:"APPROVALS_FRONTEND_URL,required"`
	AccessHandlerURL  string `env:"ACCESS_HANDLER_URL,default=http://0.0.0.0:9092"`
	RunAccessHandler  bool   `env:"RUN_ACCESS_HANDLER,default=true"`
	MockAccessHandler bool   `env:"MOCK_ACCESS_HANDLER,default=false"`
	SentryDSN         string `env:"SENTRY_DSN"`
	EventBusArn       string `env:"EVENT_BUS_ARN,required"`
	EventBusSource    string `env:"EVENT_BUS_SOURCE,required"`
	IdpProvider       string `env:"IDENTITY_PROVIDER,required"`
	DeploymentSuffix  string `env:"DEPLOYMENT_SUFFIX"`
	// This should be an instance of deploy.FeatureMap which is a specific json format for this
	// Use deploy.UnmarshalFeatureMap to unmarshal this data into a FeatureMap
	IdentitySettings              string `env:"IDENTITY_SETTINGS,default={}"`
	PaginationKMSKeyARN           string `env:"PAGINATION_KMS_KEY_ARN,required"`
	AccessHandlerExecutionRoleARN string `env:"ACCESS_HANDLER_EXECUTION_ROLE_ARN,required"`
	RemoteConfigURL               string `env:"REMOTE_CONFIG_URL"`
	RemoteConfigHeaders           string `env:"REMOTE_CONFIG_HEADERS"`
}

type NotificationsConfig struct {
	LogLevel    string `env:"LOG_LEVEL,default=info"`
	DynamoTable string `env:"APPROVALS_TABLE_NAME,required"`
	FrontendURL string `env:"APPROVALS_FRONTEND_URL,required"`
	// This should be an instance of deploy.FeatureMap which is a specific json format for this
	// Use deploy.UnmarshalFeatureMap to unmarshal this data into a FeatureMap
	NotificationsConfig string `env:"NOTIFICATIONS_SETTINGS,default={}"`
	RemoteConfigURL     string `env:"REMOTE_CONFIG_URL"`
	RemoteConfigHeaders string `env:"REMOTE_CONFIG_HEADERS"`
}

type EventHandlerConfig struct {
	LogLevel    string `env:"LOG_LEVEL,default=info"`
	DynamoTable string `env:"APPROVALS_TABLE_NAME,required"`
}

type SyncConfig struct {
	TableName   string `env:"APPROVALS_TABLE_NAME,required"`
	IdpProvider string `env:"IDENTITY_PROVIDER,required"`
	UserPoolId  string `env:"APPROVALS_COGNITO_USER_POOL_ID,required"`
	LogLevel    string `env:"LOG_LEVEL,default=info"`
	// This should be an instance of deploy.FeatureMap which is a specific json format for this
	// Use deploy.UnmarshalFeatureMap to unmarshal this data into a FeatureMap
	IdentitySettings string `env:"IDENTITY_SETTINGS,default={}"`
}

type FrontendDeployerConfig struct {
	LogLevel                             string `env:"LOG_LEVEL,default=info"`
	Region                               string `env:"AWS_REGION,required"`
	CFReleasesBucket                     string `env:"CF_RELEASES_BUCKET,required"`
	CFReleasesFrontendAssetsObjectPrefix string `env:"CF_RELEASES_FRONTEND_ASSET_OBJECT_PREFIX,required"`
	FrontendBucket                       string `env:"FRONTEND_BUCKET,required"`
	UserPoolID                           string `env:"COGNITO_USER_POOL_ID,required"`
	CognitoClientID                      string `env:"COGNITO_CLIENT_ID,required"`
	UserPoolDomain                       string `env:"USER_POOL_DOMAIN,required"`
	FrontendDomain                       string `env:"FRONTEND_DOMAIN,required"`
	CloudFrontDistributionID             string `env:"CLOUDFRONT_DISTRIBUTION_ID,required"`
	APIURL                               string `env:"API_URL,required"`
}
