package mailcore

import (
	"time"

	"github.com/kelseyhightower/envconfig"
)

type Config struct {
	Region         string        `split_words:"true"`
	CodeTTLDefault time.Duration `split_words:"true"`
}

func parseEnvVar(envVarPrefix string) *Config {
	var cfg Config
	err := envconfig.Process(envVarPrefix, &cfg)

	if err != nil {
		panic(err)
	}

	return &cfg
}
