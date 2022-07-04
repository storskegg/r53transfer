package application

import (
	"fmt"

	"github.com/storskegg/r53transfer/internal/clients"
	"github.com/storskegg/r53transfer/internal/profiles"
)

type Application interface {
	Run([]string) error
}

type app struct {
	Profiles     profiles.Profiles
	OmitProfiles profiles.Profiles
	Clients      clients.Clients
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

	printSelectedProfiles(sourceProfile, targetProfile)

	fmt.Print("Initializing Clients...")
	client, err := clients.New(sourceProfile, targetProfile, a.Profiles)
	if err != nil {
		fmt.Println("ERROR")
		return err
	}
	fmt.Println("OK")
	fmt.Println()

	//fmt.Println("Updated Profiles:")
	//fmt.Printf("    Source: %s - %s\n", a.Profiles[sourceProfile], sourceProfile)
	//fmt.Printf("    Target: %s - %s\n", a.Profiles[targetProfile], targetProfile)
	//fmt.Println()

	fmt.Println("Fetching Domains...")
	out, err := client.ListSourceDomains()
	if err != nil {
		return err
	}
	fmt.Println(out)

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
