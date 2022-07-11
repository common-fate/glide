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
	IdpProvider       string `env:"IDENTITY_PROVIDER,default=COGNITO"`
	IdentitySettings  string `env:"IDENTITY_SETTINGS,default={}"`
}

type SlackNotifierConfig struct {
	LogLevel      string `env:"LOG_LEVEL,default=info"`
	DynamoTable   string `env:"APPROVALS_TABLE_NAME,required"`
	FrontendURL   string `env:"APPROVALS_FRONTEND_URL,required"`
	SlackSettings string `env:"SLACK_SETTINGS,required"`
}

type EventHandlerConfig struct {
	LogLevel    string `env:"LOG_LEVEL,default=info"`
	DynamoTable string `env:"APPROVALS_TABLE_NAME,required"`
}

type SyncConfig struct {
	TableName        string `env:"APPROVALS_TABLE_NAME,required"`
	IdpProvider      string `env:"IDENTITY_PROVIDER,default=COGNITO"`
	UserPoolId       string `env:"APPROVALS_COGNITO_USER_POOL_ID"`
	LogLevel         string `env:"LOG_LEVEL,default=info"`
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
