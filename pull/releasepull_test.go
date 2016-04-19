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
			var release = "concourse?v=1.1.0"
			var controlReleaseURL = "https://bosh.io/d/github.com/concourse/" + release
			var controlCacheDir = "fixtures"
			var filename string

			BeforeEach(func() {
				release := NewRelease(controlCacheDir)
				filename = release.Pull(controlReleaseURL)
			})

			It("then it should return a valid filename", func() {
				_, err := os.Stat(filename)
				Ω(err).ShouldNot(HaveOccurred())
				Ω(filename).ShouldNot(BeEmpty())
				Ω(filename).Should(Equal(path.Join(controlCacheDir, release)))
			})
		})
	})
})
