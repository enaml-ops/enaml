package diff

import "github.com/xchapter7x/enaml/pull"

type pivnetReleaseDiffer struct {
	ReleaseRepo pull.Release
	R1Path      string
	R2Path      string
}
