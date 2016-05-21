package diff

import (
	"github.com/enaml-ops/enaml"
	"github.com/enaml-ops/enaml/release"
)

type boshReleaseDiffer struct {
	release1 *release.BoshRelease
	release2 *release.BoshRelease
}

func (d boshReleaseDiffer) Diff() (*Result, error) {
	result := &Result{}
	for _, jname := range d.allJobNames() {
		jresult, err := d.DiffJob(jname)
		if err != nil {
			return nil, err
		}
		result.Concat(jresult)
	}
	return result, nil
}

func (d boshReleaseDiffer) DiffJob(job string) (result *Result, err error) {
	var job1, job2 enaml.JobManifest
	job1 = d.release1.JobManifests[job]
	job2 = d.release2.JobManifests[job]

	dj := newDeltaJob(d.release1.ReleaseManifest.Name, job)
	result = &Result{}
	result.AddDeltaJob(dj)

	// added properties or new jobs
	for pname, prop := range job2.Properties {
		if _, ok := job1.Properties[pname]; !ok {
			dj.AddedProperty(pname, &prop)
		}
	}

	// removed properties or removed jobs
	for pname, prop := range job1.Properties {
		if _, ok := job2.Properties[pname]; !ok {
			dj.RemovedProperty(pname, &prop)
		}
	}

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
