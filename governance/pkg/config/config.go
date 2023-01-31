package config

type Config struct {
	GovernanceURL       string `env:"COMMONFATE_GOVERNANCE_URL,default=0.0.0.0:8889"`
	AccessHandlerURL    string `env:"COMMONFATE_ACCESS_HANDLER_URL,default=http://0.0.0.0:9092"`
	LogLevel            string `env:"LOG_LEVEL,default=info"`
	DynamoTable         string `env:"COMMONFATE_TABLE_NAME,required"`
	Region              string `env:"AWS_REGION,required"`
	PaginationKMSKeyARN string `env:"COMMONFATE_PAGINATION_KMS_KEY_ARN,required"`
	MockAccessHandler   bool   `env:"COMMONFATE_MOCK_ACCESS_HANDLER,default=false"`
	RunAccessHandler    bool   `env:"COMMONFATE_RUN_ACCESS_HANDLER,default=true"`
}
