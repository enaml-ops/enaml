package diff

import (
	"io"
	"io/ioutil"
	"os"

	"gopkg.in/yaml.v2"

	"github.com/kr/pretty"
	"github.com/xchapter7x/enaml"
	"github.com/xchapter7x/enaml/pkg"
	"github.com/xchapter7x/enaml/pull"
)

type boshReleaseDiffer struct {
	ReleaseRepo pull.Release
	R1Path      string
	R2Path      string
}

func (d boshReleaseDiffer) Diff() (Result, error) {
	result := Result{}
	// TODO: close file?
	// TODO: release repo return an io.Reader
	// TODO: move decodeYaml to yaml.go or something
	// TODO: move release to another file, rename boshrelease?
	// TODO: move reading a bosh release to ^ ?
	r1, err := d.boshReleaseReader(d.R1Path)
	if err != nil {
		return result, err
	}
	r2, err := d.boshReleaseReader(d.R2Path)
	if err != nil {
		return result, err
	}
	return d.diffBoshReleases(r1, r2)
}

func (d boshReleaseDiffer) diffBoshReleases(r1, r2 io.Reader) (Result, error) {
	result := Result{}
	release1, err := readRelease(r1)
	if err != nil {
		return result, err
	}
	release2, err := readRelease(r2)
	if err != nil {
		return result, err
	}
	result.Deltas = pretty.Diff(release1, release2)
	return result, nil
}

// release wraps a release manifest and its job manifests neatly together.
type release struct {
	ReleaseManifest enaml.ReleaseManifest
	JobManifests    map[string]enaml.JobManifest
}

func newRelease() *release {
	return &release{
		JobManifests: make(map[string]enaml.JobManifest),
	}
}

func decodeYaml(r io.Reader, v interface{}) error {
	bytes, err := ioutil.ReadAll(r)
	if err == nil {
		yaml.Unmarshal(bytes, &v)
	}
	return err
}

func readRelease(rr io.Reader) (*release, error) {
	diffRelease := newRelease()
	w := pkg.NewWalker(rr)
	w.OnMatch("release.MF", func(file pkg.FileEntry) error {
		return decodeYaml(file.Reader, diffRelease.ReleaseManifest)
	})
	w.OnMatch("/jobs/", func(file pkg.FileEntry) error {
		job, jerr := readJob(file.Reader)
		if jerr == nil {
			diffRelease.JobManifests[job.Name] = job
		}
		return jerr
	})
	err := w.Walk()
	return diffRelease, err
}

// readJob reads the job manifest.
func readJob(jr io.Reader) (enaml.JobManifest, error) {
	var job enaml.JobManifest
	jw := pkg.NewWalker(jr)
	jw.OnMatch("job.MF", func(file pkg.FileEntry) error {
		return decodeYaml(file.Reader, job)
	})
	err := jw.Walk()
	return job, err
}

func (d boshReleaseDiffer) boshReleaseReader(path string) (io.Reader, error) {
	local, err := d.ReleaseRepo.Pull(path)
	if err != nil {
		return nil, err
	}
	r, err := os.Open(local)
	if err != nil {
		return nil, err
	}
	return r, nil
}
