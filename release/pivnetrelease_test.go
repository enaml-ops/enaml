package release

import (
	"github.com/xchapter7x/enaml/pull"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Pivnetrelease", func() {
	var (
		err     error
		release *PivnetRelease
	)
	Context("Redis BOSH release 1.5.0", func() {
		BeforeEach(func() {
			releaseRepo := pull.Release{CacheDir: ".cache"}
			release, err = LoadPivnetRelease(releaseRepo, "../fixtures/p-redis-1.5.0.pivotal")
		})
		It("should not have errored", func() {
			Expect(err).NotTo(HaveOccurred())
		})
		It("should contain only the redis BOSH release", func() {
			Expect(release.BoshRelease).To(HaveLen(1))
			Expect(release.BoshRelease).To(HaveKey("redis"))
			Expect(release.BoshRelease["redis"].ReleaseManifest.Name).To(Equal("redis"))
		})
	})
})
