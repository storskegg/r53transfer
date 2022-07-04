package clients

import (
	"fmt"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/route53domains"
	"github.com/aws/aws-sdk-go/service/sts"
	"github.com/storskegg/r53transfer/internal/profiles"
)

type Clients struct {
	Source *route53domains.Route53Domains
	Target *route53domains.Route53Domains
}

func (c *Clients) InitSource(profile string, p profiles.Profiles) (err error) {
	var credPath string

	if credPath, err = profiles.CredentialsPath(); err != nil {
		return
	}

	creds := credentials.NewSharedCredentials(credPath, profile)

	cfg := aws.NewConfig().WithCredentials(creds).WithRegion("us-east-1")
	sess := session.Must(session.NewSession(cfg))
	c.Source = route53domains.New(sess)

	accountNumber, err := AccountNumberWithConfig(cfg)
	if err != nil {
		return err
	}
	if err = p.AddAccountNumber(profile, accountNumber); err != nil {
		return
	}

	return
}

func (c *Clients) InitTarget(profile string, p profiles.Profiles) (err error) {
	var credPath string

	if credPath, err = profiles.CredentialsPath(); err != nil {
		return
	}

	creds := credentials.NewSharedCredentials(credPath, profile)

	cfg := aws.NewConfig().WithCredentials(creds).WithRegion("us-east-1")
	sess := session.Must(session.NewSession(cfg))
	c.Target = route53domains.New(sess)

	accountNumber, err := AccountNumberWithConfig(cfg)
	if err != nil {
		return err
	}
	if err = p.AddAccountNumber(profile, accountNumber); err != nil {
		return
	}

	return
}

func (c *Clients) ListSourceDomains() (*route53domains.ListDomainsOutput, error) {
	out, err := c.Source.ListDomains(&route53domains.ListDomainsInput{})
	if err != nil {
		return nil, err
	}

	return out, nil
}

func New(source string, target string, p profiles.Profiles) (c *Clients, err error) {
	c = &Clients{}

	if err = c.InitSource(source, p); err != nil {
		return nil, err
	}

	if err = c.InitTarget(target, p); err != nil {
		return nil, err
	}

	return
}

func AccountNumberWithConfig(cfg *aws.Config) (accountNumber string, err error) {
	sess, err := session.NewSession(cfg)
	if err != nil {
		return "", err
	}

	svc := sts.New(sess)
	input := &sts.GetCallerIdentityInput{}

	result, err := svc.GetCallerIdentity(input)
	if err != nil {
		if aerr, ok := err.(awserr.Error); ok {
			switch aerr.Code() {
			default:
				fmt.Println(aerr.Error())
			}
		} else {
			// Print the error, cast err to awserr.Error to get the Code and
			// Message from an error.
			fmt.Println(err.Error())
		}
		return "", nil
	}

	return *result.Account, nil
}
