package release

import (
	"io"

	"github.com/enaml-ops/enaml"
	"github.com/enaml-ops/enaml/pkg"
	"github.com/enaml-ops/enaml/pull"
)

// BoshRelease wraps a release manifest and its job manifests neatly together.
type BoshRelease struct {
	ReleaseManifest enaml.ReleaseManifest
	JobManifests    map[string]enaml.JobManifest
}

// emptyBoshRelease is a null boshRelease object. This is useful when a pivnet
// release from one release to another has missing BOSH releases.
var emptyBoshRelease = &BoshRelease{
	ReleaseManifest: enaml.ReleaseManifest{},
	JobManifests:    make(map[string]enaml.JobManifest),
}

// LoadBoshRelease creates an initialized boshRelease instance from the
// specifed local or remote .tgz file
func LoadBoshRelease(releaseRepo pull.Release, path string) (release *BoshRelease, err error) {
	var rr io.ReadCloser
	rr, err = releaseRepo.Read(path)
	if err != nil {
		return
	}
	defer func() {
		if cerr := rr.Close(); cerr != nil {
			err = cerr
		}
	}()
	release, err = readBoshRelease(rr)
	return
}

// readBoshRelease creates an initialized boshRelease instance from the
// specifed .tgz reader
func readBoshRelease(rr io.Reader) (*BoshRelease, error) {
	release := &BoshRelease{
		JobManifests: make(map[string]enaml.JobManifest),
	}
	err := release.readBoshRelease(rr)
	return release, err
}

// readBoshRelease reads a bosh release out of the given reader into a new
// boshRelease struct
func (r *BoshRelease) readBoshRelease(rr io.Reader) error {
	w := pkg.NewTgzWalker(rr)
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
func (r *BoshRelease) readBoshJob(jr io.Reader) (enaml.JobManifest, error) {
	var job enaml.JobManifest
	jw := pkg.NewTgzWalker(jr)
	jw.OnMatch("job.MF", func(file pkg.FileEntry) error {
		return decodeYaml(file.Reader, &job)
	})
	err := jw.Walk()
	return job, err
}
