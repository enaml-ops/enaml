package diff

import "github.com/xchapter7x/enaml"

// Result is returned from a diff operation
type Result struct {
	DeltaJob []DeltaJob
}

// DeltaJob contains all the add/remove property differences of a job between
// releases.
type DeltaJob struct {
	ReleaseName       string
	JobName           string
	AddedProperties   map[string]enaml.JobManifestProperty
	RemovedProperties map[string]enaml.JobManifestProperty
}

func newDeltaJob(release, job string) *DeltaJob {
	return &DeltaJob{
		ReleaseName:       release,
		JobName:           job,
		AddedProperties:   make(map[string]enaml.JobManifestProperty),
		RemovedProperties: make(map[string]enaml.JobManifestProperty),
	}
}

// AddedProperty adds a new "added" job property to the list of differences
func (dj *DeltaJob) AddedProperty(name string, p *enaml.JobManifestProperty) {
	dj.AddedProperties[name] = *p
}

// RemovedProperty adds a new "removed" job property to the list of differences
func (dj *DeltaJob) RemovedProperty(name string, p *enaml.JobManifestProperty) {
	dj.RemovedProperties[name] = *p
}

// AddDeltaJob adds a new delta for a specific job
func (r *Result) AddDeltaJob(dj *DeltaJob) {
	r.DeltaJob = append(r.DeltaJob, *dj)
}

// Concat adds the other result to this result
func (r *Result) Concat(other *Result) {
	for _, dj := range other.DeltaJob {
		r.DeltaJob = append(r.DeltaJob, dj)
	}
}
