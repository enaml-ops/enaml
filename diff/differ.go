package diff

import (
	"fmt"
	"path/filepath"

	"github.com/xchapter7x/enaml/pull"
)

// Result is returned from a diff operation
type Result struct {
	Deltas []string
}

// Differ implements diffing BOSH or Pivnet releases and their contained entities.
type Differ interface {
	Diff() (Result, error)
	DiffJob(job string) (Result, error)
}

// New creates a Differ instance for comparing two releases
func New(releaseRepo pull.Release, r1Path, r2Path string) (differ Differ, err error) {
	if filepath.Ext(r1Path) != filepath.Ext(r2Path) {
		err = fmt.Errorf("The specified releases didn't have matching file extensions, " +
			"assuming different release types.")
		return
	}
	if filepath.Ext(r1Path) == ".pivotal" {
		differ = pivnetReleaseDiffer{
			ReleaseRepo: releaseRepo,
			R1Path:      r1Path,
			R2Path:      r2Path,
		}
	} else {
		differ = boshReleaseDiffer{
			ReleaseRepo: releaseRepo,
			R1Path:      r1Path,
			R2Path:      r2Path,
		}
	}
	return
}
