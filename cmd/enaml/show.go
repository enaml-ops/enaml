package main

import (
	"fmt"
	"io"
	"path/filepath"

	"github.com/xchapter7x/enaml"
	"github.com/xchapter7x/enaml/pull"
	"github.com/xchapter7x/enaml/release"
)

type show struct {
	releaseRepo pull.Release
	release     string
}

func (s *show) All(w io.Writer) error {
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

func (s *show) printBoshRelease(w io.Writer, br *release.BoshRelease) {
	for _, j := range br.JobManifests {
		s.printBoshJob(w, j, br.ReleaseManifest.Name)
	}
	fmt.Fprintln(w)
}

func (s *show) printBoshJob(w io.Writer, j enaml.JobManifest, boshReleaseName string) {
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

func (s *show) printBoshJobProperty(w io.Writer, p enaml.JobManifestProperty) {
	if len(p.Description) > 0 {
		fmt.Fprintln(w, fmt.Sprintf("  Description: %s", p.Description))
	}
	if p.Default != nil {
		fmt.Fprintln(w, fmt.Sprintf("  Default: %v", p.Default))
	}
}
