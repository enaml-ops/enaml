package pull_test

import (
	"os"
	"path"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	. "github.com/xchapter7x/enaml/pull"
)

var _ = Describe("given Release object", func() {
	Describe("given a Pull method", func() {
		Context("when called on a valid release in the cache", func() {
			var (
				releaseName       = "concourse?v=1.1.0"
				controlReleaseURL = "https://bosh.io/d/github.com/concourse/" + releaseName
				controlCacheDir   = "fixtures"
				release           *Release
				filename          string
				err               error
			)

			BeforeEach(func() {
				release = NewRelease(controlCacheDir)
				filename, err = release.Pull(controlReleaseURL)
				Ω(err).ShouldNot(HaveOccurred())
			})

			It("then it should return a valid filename", func() {
				_, err = os.Stat(filename)
				Ω(err).ShouldNot(HaveOccurred())
				Ω(filename).ShouldNot(BeEmpty())
				Ω(filename).Should(Equal(path.Join(controlCacheDir, releaseName)))
			})
		})
	})
})
