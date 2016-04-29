package diff

import "github.com/xchapter7x/enaml/release"

type pivnetReleaseDiffer struct {
	release1 *release.PivnetRelease
	release2 *release.PivnetRelease
}

func (d pivnetReleaseDiffer) Diff() (Result, error) {
	return d.doDiff(func(brd boshReleaseDiffer) (Result, error) {
		return brd.Diff()
	})
}

func (d pivnetReleaseDiffer) DiffJob(job string) (Result, error) {
	return d.doDiff(func(brd boshReleaseDiffer) (Result, error) {
		return brd.DiffJob(job)
	})
}

type diffFunc func(brd boshReleaseDiffer) (Result, error)

func (d pivnetReleaseDiffer) doDiff(fn diffFunc) (Result, error) {
	result := Result{}
	for _, rname := range d.allBoshReleaseNames() {
		br1 := d.release1.BoshReleaseOrEmpty(rname)
		br2 := d.release2.BoshReleaseOrEmpty(rname)
		boshDiffer := boshReleaseDiffer{
			release1: br1,
			release2: br2,
		}
		boshDiffResult, err := fn(boshDiffer)
		if err != nil {
			return Result{}, err
		}
		result.Deltas = append(result.Deltas, boshDiffResult.Deltas...)
	}
	return result, nil
}

// allBoshReleaseNames returns a union of unique BOSH release names across all
// contained BOSH releases.
func (d pivnetReleaseDiffer) allBoshReleaseNames() []string {
	boshReleaseNamesMap := make(map[string]string)
	var addReleaseNames = func(br map[string]*release.BoshRelease) {
		for brname := range br {
			boshReleaseNamesMap[brname] = brname
		}
	}
	addReleaseNames(d.release1.BoshRelease)
	addReleaseNames(d.release2.BoshRelease)
	var releaseNames []string
	for brname := range boshReleaseNamesMap {
		releaseNames = append(releaseNames, brname)
	}
	return releaseNames
}
