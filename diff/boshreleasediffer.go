package diff

import (
	"fmt"
	"io"

	"github.com/kr/pretty"
	"github.com/xchapter7x/enaml"
	"github.com/xchapter7x/enaml/pull"
)

type boshReleaseDiffer struct {
	ReleaseRepo pull.Release
	R1Path      string
	R2Path      string
}

func (d boshReleaseDiffer) Diff() (result Result, err error) {
	result = Result{}
	err = d.withBoshRelease(d.R1Path, func(r1 *boshRelease) error {
		err = d.withBoshRelease(d.R2Path, func(r2 *boshRelease) error {
			result.Deltas = pretty.Diff(r1, r2)
			return nil
		})
		return nil
	})
	return
}

func (d boshReleaseDiffer) DiffJob(job string) (result Result, err error) {
	result = Result{}
	err = d.withBoshRelease(d.R1Path, func(r1 *boshRelease) error {
		err = d.withBoshRelease(d.R2Path, func(r2 *boshRelease) error {
			var (
				job1, job2 enaml.JobManifest
				ok         bool
			)
			if job1, ok = r1.JobManifests[job]; !ok {
				return fmt.Errorf("Couldn't find job '%s' in release 1", job)
			}
			if job2, ok = r2.JobManifests[job]; !ok {
				return fmt.Errorf("Couldn't find job '%s' in release 2", job)
			}
			result.Deltas = pretty.Diff(job1, job2)
			return nil
		})
		return nil
	})
	return
}

type withBoshReleaseFunc func(r *boshRelease) error

func (d boshReleaseDiffer) withBoshRelease(path string, fn withBoshReleaseFunc) (err error) {
	var rr io.ReadCloser
	rr, err = d.ReleaseRepo.Read(path)
	if err != nil {
		return err
	}
	defer func() {
		if cerr := rr.Close(); cerr != nil {
			err = cerr
		}
	}()
	release := newBoshRelease()
	if err = release.readBoshRelease(rr); err != nil {
		return err
	}
	return fn(release)
}
