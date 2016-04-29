package diff_test

import (
	"fmt"

	. "github.com/xchapter7x/enaml/diff"
	"github.com/xchapter7x/enaml/pull"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Differ", func() {
	var (
		err    error
		differ Differ
		result *Result
	)

	Describe("Given a Diff func", func() {
		var releaseRepo pull.Release
		BeforeEach(func() {
			releaseRepo = pull.Release{CacheDir: "./cache"}
		})
		Context("When comparing BOSH Redis release 1 to 12", func() {
			BeforeEach(func() {
				differ, err = New(releaseRepo, "../fixtures/redis-boshrelease-1.tgz", "../fixtures/redis-boshrelease-12.tgz")
				Expect(err).NotTo(HaveOccurred())
				result, err = differ.Diff()
			})
			It("should not have errored", func() {
				Expect(err).NotTo(HaveOccurred())
			})
			It("should have differences", func() {
				Expect(len(result.DeltaJob)).To(BeNumerically(">", 0))
			})
		})
		Context("When comparing Pivnet Redis release 1.4.0 to 1.5.0", func() {
			BeforeEach(func() {
				differ, err = New(releaseRepo, "../fixtures/p-redis-1.4.0.pivotal", "../fixtures/p-redis-1.5.0.pivotal")
				Expect(err).NotTo(HaveOccurred())
				result, err = differ.Diff()
			})
			It("should not have errored", func() {
				Expect(err).NotTo(HaveOccurred())
			})
			It("should have differences", func() {
				Expect(len(result.DeltaJob)).To(BeNumerically(">", 0))
			})
		})
		Context("When comparing Pivnet Redis release 1.5.0 to Xip release 2.0.0", func() {
			BeforeEach(func() {
				differ, err = New(releaseRepo, "../fixtures/p-redis-1.5.0.pivotal", "../fixtures/p-xip-2.0.0.pivotal")
				Expect(err).NotTo(HaveOccurred())
				result, err = differ.Diff()
			})
			It("should not error", func() {
				Expect(err).ToNot(HaveOccurred())
			})
			It("should have differences", func() {
				Expect(len(result.DeltaJob)).To(BeNumerically(">", 0))
			})
		})
		Context("When comparing Redis job between Pivnet Redis release 1.4.0 and 1.5.0", func() {
			BeforeEach(func() {
				differ, err = New(releaseRepo, "../fixtures/p-redis-1.4.0.pivotal", "../fixtures/p-redis-1.5.0.pivotal")
				Expect(err).NotTo(HaveOccurred())
				result, err = differ.DiffJob("redis")
			})
			It("should not have errored", func() {
				Expect(err).NotTo(HaveOccurred())
			})
			It("should have differences", func() {
				Expect(len(result.DeltaJob)).To(BeNumerically(">", 0))
				for _, d := range result.DeltaJob {
					fmt.Println(d)
				}
			})
		})
	})
})
