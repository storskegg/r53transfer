package profiles

import (
	"bufio"
	"os"
	"path"
	"regexp"
	"sort"

	"github.com/manifoldco/promptui"
)

var (
	AwsCredentialsPath = []string{".aws", "credentials"}

	DefaultProfile = "[default]"

	TestProfileHeader = regexp.MustCompile(`^\[[A-Za-z0-9_-]+\]$`) // TODO: check against actual alias rules
)

type Profiles map[string]struct{}

func (p Profiles) Sort() (profiles []string) {
	for key := range p {
		profiles = append(profiles, key)
	}

	sort.Strings(profiles)
	return
}

func (p Profiles) Add(key string) {
	p[key] = struct{}{}
}

func (p Profiles) Delete(key string) {
	delete(p, key)
}

func (p Profiles) Exists(key string) (ok bool) {
	_, ok = p[key]
	return
}

func New() Profiles {
	return make(Profiles)
}

func ReadProfiles() (profiles Profiles, err error) {
	var credPath string
	var credFile *os.File

	if credPath, err = CredentialsPath(); err != nil {
		return nil, err
	}

	if credFile, err = os.Open(credPath); err != nil {
		return nil, err
	}
	defer credFile.Close()

	profiles = New()

	profileScanner := bufio.NewScanner(credFile)
	for profileScanner.Scan() {
		t := profileScanner.Text()
		if TestProfileHeader.MatchString(t) && t != DefaultProfile {
			profiles.Add(t[1 : len(t)-1])
		}
	}

	return
}

func SelectSourceProfile(haystack Profiles, omitProfiles Profiles) (profile string, err error) {
	displayProfiles := New()

	for h := range haystack {
		if _, omit := omitProfiles[h]; !omit {
			displayProfiles.Add(h)
		}
	}

	prompt := promptui.Select{
		Label: "Source Profile",
		Items: displayProfiles.Sort(),
	}

	_, profile, err = prompt.Run()

	return
}

func SelectTargetProfile(haystack Profiles, omitProfiles Profiles) (profile string, err error) {
	displayProfiles := New()

	for h := range haystack {
		if _, omit := omitProfiles[h]; !omit {
			displayProfiles.Add(h)
		}
	}

	prompt := promptui.Select{
		Label: "Target Profile",
		Items: displayProfiles.Sort(),
	}

	_, profile, err = prompt.Run()

	return
}

func CredentialsPath() (credPath string, err error) {
	var userDir string
	if userDir, err = os.UserHomeDir(); err != nil {
		return "", err
	}

	return path.Join(append([]string{userDir}, AwsCredentialsPath...)...), nil
}
