package diff_test

import (
	. "github.com/xchapter7x/enaml/diff"
	"github.com/xchapter7x/enaml/pull"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Differ", func() {
	var (
		err    error
		differ Differ
		result Result
	)
	Describe("Given a Diff func", func() {
		Context("Redis BOSH release 1 compared to 12", func() {
			BeforeEach(func() {
				releaseRepo := pull.Release{CacheDir: "./cache"}
				differ, err = New(releaseRepo, "./fixtures/redis-boshrelease-1.tgz", "./fixtures/redis-boshrelease-12.tgz")
				Expect(err).NotTo(HaveOccurred())
				result, err = differ.Diff()
			})
			It("should not have errored", func() {
				Expect(err).NotTo(HaveOccurred())
			})
			It("should not have errored", func() {
				Expect(len(result.Deltas)).To(BeNumerically(">", 0))
			})
		})
	})
})
