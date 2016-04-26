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
	var r1, r2 *boshRelease
	r1, err = d.loadBoshRelease(d.R1Path)
	if err != nil {
		return
	}
	r2, err = d.loadBoshRelease(d.R2Path)
	if err != nil {
		return
	}
	result.Deltas = pretty.Diff(r1, r2)
	return
}

func (d boshReleaseDiffer) DiffJob(job string) (result Result, err error) {
	result = Result{}
	var r1, r2 *boshRelease
	if r1, err = d.loadBoshRelease(d.R1Path); err == nil {
		if r2, err = d.loadBoshRelease(d.R2Path); err == nil {
			var (
				job1, job2 enaml.JobManifest
				ok         bool
			)
			if job1, ok = r1.JobManifests[job]; !ok {
				err = fmt.Errorf("Couldn't find job '%s' in release 1", job)
				return
			}
			if job2, ok = r2.JobManifests[job]; !ok {
				err = fmt.Errorf("Couldn't find job '%s' in release 2", job)
				return
			}
			result.Deltas = pretty.Diff(job1, job2)
		}
	}
	return
}

func (d boshReleaseDiffer) loadBoshRelease(path string) (release *boshRelease, err error) {
	var rr io.ReadCloser
	rr, err = d.ReleaseRepo.Read(path)
	if err != nil {
		return
	}
	defer func() {
		if cerr := rr.Close(); cerr != nil {
			err = cerr
		}
	}()
	release = newBoshRelease()
	err = release.readBoshRelease(rr)
	return
}
