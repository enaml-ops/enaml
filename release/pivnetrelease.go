package release

import (
	"github.com/xchapter7x/enaml/pkg"
	"github.com/xchapter7x/enaml/pull"
)

// PivnetRelease wraps a .pivotal release and all of its contained BOSH
// releases and jobs
type PivnetRelease struct {
	BoshRelease map[string]*BoshRelease
}

// LoadPivnetRelease creates an initialized pivnetRelease instance from the
// specified .pivotal file.
func LoadPivnetRelease(releaseRepo pull.Release, path string) (release *PivnetRelease, err error) {
	release = &PivnetRelease{}
	var localPath string
	localPath, err = releaseRepo.Pull(path)
	if err != nil {
		return
	}
	release = &PivnetRelease{
		BoshRelease: make(map[string]*BoshRelease),
	}
	err = release.readPivnetRelease(localPath)
	return
}

// BoshReleaseOrEmpty returns the named BOSH release from this pivnet release
// if it exists, otherwise emptyBoshRelease is returned.
func (r *PivnetRelease) BoshReleaseOrEmpty(name string) *BoshRelease {
	br := r.BoshRelease[name]
	if br == nil {
		br = emptyBoshRelease
	}
	return br
}

// readPivnetRelease reads a pivnet release out of the given reader into a new
// pivnetRelease struct
func (r *PivnetRelease) readPivnetRelease(path string) error {
	walker := pkg.NewZipWalker(path)
	walker.OnMatch("releases/", func(file pkg.FileEntry) error {
		br, berr := readBoshRelease(file.Reader)
		if berr != nil {
			return berr
		}
		r.BoshRelease[br.ReleaseManifest.Name] = br
		return nil
	})
	return walker.Walk()
}
