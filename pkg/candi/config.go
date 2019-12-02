package candi

import (
	"strings"

	"github.com/kelseyhightower/envconfig"
)

type Config struct {
	ApiBaseURL string `split_words:"true" default:"https://api.helochat.id"`
	BasicAuth  string `split_words:"true"`
}

func parseFromEnv(envPrefixVar string) (*Config, error) {
	var cfg Config
	err := envconfig.Process(strings.TrimRight(envPrefixVar, "_"), &cfg)
	if err != nil {
		return nil, err
	}
	return &cfg, nil
}
