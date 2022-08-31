package config

type Config struct {
	Host                   string `env:"ACCESS_HANDLER_HOST,default=0.0.0.0:9092"`
	LogLevel               string `env:"LOG_LEVEL,default=info"`
	Runtime                string `env:"GRANTED_RUNTIME,required"`
	EventBusArn            string `env:"EVENT_BUS_ARN"`
	EventBusSource         string `env:"EVENT_BUS_SOURCE"`
	AssumeExecutionRoleARN string `env:"ASSUME_EXECUTION_ROLE_ARN"`
}

type Runtime struct {
	Runtime string `env:"GRANTED_RUNTIME,required"`
}

type GranterConfig struct {
	LogLevel               string `env:"LOG_LEVEL,default=info"`
	EventBusArn            string `env:"EVENT_BUS_ARN"`
	EventBusSource         string `env:"EVENT_BUS_SOURCE"`
	AssumeExecutionRoleARN string `env:"ASSUME_EXECUTION_ROLE_ARN"`
}
