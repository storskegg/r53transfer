package clients

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/route53domains"
	"github.com/storskegg/r53transfer/internal/profiles"
)

type Clients struct {
	Source *route53domains.Route53Domains
	Target *route53domains.Route53Domains
}

func (c *Clients) InitSource(profile string) (err error) {
	var credPath string

	if credPath, err = profiles.CredentialsPath(); err != nil {
		return
	}

	cfg := aws.NewConfig().WithCredentials(credentials.NewSharedCredentials(credPath, profile))
	s := session.Must(session.NewSession(cfg))
	c.Source = route53domains.New(s)

	return
}

func (c *Clients) InitTarget(profile string) (err error) {
	var credPath string

	if credPath, err = profiles.CredentialsPath(); err != nil {
		return
	}

	cfg := aws.NewConfig().WithCredentials(credentials.NewSharedCredentials(credPath, profile))
	s := session.Must(session.NewSession(cfg))
	c.Target = route53domains.New(s)

	return
}

func New(source string, target string) (c *Clients, err error) {
	if err = c.InitSource(source); err != nil {
		return nil, err
	}

	if err = c.InitTarget(target); err != nil {
		return nil, err
	}

	return
}
