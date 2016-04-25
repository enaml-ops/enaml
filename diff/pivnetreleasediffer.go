package diff

import "github.com/xchapter7x/enaml/pull"

type pivnetReleaseDiffer struct {
	ReleaseRepo pull.Release
	R1Path      string
	R2Path      string
}

func (d pivnetReleaseDiffer) Diff() (Result, error) {
	return Result{}, nil
}

func (d pivnetReleaseDiffer) DiffJob(job string) (Result, error) {
	return Result{}, nil
}
