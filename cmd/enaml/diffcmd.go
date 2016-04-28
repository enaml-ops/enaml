package main

import (
	"fmt"
	"io"

	"github.com/xchapter7x/enaml"
	"github.com/xchapter7x/enaml/diff"
	"github.com/xchapter7x/enaml/pull"
)

type diffCmd struct {
	releaseRepo pull.Release
	release1    string
	release2    string
}

func (s *diffCmd) All(w io.Writer) error {
	differ, err := diff.New(s.releaseRepo, s.release1, s.release2)
	if err != nil {
		return err
	}
	d, err := differ.Diff()
	if err != nil {
		return err
	}
	s.printDiffResult(w, d)
	return nil
}

func (s *diffCmd) Job(job string, w io.Writer) error {
	differ, err := diff.New(s.releaseRepo, s.release1, s.release2)
	if err != nil {
		return err
	}
	d, err := differ.DiffJob(job)
	if err != nil {
		return err
	}
	s.printDiffResult(w, d)
	return nil
}

func (s *diffCmd) printDiffResult(w io.Writer, d *diff.Result) {
	for _, j := range d.DeltaJob {
		s.printDeltaJob(w, &j)
	}
	fmt.Fprintln(w)
}

func (s *diffCmd) printDeltaJob(w io.Writer, j *diff.DeltaJob) {
	fmt.Fprintln(w, "------------------------------------------------------")
	fmt.Fprintln(w, fmt.Sprintf("Release: %s", j.ReleaseName))
	fmt.Fprintln(w, fmt.Sprintf("Job:     %s", j.JobName))
	fmt.Fprintln(w, "------------------------------------------------------")
	for pname, prop := range j.AddedProperties {
		fmt.Fprintln(w, fmt.Sprintf("+ %s", pname))
		s.printBoshJobProperty(w, "+", prop)
		fmt.Fprintln(w)
	}
	for pname, prop := range j.RemovedProperties {
		fmt.Fprintln(w, fmt.Sprintf("- %s", pname))
		s.printBoshJobProperty(w, "-", prop)
		fmt.Fprintln(w)
	}
	fmt.Fprintln(w)
}

func (s *diffCmd) printBoshJobProperty(w io.Writer, addedRemoved string, p enaml.JobManifestProperty) {
	if len(p.Description) > 0 {
		fmt.Fprintln(w, fmt.Sprintf("%s   Description: %s", addedRemoved, p.Description))
	}
	if p.Default != nil {
		fmt.Fprintln(w, fmt.Sprintf("%s   Default: %v", addedRemoved, p.Default))
	}
}
