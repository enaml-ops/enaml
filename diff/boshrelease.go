package diff

import (
	"io"

	"github.com/xchapter7x/enaml"
	"github.com/xchapter7x/enaml/pkg"
)

// boshRelease wraps a release manifest and its job manifests neatly together.
type boshRelease struct {
	ReleaseManifest enaml.ReleaseManifest
	JobManifests    map[string]enaml.JobManifest
}

// newBoshRelease creates a new empty boshRelease struct
func newBoshRelease() *boshRelease {
	return &boshRelease{
		JobManifests: make(map[string]enaml.JobManifest),
	}
}

// readBoshRelease reads a bosh release out of the given reader into a new
// boshRelease struct
func (r *boshRelease) readBoshRelease(rr io.Reader) error {
	w := pkg.NewWalker(rr)
	w.OnMatch("release.MF", func(file pkg.FileEntry) error {
		return decodeYaml(file.Reader, &r.ReleaseManifest)
	})
	w.OnMatch("/jobs/", func(file pkg.FileEntry) error {
		job, jerr := r.readBoshJob(file.Reader)
		if jerr == nil {
			r.JobManifests[job.Name] = job
		}
		return jerr
	})
	err := w.Walk()
	return err
}

// readBoshJob reads a BOSH job manifest out of the given reader into a new
// JobManifest struct
func (r *boshRelease) readBoshJob(jr io.Reader) (enaml.JobManifest, error) {
	var job enaml.JobManifest
	jw := pkg.NewWalker(jr)
	jw.OnMatch("job.MF", func(file pkg.FileEntry) error {
		return decodeYaml(file.Reader, &job)
	})
	err := jw.Walk()
	return job, err
}
