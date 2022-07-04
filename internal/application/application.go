package application

import (
	"fmt"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/route53domains"
	"github.com/storskegg/r53transfer/internal/clients"
	"github.com/storskegg/r53transfer/internal/profiles"
	"github.com/storskegg/r53transfer/internal/transfers"
)

type Application interface {
	Run([]string) error
}

type app struct {
	Profiles         profiles.Profiles
	OmitProfiles     profiles.Profiles
	Clients          clients.Clients
	Domains2Transfer []transfers.Transfer
}

func (a *app) Run(args []string) (err error) {
	a.OmitProfiles = profiles.New()

	if a.Profiles, err = profiles.ReadProfiles(); err != nil {
		return err
	}

	sourceProfile, err := profiles.SelectSourceProfile(a.Profiles, a.OmitProfiles)
	if err != nil {
		return err
	}

	a.OmitProfiles.Add(sourceProfile)

	targetProfile, err := profiles.SelectTargetProfile(a.Profiles, a.OmitProfiles)
	if err != nil {
		return err
	}

	// The profiles are valued with empty strings, at this point, and will be backfilled as account numbers as clients
	// are instantiated.

	printSelectedProfiles(sourceProfile, targetProfile)

	fmt.Print("Initializing Clients...")
	client, err := clients.New(sourceProfile, targetProfile, a.Profiles)
	if err != nil {
		fmt.Println("ERROR")
		return err
	}
	fmt.Println("OK")
	fmt.Println()

	fmt.Println("Updated Profiles:")
	fmt.Printf("    Source: %s - %s\n", a.Profiles[sourceProfile], sourceProfile)
	fmt.Printf("    Target: %s - %s\n", a.Profiles[targetProfile], targetProfile)
	fmt.Println()

	fmt.Println("Fetching Domains...")
	listedDomains, err := client.ListSourceDomains()
	if err != nil {
		return err
	}

	for _, ds := range listedDomains.Domains {
		fmt.Println("    ", *ds.DomainName)
		a.Domains2Transfer = append(a.Domains2Transfer, transfers.Transfer{
			TransferInput: &route53domains.TransferDomainToAnotherAwsAccountInput{
				DomainName: ds.DomainName,
				AccountId:  aws.String(a.Profiles[sourceProfile]),
			},
		})
	}

	// Transfer the domains
	var transferError error
	for _, d := range a.Domains2Transfer {
		if d.TransferResponse, transferError = a.Clients.Source.TransferDomainToAnotherAwsAccount(d.TransferInput); transferError != nil {
			fmt.Println(err)
			continue
		}
		d.GenerateAcceptance()
	}

	// ListOperations - Source
	sourceOperations, operationsError := a.Clients.Source.ListOperations(&route53domains.ListOperationsInput{})
	if operationsError != nil {
		fmt.Println(err)
	}
	fmt.Println("Source Operations:\n", sourceOperations)

	targetOperations, operationsError := a.Clients.Target.ListOperations(&route53domains.ListOperationsInput{})
	if operationsError != nil {
		fmt.Println(err)
	}
	fmt.Println("Target Operations:\n", targetOperations)

	return
}

func New() Application {
	a := &app{}

	return a
}

func printSelectedProfiles(source, target string) {
	fmt.Println("Selected Profiles:")
	fmt.Printf("    Source: %s\n", source)
	fmt.Printf("    Target: %s\n", target)
	fmt.Println()
}
