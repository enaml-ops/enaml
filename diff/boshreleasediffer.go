package diff

import (
	"github.com/kr/pretty"
	"github.com/xchapter7x/enaml"
	"github.com/xchapter7x/enaml/release"
)

type boshReleaseDiffer struct {
	release1 *release.BoshRelease
	release2 *release.BoshRelease
}

func (d boshReleaseDiffer) Diff() (result Result, err error) {
	result = Result{}
	var jresult Result
	for _, jname := range d.allJobNames() {
		jresult, err = d.DiffJob(jname)
		result.Deltas = append(result.Deltas, jresult.Deltas...)
	}
	return
}

func (d boshReleaseDiffer) DiffJob(job string) (result Result, err error) {
	result = Result{}
	var job1, job2 enaml.JobManifest
	job1 = d.release1.JobManifests[job]
	job2 = d.release2.JobManifests[job]
	result.Deltas = pretty.Diff(job1, job2)
	return
}

// allJobNames returns a union of unique job names across both BOSH releases
func (d boshReleaseDiffer) allJobNames() []string {
	jobNamesMap := make(map[string]string)
	var addJobNames = func(br *release.BoshRelease) {
		if br != nil {
			for jbname := range br.JobManifests {
				jobNamesMap[jbname] = jbname
			}
		}
	}
	addJobNames(d.release1)
	addJobNames(d.release2)
	var jobNames []string
	for jname := range jobNamesMap {
		jobNames = append(jobNames, jname)
	}
	return jobNames
}
