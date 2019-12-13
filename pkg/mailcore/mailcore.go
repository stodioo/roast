package mailcore

import (
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ses"
)

type MailCore struct {
	sesClient      *ses.SES
	codeTTLDefault time.Duration
}

func New(envPrefixVar string) (*MailCore, error) {
	cfg := parseEnvVar(envPrefixVar)
	sess, err := session.NewSession(&aws.Config{
		Region: aws.String(cfg.Region),
	})

	if err != nil {
		return nil, err
	}

	svc := ses.New(sess)
	return &MailCore{
		sesClient:      svc,
		codeTTLDefault: cfg.CodeTTLDefault,
	}, nil
}
