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
			It("should have differences", func() {
				Expect(len(result.Deltas)).To(BeNumerically(">", 0))
			})
		})
		Context("Redis PivNet release 1.4.0 compared to 1.5.0", func() {
			BeforeEach(func() {
				releaseRepo := pull.Release{CacheDir: "./cache"}
				differ, err = New(releaseRepo, "./fixtures/p-redis-1.4.0.pivotal", "./fixtures/p-redis-1.5.0.pivotal")
				Expect(err).NotTo(HaveOccurred())
				result, err = differ.Diff()
			})
			It("should not have errored", func() {
				Expect(err).NotTo(HaveOccurred())
			})
			It("should have differences", func() {
				Expect(len(result.Deltas)).To(BeNumerically(">", 0))
			})
		})
	})
})
