package config

type Config struct {
	Host           string `env:"ACCESS_HANDLER_HOST,default=0.0.0.0:9092"`
	LogLevel       string `env:"LOG_LEVEL,default=info"`
	Runtime        string `env:"COMMONFATE_ACCESS_HANDLER_RUNTIME,required"`
	EventBusArn    string `env:"COMMONFATE_EVENT_BUS_ARN"`
	EventBusSource string `env:"COMMONFATE_EVENT_BUS_SOURCE"`
}

type Runtime struct {
	Runtime string `env:"COMMONFATE_ACCESS_HANDLER_RUNTIME,required"`
}

type GranterConfig struct {
	LogLevel       string `env:"LOG_LEVEL,default=info"`
	EventBusArn    string `env:"COMMONFATE_EVENT_BUS_ARN"`
	EventBusSource string `env:"COMMONFATE_EVENT_BUS_SOURCE"`
}
