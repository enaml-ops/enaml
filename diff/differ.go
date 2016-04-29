package diff

import (
	"fmt"
	"path/filepath"

	"github.com/xchapter7x/enaml/pull"
	"github.com/xchapter7x/enaml/release"
)

// Differ implements diffing BOSH or Pivnet releases and their contained entities.
type Differ interface {
	Diff() (*Result, error)
	DiffJob(job string) (*Result, error)
}

// New creates a Differ instance for comparing two releases
func New(releaseRepo pull.Release, r1Path, r2Path string) (differ Differ, err error) {
	if filepath.Ext(r1Path) != filepath.Ext(r2Path) {
		err = fmt.Errorf("The specified releases didn't have matching file extensions, " +
			"assuming different release types.")
		return
	}
	if filepath.Ext(r1Path) == ".pivotal" {
		var r1, r2 *release.PivnetRelease
		if r1, err = release.LoadPivnetRelease(releaseRepo, r1Path); err == nil {
			if r2, err = release.LoadPivnetRelease(releaseRepo, r2Path); err == nil {
				differ = pivnetReleaseDiffer{
					release1: r1,
					release2: r2,
				}
			}
		}
	} else {
		var r1, r2 *release.BoshRelease
		if r1, err = release.LoadBoshRelease(releaseRepo, r1Path); err == nil {
			if r2, err = release.LoadBoshRelease(releaseRepo, r2Path); err == nil {
				differ = boshReleaseDiffer{
					release1: r1,
					release2: r2,
				}
			}
		}
	}
	return
}
