package diff

import (
	"github.com/xchapter7x/enaml/pkg"
	"github.com/xchapter7x/enaml/pull"
)

// pivnetRelease wraps a .pivotal release and all of its contained BOSH
// releases and jobs
type pivnetRelease struct {
	boshRelease map[string]*boshRelease
}

// loadPivnetRelease creates an initialized pivnetRelease instance from the
// specified .pivotal file.
func loadPivnetRelease(releaseRepo pull.Release, path string) (release *pivnetRelease, err error) {
	release = &pivnetRelease{}
	var localPath string
	localPath, err = releaseRepo.Pull(path)
	if err != nil {
		return
	}
	release = &pivnetRelease{
		boshRelease: make(map[string]*boshRelease),
	}
	err = release.readPivnetRelease(localPath)
	return
}

// readPivnetRelease reads a pivnet release out of the given reader into a new
// pivnetRelease struct
func (r *pivnetRelease) readPivnetRelease(path string) error {
	walker := pkg.NewZipWalker(path)
	walker.OnMatch("releases/", func(file pkg.FileEntry) error {
		br, berr := readBoshRelease(file.Reader)
		if berr != nil {
			return berr
		}
		r.boshRelease[br.ReleaseManifest.Name] = br
		return nil
	})
	return walker.Walk()
}
