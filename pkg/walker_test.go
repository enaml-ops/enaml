package pkg_test

import (
	"io"
	"os"

	. "github.com/xchapter7x/enaml/pkg"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Walker", func() {
	Describe("Walking a package", func() {
		var (
			walker Walker
			file   FileEntry
			err    error
		)
		BeforeEach(func() {
			walker = NewWalker(readFixtureFile("redis-boshrelease-1.tgz"))
		})
		JustBeforeEach(func() {
			err = walker.Walk()
		})
		Context("When matching on release.MF", func() {
			BeforeEach(func() {
				walker.OnMatch("release.MF", func(f FileEntry) error {
					file = f
					return nil
				})
			})
			It("Calls back with release.MF", func() {
				Expect(file.FileName).To(Equal("./release.MF"))
			})
		})
		Context("When matching on ./jobs", func() {
			BeforeEach(func() {
				walker.OnMatch("/jobs/", func(f FileEntry) error {
					file = f
					return nil
				})
			})
			It("Calls back with redis.tgz", func() {
				Expect(file.FileName).To(Equal("./jobs/redis.tgz"))
			})
		})
	})
})

func readFixtureFile(filename string) io.Reader {
	f, err := os.Open("./fixtures/" + filename)
	Expect(err).NotTo(HaveOccurred())
	return f
}
