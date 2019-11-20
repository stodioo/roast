package blogapp

import "github.com/kelseyhightower/envconfig"

type Config struct {
	BlogDBUrl string `required:"true" split_words:"true"`
}

func ParseConfigFromEnv(envVarPrefix string) *Config {
	var config Config
	err := envconfig.Process(envVarPrefix, &config)

	if err != nil {
		panic(err)
	}

	return &config
}
