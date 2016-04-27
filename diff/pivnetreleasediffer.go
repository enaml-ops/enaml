package diff

type pivnetReleaseDiffer struct {
	release1 *pivnetRelease
	release2 *pivnetRelease
}

func (d pivnetReleaseDiffer) Diff() (Result, error) {
	result := Result{}
	for rname, br1 := range d.release1.boshRelease {
		// TODO: need to handle the case where the BOSH release doesn't exist in the other .pivotal file
		br2 := d.release2.boshRelease[rname]
		boshDiffer := boshReleaseDiffer{
			release1: br1,
			release2: br2,
		}
		boshDiffResult, err := boshDiffer.Diff()
		if err != nil {
			return Result{}, err
		}
		result.Deltas = append(result.Deltas, boshDiffResult.Deltas...)
	}
	return result, nil
}

func (d pivnetReleaseDiffer) DiffJob(job string) (Result, error) {
	// TODO: implement
	return Result{}, nil
}
