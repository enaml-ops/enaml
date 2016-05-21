package run

import (
	"fmt"
	"io"
	"path/filepath"

	"github.com/enaml-ops/enaml"
	"github.com/enaml-ops/enaml/pull"
	"github.com/enaml-ops/enaml/release"
)

// ShowCmd runs the show CLI command
type ShowCmd struct {
	releaseRepo pull.Release
	release     string
}

// NewShowCmd creates a new ShowCmd instance.
func NewShowCmd(releaseRepo pull.Release, release string) *ShowCmd {
	return &ShowCmd{
		releaseRepo: releaseRepo,
		release:     release,
	}
}

// All writes out all the release data to writer.
func (s *ShowCmd) All(w io.Writer) error {
	if filepath.Ext(s.release) == ".pivotal" {
		pivnetRelease, err := release.LoadPivnetRelease(s.releaseRepo, s.release)
		if err != nil {
			return err
		}
		for _, br := range pivnetRelease.BoshRelease {
			s.printBoshRelease(w, br)
		}
		return nil
	}
	boshRelease, err := release.LoadBoshRelease(s.releaseRepo, s.release)
	if err != nil {
		return err
	}
	s.printBoshRelease(w, boshRelease)
	return nil
}

func (s *ShowCmd) printBoshRelease(w io.Writer, br *release.BoshRelease) {
	for _, j := range br.JobManifests {
		s.printBoshJob(w, j, br.ReleaseManifest.Name)
	}
	fmt.Fprintln(w)
}

func (s *ShowCmd) printBoshJob(w io.Writer, j enaml.JobManifest, boshReleaseName string) {
	fmt.Fprintln(w, "------------------------------------------------------")
	fmt.Fprintln(w, fmt.Sprintf("Release: %s", boshReleaseName))
	fmt.Fprintln(w, fmt.Sprintf("Job:     %s", j.Name))
	fmt.Fprintln(w, "------------------------------------------------------")
	for pname, p := range j.Properties {
		fmt.Fprintln(w, pname)
		s.printBoshJobProperty(w, p)
		fmt.Fprintln(w)
	}
	fmt.Fprintln(w)
}

func (s *ShowCmd) printBoshJobProperty(w io.Writer, p enaml.JobManifestProperty) {
	if len(p.Description) > 0 {
		fmt.Fprintln(w, fmt.Sprintf("  Description: %s", p.Description))
	}
	if p.Default != nil {
		fmt.Fprintln(w, fmt.Sprintf("  Default: %v", p.Default))
	}
}
