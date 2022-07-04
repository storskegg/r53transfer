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
	Clients clients.Clients
}

func (a *app) Run(args []string) (err error) {
	omitProfiles := profiles.New()
	p, err := profiles.ReadProfiles()
	if err != nil {
		return err
	}

	sourceProfile, err := profiles.SelectSourceProfile(p, omitProfiles)
	if err != nil {
		return err
	}

	omitProfiles[sourceProfile] = struct{}{}

	targetProfile, err := profiles.SelectTargetProfile(p, omitProfiles)
	if err != nil {
		return err
	}

	fmt.Printf("Source Profile: %s\n", sourceProfile)
	fmt.Printf("Target Profile: %s\n", targetProfile)

	return
}

func New() Application {
	a := &app{}

	return a
}
