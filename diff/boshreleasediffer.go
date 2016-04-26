package diff

import (
	"fmt"

	"github.com/kr/pretty"
	"github.com/xchapter7x/enaml"
)

type boshReleaseDiffer struct {
	release1 *boshRelease
	release2 *boshRelease
}

func (d boshReleaseDiffer) Diff() (result Result, err error) {
	result = Result{}
	result.Deltas = pretty.Diff(d.release1, d.release2)
	return
}

func (d boshReleaseDiffer) DiffJob(job string) (result Result, err error) {
	result = Result{}
	var (
		job1, job2 enaml.JobManifest
		ok         bool
	)
	if job1, ok = d.release1.JobManifests[job]; !ok {
		err = fmt.Errorf("Couldn't find job '%s' in release 1", job)
		return
	}
	if job2, ok = d.release2.JobManifests[job]; !ok {
		err = fmt.Errorf("Couldn't find job '%s' in release 2", job)
		return
	}
	result.Deltas = pretty.Diff(job1, job2)
	return
}
